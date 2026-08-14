"""TaskConsumer — 消费 Redis Stream 任务。

支持单流和多流消费模式：
- TaskConsumer: 单流消费（向后兼容，直跑测试用）
- MultiStreamConsumer: 多流并发消费（生产用）

使用 Redis Stream Consumer Group (XREADGROUP) 实现分布式消费。
"""
from __future__ import annotations

import logging
import threading
import time
from typing import Iterator

import redis

from scrapers import config
from scrapers.contracts.task_event import CrawlTaskMessage

logger = logging.getLogger(__name__)


class TaskConsumer:
    """单流消费者：消费 crawl:task:dispatch Stream，yield (msg_id, CrawlTaskMessage)。"""

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
        """持续消费消息，yield (msg_id, CrawlTaskMessage)。"""
        self.ensure_group()
        logger.info(
            "TaskConsumer started: stream=%s group=%s consumer=%s",
            self._stream, self._group, self._consumer,
        )

        yield from self._read_pending()

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


class MultiStreamConsumer:
    """多流消费者：同时消费多个 Stream，支持动态增减流。

    每个流在独立线程中消费，通过共享队列汇总消息。
    """

    def __init__(
        self,
        streams: list[str],
        group: str | None = None,
        consumer_name: str | None = None,
        redis_url: str | None = None,
    ):
        self._streams = list(streams)
        self._group = group or config.settings.consumer_group
        self._consumer = consumer_name or config.settings.consumer_name
        self._redis_url = redis_url or config.settings.redis_url
        self._conn: redis.Redis | None = None
        self._stream_consumers: dict[str, TaskConsumer] = {}
        self._lock = threading.Lock()

    @property
    def conn(self) -> redis.Redis:
        if self._conn is None:
            self._conn = redis.Redis.from_url(self._redis_url, decode_responses=True)
        return self._conn

    def _ensure_group(self, stream: str) -> None:
        """确保指定 Stream 上的 Consumer Group 存在。"""
        try:
            self.conn.xgroup_create(stream, self._group, id="0", mkstream=True)
            logger.info("Created consumer group %s on %s", self._group, stream)
        except redis.ResponseError as e:
            if "BUSYGROUP" in str(e):
                pass
            else:
                raise

    def consume(self, count: int = 5, block: int = 5000) -> Iterator[tuple[str, str, CrawlTaskMessage]]:
        """多流轮询消费，yield (stream_name, msg_id, CrawlTaskMessage)。

        使用 XREADGROUP 多流模式同时监听所有流。
        """
        if not self._streams:
            logger.warning("No streams to consume")
            return

        # 确保所有流的 Consumer Group 存在
        for stream in self._streams:
            self._ensure_group(stream)

        logger.info(
            "MultiStreamConsumer started: streams=%s group=%s consumer=%s",
            self._streams, self._group, self._consumer,
        )

        # 先消费各流的 PEL
        yield from self._read_all_pending()

        # 构建 streams dict 用于 XREADGROUP
        while True:
            try:
                with self._lock:
                    current_streams = list(self._streams)

                if not current_streams:
                    time.sleep(1)
                    continue

                streams_dict = {s: ">" for s in current_streams}
                results = self.conn.xreadgroup(
                    groupname=self._group,
                    consumername=self._consumer,
                    streams=streams_dict,
                    count=count,
                    block=block,
                )
                if not results:
                    continue

                for stream_name, messages in results:
                    for msg_id, fields in messages:
                        try:
                            task = CrawlTaskMessage.from_stream_message(fields)
                            yield stream_name, msg_id, task
                        except Exception:
                            logger.exception("Failed to parse message %s from %s", msg_id, stream_name)
                            self.ack(stream_name, msg_id)
            except redis.ConnectionError:
                logger.warning("Redis connection lost, reconnecting in 3s...")
                time.sleep(3)
                self._conn = None

    def _read_all_pending(self) -> Iterator[tuple[str, str, CrawlTaskMessage]]:
        """消费所有流的 PEL。"""
        with self._lock:
            current_streams = list(self._streams)

        for stream in current_streams:
            while True:
                results = self.conn.xreadgroup(
                    groupname=self._group,
                    consumername=self._consumer,
                    streams={stream: "0"},
                    count=50,
                )
                if not results:
                    break
                messages = results[0][1]
                if not messages:
                    break
                for msg_id, fields in messages:
                    if not fields:
                        continue
                    try:
                        task = CrawlTaskMessage.from_stream_message(fields)
                        logger.info("Replaying pending: stream=%s task=%s", stream, task.task_id)
                        yield stream, msg_id, task
                    except Exception:
                        logger.exception("Failed to parse pending message %s", msg_id)
                        self.ack(stream, msg_id)

    def ack(self, stream_name: str, msg_id: str) -> None:
        """确认消息已处理。"""
        self.conn.xack(stream_name, self._group, msg_id)

    def update_streams(self, new_streams: list[str]) -> None:
        """动态更新消费的 Stream 列表。"""
        with self._lock:
            old = set(self._streams)
            new = set(new_streams)
            added = new - old
            removed = old - new

            if added:
                for stream in added:
                    self._ensure_group(stream)
                logger.info("Adding streams: %s", added)
            if removed:
                logger.info("Removing streams: %s", removed)

            self._streams = list(new_streams)
            logger.info("Updated streams: %s", self._streams)

    def close(self) -> None:
        if self._conn is not None:
            self._conn.close()
            self._conn = None
