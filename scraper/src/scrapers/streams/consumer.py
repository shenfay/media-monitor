"""TaskConsumer — 消费 crawl:task:dispatch Stream，产出 CrawlTaskMessage。

使用 Redis Stream Consumer Group (XREADGROUP) 实现分布式消费。
"""
from __future__ import annotations

import logging
from typing import Iterator

import redis

from scrapers import config
from scrapers.contracts.task_event import CrawlTaskMessage

logger = logging.getLogger(__name__)


class TaskConsumer:
    """消费 crawl:task:dispatch，yield (msg_id, CrawlTaskMessage)。"""

    def __init__(
        self,
        redis_url: str | None = None,
        stream: str | None = None,
        group: str | None = None,
        consumer_name: str | None = None,
    ):
        self._redis_url = redis_url or config.settings.redis_url
        self._stream = stream or config.settings.dispatch_stream
        self._group = group or config.settings.consumer_group
        self._consumer = consumer_name or config.settings.consumer_name
        self._conn: redis.Redis | None = None

    @property
    def conn(self) -> redis.Redis:
        if self._conn is None:
            self._conn = redis.Redis.from_url(self._redis_url, decode_responses=True)
        return self._conn

    def ensure_group(self) -> None:
        """确保 Consumer Group 存在，已存在则忽略。"""
        try:
            self.conn.xgroup_create(self._stream, self._group, id="0", mkstream=True)
            logger.info("Created consumer group %s on %s", self._group, self._stream)
        except redis.ResponseError as e:
            if "BUSYGROUP" in str(e):
                logger.debug("Consumer group %s already exists", self._group)
            else:
                raise

    def consume(self, count: int = 5, block: int = 5000) -> Iterator[tuple[str, CrawlTaskMessage]]:
        """持续消费消息，yield (msg_id, CrawlTaskMessage)。

        调用方处理完成后需显式调用 ack()。
        先消费历史未 ACK 消息（id=">"之前用 "0"），再消费新消息（id=">"）。
        """
        self.ensure_group()
        logger.info(
            "TaskConsumer started: stream=%s group=%s consumer=%s",
            self._stream, self._group, self._consumer,
        )

        # 先消费 PEL 中未 ACK 的历史消息
        yield from self._read_pending()

        # 持续消费新消息
        while True:
            try:
                results = self.conn.xreadgroup(
                    groupname=self._group,
                    consumername=self._consumer,
                    streams={self._stream: ">"},
                    count=count,
                    block=block,
                )
                if not results:
                    continue

                for _stream_name, messages in results:
                    for msg_id, fields in messages:
                        try:
                            task = CrawlTaskMessage.from_stream_message(fields)
                            yield msg_id, task
                        except Exception:
                            logger.exception("Failed to parse message %s", msg_id)
                            self.ack(msg_id)
            except redis.ConnectionError:
                logger.warning("Redis connection lost, reconnecting in 3s...")
                import time
                time.sleep(3)
                self._conn = None

    def _read_pending(self) -> Iterator[tuple[str, CrawlTaskMessage]]:
        """消费 PEL 中尚未 ACK 的历史消息（重启恢复用）。"""
        while True:
            results = self.conn.xreadgroup(
                groupname=self._group,
                consumername=self._consumer,
                streams={self._stream: "0"},
                count=50,
            )
            if not results:
                break

            messages = results[0][1]
            if not messages:
                break

            for msg_id, fields in messages:
                if not fields:
                    # 已被 ACK 但仍在 PEL 中（空结果）
                    continue
                try:
                    task = CrawlTaskMessage.from_stream_message(fields)
                    logger.info("Replaying pending message %s (task=%s)", msg_id, task.task_id)
                    yield msg_id, task
                except Exception:
                    logger.exception("Failed to parse pending message %s", msg_id)
                    self.ack(msg_id)

    def ack(self, msg_id: str) -> None:
        """确认消息已处理。"""
        self.conn.xack(self._stream, self._group, msg_id)

    def close(self) -> None:
        if self._conn is not None:
            self._conn.close()
            self._conn = None
