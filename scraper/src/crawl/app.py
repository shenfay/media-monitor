"""CLI 入口（typer）。

子命令：
- serve:  常驻消费 Redis Stream（生产用）
- run:    直跑某个 source（本地测试用）
- health: 健康检查（Docker HEALTHCHECK 用）
"""
from __future__ import annotations

import json
import logging
import sys

import typer

app = typer.Typer(help="MediaMonitor 抓取服务")


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
