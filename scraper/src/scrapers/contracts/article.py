"""Article 文章契约。

Python 通过 crawl:article:ingest Stream 分批回传给 Go。
"""
from __future__ import annotations

import dataclasses
import datetime


def utc_now_iso() -> str:
    """返回当前 UTC 时间的 ISO 8601 字符串。"""
    return datetime.datetime.now(datetime.timezone.utc).isoformat()


@dataclasses.dataclass
class Article:
    """单篇文章。"""

    source_id: str
    platform: str
    title: str
    url: str
    external_id: str = ""                 # 平台侧文章 ID（原 aid）
    author: str = ""
    source_name: str = ""
    summary: str = ""
    content: str = ""                     # 正文（HTML 或 Markdown）
    published_at: str = ""                # ISO 8601
    images: list = dataclasses.field(default_factory=list)
    extra: dict = dataclasses.field(default_factory=dict)   # 互动指标等
    fetched_at: str = dataclasses.field(default_factory=utc_now_iso)

    def to_dict(self) -> dict:
        return {
            "source_id": self.source_id,
            "platform": self.platform,
            "title": self.title,
            "url": self.url,
            "external_id": self.external_id,
            "author": self.author,
            "source_name": self.source_name,
            "summary": self.summary,
            "content": self.content,
            "published_at": self.published_at,
            "images": self.images,
            "extra": self.extra,
            "fetched_at": self.fetched_at,
        }

    @classmethod
    def from_dict(cls, d: dict) -> Article:
        return cls(
            source_id=d.get("source_id", ""),
            platform=d.get("platform", ""),
            title=d.get("title", ""),
            url=d.get("url", ""),
            external_id=d.get("external_id", ""),
            author=d.get("author", ""),
            source_name=d.get("source_name", ""),
            summary=d.get("summary", ""),
            content=d.get("content", ""),
            published_at=d.get("published_at", ""),
            images=d.get("images") or [],
            extra=d.get("extra") or {},
            fetched_at=d.get("fetched_at", ""),
        )


@dataclasses.dataclass
class ArticleBatch:
    """分批回传的文章批次（对应一条 crawl:article:ingest 消息）。"""

    task_id: str
    source_id: str
    phase: str                            # "list" | "detail"
    batch_seq: int
    articles: list[Article] = dataclasses.field(default_factory=list)

    def to_dict(self) -> dict:
        return {
            "task_id": self.task_id,
            "source_id": self.source_id,
            "phase": self.phase,
            "batch_seq": self.batch_seq,
            "articles": [a.to_dict() for a in self.articles],
        }
