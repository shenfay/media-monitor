"""CLI 入口（typer）。

子命令：
- serve:          常驻消费 Redis Stream（生产用）
- run:            直跑某个 source（本地测试用）
- history:        历史数据回刷（从 DB 读配置，深度翻页，经 Stream 入库）
- detail-worker:  详情抓取 Worker（消费 crawl:detail:queue，抓取正文后回传）
- health:         健康检查（Docker HEALTHCHECK 用）
"""
from __future__ import annotations

import json
import logging
import sys
import uuid

import typer

app = typer.Typer(help="MediaMonitor 抓取服务")


@app.callback()
def main(
    ctx: typer.Context,
    source: str = typer.Option(None, "--source", "-s", help="数据源 ID（如 source_hq_world）"),
) -> None:
    """全局选项。"""
    ctx.ensure_object(dict)
    ctx.obj["source"] = source


def _setup_logging(verbose: bool = False) -> None:
    level = logging.DEBUG if verbose else logging.INFO
    logging.basicConfig(
        level=level,
        format="%(asctime)s %(levelname)-5s %(name)s: %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )


@app.command()
def serve(
    verbose: bool = typer.Option(False, "--verbose", "-v", help="Debug 日志"),
) -> None:
    """常驻消费 Redis Stream，按 task 调对应 adapter，回写数据。"""
    _setup_logging(verbose)
    logger = logging.getLogger("crawl.serve")

    # 触发适配器自动注册
    import crawl.adapters  # noqa: F401

    from crawl.adapters.registry import get_all_required_tags, list_metadata
    from crawl.core.executor import TaskExecutor
    from crawl.streams.consumer import MultiStreamConsumer
    from crawl.streams.event_emitter import TaskEventEmitter
    from crawl.streams.publisher import ArticlePublisher
    from crawl.worker import ControlListener, StreamRouter, WorkerRegistration

    # 收集 adapter 元数据与能力标签
    adapters_meta = list_metadata()
    capabilities = get_all_required_tags()
    logger.info("Registered adapters: %s", list(adapters_meta.keys()))
    logger.info("Worker capabilities: %s", capabilities)

    # Worker 注册与心跳
    registration = WorkerRegistration(
        adapters_metadata=list(adapters_meta.values()),
        capabilities=capabilities,
    )
    registration.register()
    registration.start_heartbeat()

    # Stream 路由：计算当前 Worker 需要消费的流
    import redis as redis_lib

    from crawl import config
    redis_conn = redis_lib.Redis.from_url(config.settings.redis_url, decode_responses=True)
    adapter_required_tags = {
        name: meta.get("required_tags", [])
        for name, meta in adapters_meta.items()
    }
    stream_router = StreamRouter(
        redis_conn=redis_conn,
        capabilities=capabilities,
        adapter_required_tags=adapter_required_tags,
    )
    initial_streams = stream_router.compute_streams()
    logger.info("Initial streams: %s", initial_streams)

    consumer = MultiStreamConsumer(
        streams=initial_streams,
        consumer_name=registration.worker_id,
    )
    publisher = ArticlePublisher()
    emitter = TaskEventEmitter()
    executor = TaskExecutor(publisher, emitter)

    # 控制状态
    paused = False
    shutdown_requested = False

    def on_pause():
        nonlocal paused
        paused = True
        logger.info("Worker paused")

    def on_resume():
        nonlocal paused
        paused = False
        logger.info("Worker resumed")

    def on_shutdown():
        nonlocal shutdown_requested
        shutdown_requested = True
        logger.info("Shutdown requested")

    def on_recalculate():
        new_streams = stream_router.compute_streams()
        consumer.update_streams(new_streams)
        logger.info("Streams recalculated: %s", new_streams)

    # 控制监听器
    control = ControlListener(
        worker_id=registration.worker_id,
        redis_conn=redis_conn,
        on_pause=on_pause,
        on_resume=on_resume,
        on_shutdown=on_shutdown,
        on_recalculate=on_recalculate,
    )
    control.start()

    logger.info("Scraper worker starting...")
    processed = 0
    try:
        for stream_name, msg_id, task in consumer.consume():
            if shutdown_requested:
                logger.info("Shutdown requested, stopping...")
                break
            # 暂停时跳过任务执行，但继续消费（消息会留在 PEL 中）
            if paused:
                import time as _time
                while paused and not shutdown_requested:
                    _time.sleep(1)
                if shutdown_requested:
                    break
            logger.info(
                "Received task: stream=%s id=%s source=%s",
                stream_name, task.task_id, task.source_id,
            )
            registration.update_status(current_task=task.task_id)
            try:
                executor.execute(task)
            except Exception:
                logger.exception("Task execution failed: %s", task.task_id)
                emitter.task_failed(task.task_id, "unhandled error")
            consumer.ack(stream_name, msg_id)
            processed += 1
            registration.update_status(current_task="", processed_count=processed)
    except KeyboardInterrupt:
        logger.info("Shutting down...")
    finally:
        control.stop()
        registration.deregister()
        consumer.close()
        publisher.close()
        emitter.close()
        redis_conn.close()


@app.command()
def run(
    source: str = typer.Option("huanqiu", help="数据源名称或 platform_type"),
    limit: int = typer.Option(200, help="最大文章数"),
    with_body: bool = typer.Option(False, help="是否抓取正文"),
    out: str = typer.Option(None, help="输出 JSON 文件路径（默认 stdout）"),
    verbose: bool = typer.Option(False, "--verbose", "-v", help="Debug 日志"),
) -> None:
    """直跑某个 source（本地测试用，不依赖 Redis）。"""
    _setup_logging(verbose)
    logger = logging.getLogger("crawl.run")

    import crawl.adapters  # noqa: F401

    from crawl.cleaning.normalizer import normalize_article
    from crawl.cleaning.validator import validate_article
    from crawl.contracts.source import Source
    from crawl.adapters.registry import get_adapter

    adapter = get_adapter(source, source)
    if adapter is None:
        logger.error("No adapter found for: %s", source)
        raise typer.Exit(1)

    # 构造一个最小 Source（直跑模式）
    src = Source(id="local", name=source, platform_type=source)
    logger.info("Running adapter=%s limit=%d with_body=%s", adapter.__class__.__name__, limit, with_body)

    articles = adapter.fetch_list(src, limit=limit)
    articles = [normalize_article(a) for a in articles]
    articles = [a for a in articles if not validate_article(a)]

    if with_body:
        for a in articles:
            try:
                adapter.fetch_detail(src, a)
                normalize_article(a)
            except Exception:
                logger.exception("Detail fetch failed: %s", a.url)

    result = {
        "source": source,
        "count": len(articles),
        "articles": [a.to_dict() for a in articles],
    }
    output = json.dumps(result, ensure_ascii=False, indent=2)

    if out:
        with open(out, "w", encoding="utf-8") as f:
            f.write(output)
        logger.info("Output written to %s (%d articles)", out, len(articles))
    else:
        print(output)


@app.command()
def history(
    ctx: typer.Context,
    since: str = typer.Option(..., help="起始日期（如 2025-02-16）"),
    until: str = typer.Option("", help="截止日期（如 2025-03-01），用于分段查询"),
    delay: float = typer.Option(0.5, help="请求间隔秒数"),
    start_offset: int = typer.Option(0, help="起始 offset，用于断点续传"),
    per_node: bool = typer.Option(False, "--per-node", help="逐子节点查询，每个子节点独立翻页（绕过 offset=10000 限制）"),
    no_source_filter: bool = typer.Option(False, "--no-source-filter", help="清空 source_filter，抓取所有来源的文章"),
    with_body: bool = typer.Option(False, help="是否同步抓取正文（内联模式）"),
    async_body: bool = typer.Option(False, "--async-body", help="异步详情模式：列表入库后发布详情任务到 crawl:detail:queue，由 detail-worker 消费"),
    verbose: bool = typer.Option(False, "--verbose", "-v", help="Debug 日志"),
) -> None:
    """历史数据回刷：从数据库读取数据源配置，深度翻页抓取，经 Stream 管道写入数据库。"""
    _setup_logging(verbose)
    logger = logging.getLogger("crawl.history")

    source = ctx.obj.get("source")
    if not source:
        logger.error("请指定 --source 参数，如: python -m crawl --source source_hq_world history --since 2025-02-16")
        raise typer.Exit(1)

    from crawl import config as crawl_config

    # ── 1. 连接数据库，读取数据源配置 ──────────────────────
    dsn = crawl_config.settings.database_dsn
    if not dsn:
        logger.error("MM_DATABASE_DSN 未配置，无法读取数据源")
        raise typer.Exit(1)

    try:
        import psycopg2
    except ImportError:
        logger.error("psycopg2 未安装，请执行: pip install -e '.[db]'")
        raise typer.Exit(1)

    try:
        conn = psycopg2.connect(dsn)
        cur = conn.cursor()
        cur.execute(
            "SELECT id, name, platform_type, base_url, list_endpoint, "
            "nodes, source_filter, months, tags "
            "FROM crawl_sources WHERE id = %s",
            (source,),
        )
        row = cur.fetchone()
        if not row:
            logger.error("数据源不存在: %s", source)
            cur.close()
            conn.close()
            raise typer.Exit(1)

        cols = ["id", "name", "platform_type", "base_url", "list_endpoint",
                "nodes", "source_filter", "months", "tags"]
        data = dict(zip(cols, row))
        cur.close()
        conn.close()
    except Exception as e:
        logger.error("数据库连接失败: %s", e)
        raise typer.Exit(1)

    # 解析 JSON 字段
    import json as _json
    nodes = []
    tags = []
    try:
        nodes = _json.loads(data["nodes"]) if data["nodes"] else []
    except (ValueError, TypeError):
        pass
    try:
        tags = _json.loads(data["tags"]) if data["tags"] else []
    except (ValueError, TypeError):
        pass

    # 构造 Source 对象
    from crawl.contracts.source import Source
    extra = {"since": since, "delay": delay, "per_node": per_node}
    if until:
        extra["until"] = until
    src = Source(
        id=data["id"],
        name=data["name"],
        platform_type=data["platform_type"],
        base_url=data["base_url"] or "",
        list_endpoint=data["list_endpoint"] or "",
        nodes=nodes,
        source_filter="" if no_source_filter else (data["source_filter"] or ""),
        months=data["months"] or 6,
        tags=tags,
        extra=extra,
    )

    logger.info(
        "数据源已加载: id=%s name=%s list_endpoint=%s since=%s until=%s delay=%.1fs per_node=%s",
        src.id, src.name, src.list_endpoint, since, until or "无", delay, per_node,
    )

    # ── 2. 深度翻页抓取（边翻边写） ──────────────────────
    from crawl.adapters.huanqiu.history import HuanqiuHistoryAdapter
    from crawl.cleaning.normalizer import normalize_article
    from crawl.cleaning.validator import validate_article
    from crawl.streams.publisher import ArticlePublisher

    task_id = f"history-{uuid.uuid4().hex[:12]}"
    publisher = ArticlePublisher()
    total_fetched = 0
    total_valid = 0
    total_batches = 0

    def on_page(page_articles, page_num):
        """每页抓取完后立即标准化、校验、发布到 Stream。"""
        nonlocal total_valid, total_batches
        normalized = [normalize_article(a) for a in page_articles]
        valid = []
        for a in normalized:
            errs = validate_article(a)
            if errs:
                logger.warning("文章校验失败: url=%s errors=%s", a.url, errs)
                continue
            valid.append(a)
        if valid:
            batches = publisher.publish_list(task_id, src.id, valid)
            total_valid += len(valid)
            total_batches += batches
            logger.info(
                "  ↳ Page %d 发布 %d 条有效文章 (%d batches)",
                page_num, len(valid), batches,
            )

    adapter = HuanqiuHistoryAdapter()
    try:
        all_articles = adapter.fetch_list(src, on_page=on_page, start_offset=start_offset)
    finally:
        publisher.close()

    total_fetched = len(all_articles)
    logger.info("抓取完成: 总计 %d 条, 有效 %d 条 (已实时入库)", total_fetched, total_valid)

    if total_valid == 0:
        logger.info("无有效文章，退出")
        return

    # ── 3. 可选：抓取正文（同步内联模式） ───────────────────
    if with_body:
        logger.info("开始同步抓取正文...")
        detail_ok = 0
        detail_fail = 0
        for a in all_articles:
            try:
                adapter.fetch_detail(src, a)
                normalize_article(a)
                errs = validate_article(a)
                if errs:
                    detail_fail += 1
                else:
                    detail_ok += 1
            except Exception:
                logger.exception("正文抓取失败: %s", a.url)
                detail_fail += 1
        logger.info("正文完成: 成功 %d, 失败 %d", detail_ok, detail_fail)

    # ── 3b. 可选：异步详情模式（发布到详情队列） ─────────────
    if async_body:
        from crawl.streams.publisher import DetailTaskPublisher
        logger.info("发布详情任务到 crawl:detail:queue ...")
        detail_publisher = DetailTaskPublisher()
        try:
            batches = detail_publisher.publish(
                task_id=task_id,
                source_id=src.id,
                adapter_name="huanqiu_history",
                articles=all_articles,
            )
            logger.info("已发布 %d 批详情任务 (%d 篇文章)", batches, len(all_articles))
        finally:
            detail_publisher.close()

    # ── 4. 汇总 ──────────────────────────────────────────
    dates = [a.published_at for a in all_articles if a.published_at]
    logger.info("=" * 50)
    logger.info("回刷完成:")
    logger.info("  数据源:   %s (%s)", src.name, src.id)
    logger.info("  起始日期: %s", since)
    logger.info("  总抓取:   %d", total_fetched)
    logger.info("  有效入库: %d (%d batches)", total_valid, total_batches)
    if dates:
        logger.info("  日期范围: %s ~ %s", min(dates), max(dates))
    logger.info("  task_id:  %s", task_id)
    logger.info("=" * 50)


@app.command()
def detail_worker(
    delay: float = typer.Option(0.5, help="每篇请求间隔秒数"),
    consumer_name: str = typer.Option("", help="消费者名称（默认使用配置中的 MM_CONSUMER_NAME）"),
    verbose: bool = typer.Option(False, "--verbose", "-v", help="Debug 日志"),
) -> None:
    """详情抓取 Worker：消费 crawl:detail:queue，抓取正文后回传。

    独立运行，与列表抓取完全解耦。可启动多个实例并行消费。
    """
    _setup_logging(verbose)
    logger = logging.getLogger("crawl.detail_worker")

    # 触发适配器自动注册
    import crawl.adapters  # noqa: F401

    from crawl.detail.fetcher import DetailFetcher

    kwargs = {"delay": delay}
    if consumer_name:
        kwargs["consumer_name"] = consumer_name

    fetcher = DetailFetcher(**kwargs)
    logger.info("Starting detail worker with delay=%.1fs", delay)
    try:
        fetcher.run()
    except KeyboardInterrupt:
        logger.info("Detail worker shutting down...")


@app.command()
def health() -> None:
    """健康检查（Docker HEALTHCHECK 用）。"""
    import crawl.adapters  # noqa: F401
    from crawl.adapters.registry import list_registered

    registered = list_registered()
    if not registered:
        print("UNHEALTHY: no adapters registered")
        sys.exit(1)

    result = {"status": "ok", "adapters": registered}
    print(json.dumps(result))


def cli() -> None:
    """typer 入口。"""
    app()
