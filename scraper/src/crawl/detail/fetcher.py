"""DetailFetcher — 消费 crawl:detail:queue，抓取正文后回传。

独立运行的详情抓取 Worker：
1. 消费 crawl:detail:queue 中的详情任务
2. 按 adapter_name 查找对应适配器
3. 调用 adapter.fetch_detail() 抓取正文
4. 将结果通过 crawl:article:ingest (phase=detail) 回传
"""
from __future__ import annotations

import json
import logging
import time

import redis

from crawl import config
from crawl.adapters.registry import _REGISTRY
from crawl.cleaning.normalizer import normalize_article
from crawl.cleaning.validator import validate_article
from crawl.contracts.article import Article
from crawl.contracts.detail_task import DetailTask, DetailTaskBatch
from crawl.contracts.source import Source
from crawl.streams.publisher import ArticlePublisher

logger = logging.getLogger(__name__)


class DetailFetcher:
    """详情抓取 Worker：消费详情队列，抓取正文，回传结果。"""

    def __init__(
        self,
        redis_url: str | None = None,
        stream: str | None = None,
        group: str | None = None,
        consumer_name: str | None = None,
        delay: float = 0.5,
    ):
        self._redis_url = redis_url or config.settings.redis_url
        self._stream = stream or config.settings.detail_queue
        self._group = group or "crawl:detail:worker"
        self._consumer = consumer_name or config.settings.consumer_name
        self._delay = delay
        self._conn: redis.Redis | None = None

    @property
    def conn(self) -> redis.Redis:
        if self._conn is None:
            self._conn = redis.Redis.from_url(self._redis_url, decode_responses=True)
        return self._conn

    def _ensure_group(self) -> None:
        """确保 Consumer Group 存在。"""
        try:
            self.conn.xgroup_create(self._stream, self._group, id="0", mkstream=True)
            logger.info("Created consumer group %s on %s", self._group, self._stream)
        except redis.ResponseError as e:
            if "BUSYGROUP" in str(e):
                logger.debug("Consumer group %s already exists", self._group)
            else:
                raise

    def run(self) -> None:
        """启动消费循环。"""
        self._ensure_group()
        logger.info(
            "DetailFetcher started: stream=%s group=%s consumer=%s delay=%.1fs",
            self._stream, self._group, self._consumer, self._delay,
        )

        # 先处理积压的 pending 消息
        self._consume_pending()

        # 持续消费新消息
        while True:
            try:
                results = self.conn.xreadgroup(
                    groupname=self._group,
                    consumername=self._consumer,
                    streams={self._stream: ">"},
                    count=5,
                    block=5000,
                )
                if not results:
                    continue

                for _stream_name, messages in results:
                    for msg_id, fields in messages:
                        try:
                            self._handle_message(fields)
                            self.conn.xack(self._stream, self._group, msg_id)
                        except Exception:
                            logger.exception("Failed to handle message %s", msg_id)
            except redis.ConnectionError:
                logger.warning("Redis connection lost, reconnecting in 5s...")
                time.sleep(5)
                self._conn = None
            except KeyboardInterrupt:
                logger.info("DetailFetcher shutting down (keyboard interrupt)")
                break

    def _consume_pending(self) -> None:
        """处理积压的 pending 消息（之前未 ACK 的）。"""
        while True:
            try:
                results = self.conn.xreadgroup(
                    groupname=self._group,
                    consumername=self._consumer,
                    streams={self._stream: "0"},
                    count=10,
                )
                if not results:
                    break
                messages = results[0][1]
                if not messages:
                    break
                for msg_id, fields in messages:
                    if not fields:  # pending 但已无数据
                        self.conn.xack(self._stream, self._group, msg_id)
                        continue
                    try:
                        self._handle_message(fields)
                        self.conn.xack(self._stream, self._group, msg_id)
                    except Exception:
                        logger.exception("Failed to handle pending message %s", msg_id)
            except Exception:
                logger.exception("Error consuming pending messages")
                break

    def _handle_message(self, fields: dict) -> None:
        """处理单条详情任务批次消息。"""
        payload_str = fields.get("payload", "")
        if not payload_str:
            return

        data = json.loads(payload_str)
        batch = DetailTaskBatch(
            task_id=data.get("task_id", ""),
            source_id=data.get("source_id", ""),
            adapter_name=data.get("adapter_name", ""),
            tasks=[DetailTask.from_dict(t) for t in data.get("tasks", [])],
        )

        if not batch.tasks:
            return

        logger.info(
            "Processing detail batch: task=%s source=%s adapter=%s tasks=%d",
            batch.task_id, batch.source_id, batch.adapter_name, len(batch.tasks),
        )

        # 查找适配器
        adapter_cls = _REGISTRY.get(batch.adapter_name)
        if not adapter_cls:
            logger.error("Adapter not found: %s", batch.adapter_name)
            return

        adapter = adapter_cls()

        # 构造最小 Source（仅用于 fetch_detail）
        source = Source(
            id=batch.source_id,
            name=batch.adapter_name,
            platform_type="news",
        )

        # 逐篇抓取正文
        results: list[Article] = []
        for task in batch.tasks:
            article = Article(
                source_id=task.source_id,
                platform=task.platform,
                title=task.title,
                url=task.url,
                external_id=task.external_id,
            )
            try:
                updated = adapter.fetch_detail(source, article)
                updated = normalize_article(updated)
                errs = validate_article(updated)
                if errs:
                    logger.warning("Detail validation failed: url=%s errors=%s", task.url, errs)
                elif updated.content:
                    results.append(updated)
                else:
                    logger.debug("No content extracted: url=%s", task.url)
            except Exception:
                logger.exception("Detail fetch failed: url=%s", task.url)

            time.sleep(self._delay)

        # 回传结果
        if results:
            publisher = ArticlePublisher()
            try:
                batches = publisher.publish_detail(
                    task_id=batch.task_id,
                    source_id=batch.source_id,
                    articles=results,
                )
                logger.info(
                    "Published %d detail results (%d articles) for task=%s",
                    batches, len(results), batch.task_id,
                )
            finally:
                publisher.close()
        else:
            logger.info("No content extracted for task=%s", batch.task_id)


def normalize_articles(articles: list[Article]) -> list[Article]:
    """批量标准化文章。"""
    return [normalize_article(a) for a in articles]
