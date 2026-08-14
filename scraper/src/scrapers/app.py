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
    logger = logging.getLogger("scrapers.serve")

    # 触发适配器自动注册
    import scrapers.adapters  # noqa: F401

    from scrapers.core.executor import TaskExecutor
    from scrapers.core.registry import list_registered
    from scrapers.streams.consumer import TaskConsumer
    from scrapers.streams.event_emitter import TaskEventEmitter
    from scrapers.streams.publisher import ArticlePublisher

    registered = list_registered()
    logger.info("Registered adapters: %s", registered)

    consumer = TaskConsumer()
    publisher = ArticlePublisher()
    emitter = TaskEventEmitter()
    executor = TaskExecutor(publisher, emitter)

    logger.info("Scraper worker starting...")
    try:
        for msg_id, task in consumer.consume():
            logger.info("Received task: id=%s source=%s", task.task_id, task.source_id)
            try:
                executor.execute(task)
            except Exception:
                logger.exception("Task execution failed: %s", task.task_id)
                emitter.task_failed(task.task_id, "unhandled error")
            consumer.ack(msg_id)
    except KeyboardInterrupt:
        logger.info("Shutting down...")
    finally:
        consumer.close()
        publisher.close()
        emitter.close()


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
    logger = logging.getLogger("scrapers.run")

    import scrapers.adapters  # noqa: F401

    from scrapers.cleaning.normalizer import normalize_article
    from scrapers.cleaning.validator import validate_article
    from scrapers.contracts.source import Source
    from scrapers.core.registry import get_adapter

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
    import scrapers.adapters  # noqa: F401
    from scrapers.core.registry import list_registered

    registered = list_registered()
    if not registered:
        print("UNHEALTHY: no adapters registered")
        sys.exit(1)

    result = {"status": "ok", "adapters": registered}
    print(json.dumps(result))


def cli() -> None:
    """typer 入口。"""
    app()
