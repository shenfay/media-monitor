"""任务消息与事件契约。

- CrawlTaskMessage: Go 通过 crawl:task:dispatch 下发（含 Source 配置）。
- TaskEvent: Python 通过 crawl:task:event 回传状态。
"""
from __future__ import annotations

import dataclasses
import json

from scrapers.contracts.article import utc_now_iso
from scrapers.contracts.source import Source


@dataclasses.dataclass
class CrawlTaskMessage:
    """从 crawl:task:dispatch 消费到的任务消息。"""

    task_id: str
    source_id: str
    source: Source
    params: dict = dataclasses.field(default_factory=dict)

    @classmethod
    def from_stream_message(cls, fields: dict) -> CrawlTaskMessage:
        """从 Redis Stream 消息字段解析。

        Redis Stream 消息的 key/value 可能是 bytes 或 str，统一处理。
        """
        def get_str(key: str) -> str:
            """从 fields 中取值，兼容 bytes/str。"""
            for k in (key, key.encode("utf-8")):
                if k in fields:
                    v = fields[k]
                    if isinstance(v, bytes):
                        return v.decode("utf-8", "replace")
                    return str(v or "")
            return ""

        task_id = get_str("task_id")
        source_id = get_str("source_id")

        # payload 是 JSON 字符串，内含 source + params
        raw_payload = get_str("payload") or get_str("data") or "{}"
        try:
            payload = json.loads(raw_payload)
        except (json.JSONDecodeError, TypeError):
            payload = {}

        # Source 配置内嵌在 payload 中
        source_data = payload.get("source") or {}
        if isinstance(source_data, str):
            try:
                source_data = json.loads(source_data)
            except (json.JSONDecodeError, TypeError):
                source_data = {}
        source = Source.from_dict(source_data) if source_data else Source(id=source_id, name="")

        # 参数
        params = payload.get("params") or {}
        if isinstance(params, str):
            try:
                params = json.loads(params)
            except (json.JSONDecodeError, TypeError):
                params = {}

        return cls(task_id=task_id, source_id=source_id, source=source, params=params)


@dataclasses.dataclass
class TaskEvent:
    """任务状态事件（通过 crawl:task:event 回传）。"""

    task_id: str
    type: str                             # status | phase_start | phase_done |
                                          # list_synced | detail_progress |
                                          # task_done | task_failed
    phase: str = ""                       # list | detail | ""
    status: str = ""                      # running | success | partial | failed
    total: int = 0
    list_count: int = 0
    detail_count: int = 0
    detail_failed: int = 0
    error: str = ""
    timestamp: str = ""

    def to_dict(self) -> dict:
        return {
            "task_id": self.task_id,
            "type": self.type,
            "phase": self.phase,
            "status": self.status,
            "total": self.total,
            "list_count": self.list_count,
            "detail_count": self.detail_count,
            "detail_failed": self.detail_failed,
            "error": self.error,
            "timestamp": self.timestamp or utc_now_iso(),
        }
