"""ArticlePublisher — 分批 XADD crawl:article:ingest。

将文章列表/详情分批发送到 Go 侧消费，每批 ≤ batch_size 篇。
"""
from __future__ import annotations

import json
import logging

import redis

from scrapers import config
from scrapers.contracts.article import Article, ArticleBatch

logger = logging.getLogger(__name__)

# 单消息最大字节数（500KB）
_MAX_MESSAGE_BYTES = 500 * 1024


class ArticlePublisher:
    """分批发布文章到 crawl:article:ingest Stream。"""

    def __init__(
        self,
        redis_url: str | None = None,
        stream: str | None = None,
        max_batch: int | None = None,
        max_stream_len: int | None = None,
    ):
        self._redis_url = redis_url or config.settings.redis_url
        self._stream = stream or config.settings.article_stream
        self._max_batch = max_batch or config.settings.batch_size
        self._max_stream_len = max_stream_len or config.settings.max_stream_len
        self._conn: redis.Redis | None = None

    @property
    def conn(self) -> redis.Redis:
        if self._conn is None:
            self._conn = redis.Redis.from_url(self._redis_url, decode_responses=True)
        return self._conn

    def publish_list(
        self,
        task_id: str,
        source_id: str,
        articles: list[Article],
    ) -> int:
        """发布列表阶段的文章（phase="list"），返回发送批次数。"""
        return self._publish(task_id, source_id, "list", articles)

    def publish_detail(
        self,
        task_id: str,
        source_id: str,
        articles: list[Article],
        batch_seq_offset: int = 0,
    ) -> int:
        """发布详情阶段的文章（phase="detail"），返回发送批次数。"""
        return self._publish(task_id, source_id, "detail", articles, batch_seq_offset)

    def _publish(
        self,
        task_id: str,
        source_id: str,
        phase: str,
        articles: list[Article],
        batch_seq_offset: int = 0,
    ) -> int:
        """将文章列表分批 XADD，返回实际发送批次数。"""
        if not articles:
            return 0

        batches_sent = 0
        for i in range(0, len(articles), self._max_batch):
            chunk = articles[i : i + self._max_batch]
            batch_seq = batch_seq_offset + (i // self._max_batch)

            batch = ArticleBatch(
                task_id=task_id,
                source_id=source_id,
                phase=phase,
                batch_seq=batch_seq,
                articles=chunk,
            )
            payload = json.dumps(batch.to_dict(), ensure_ascii=False)

            # 检查消息大小
            if len(payload.encode("utf-8")) > _MAX_MESSAGE_BYTES:
                logger.warning(
                    "Batch %d too large (%d bytes), splitting further",
                    batch_seq, len(payload.encode("utf-8")),
                )
                # 降级为逐篇发送
                for article in chunk:
                    single = ArticleBatch(
                        task_id=task_id,
                        source_id=source_id,
                        phase=phase,
                        batch_seq=batch_seq,
                        articles=[article],
                    )
                    single_payload = json.dumps(single.to_dict(), ensure_ascii=False)
                    self._xadd(single_payload)
                    batches_sent += 1
            else:
                self._xadd(payload)
                batches_sent += 1

        logger.info(
            "Published %d batches for task=%s phase=%s (%d articles)",
            batches_sent, task_id, phase, len(articles),
        )
        return batches_sent

    def _xadd(self, payload: str) -> None:
        """执行 XADD，带 MAXLEN 控制。"""
        self.conn.xadd(
            self._stream,
            {"payload": payload},
            maxlen=self._max_stream_len,
            approximate=True,
        )

    def close(self) -> None:
        if self._conn is not None:
            self._conn.close()
            self._conn = None
