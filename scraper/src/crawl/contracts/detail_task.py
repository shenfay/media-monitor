"""详情任务合约：crawl:detail:queue 消息结构。

列表抓取完成后，将需要抓取正文的文章信息投入详情队列，
由独立的 DetailFetcher 消费并抓取正文。
"""
from __future__ import annotations

import dataclasses


@dataclasses.dataclass
class DetailTask:
    """单篇详情抓取任务。"""

    source_id: str
    platform: str
    url: str
    external_id: str = ""
    title: str = ""
    adapter_name: str = ""  # 适配器名称（如 huanqiu, huanqiu_history）

    def to_dict(self) -> dict:
        return {
            "source_id": self.source_id,
            "platform": self.platform,
            "url": self.url,
            "external_id": self.external_id,
            "title": self.title,
            "adapter_name": self.adapter_name,
        }

    @classmethod
    def from_dict(cls, d: dict) -> DetailTask:
        return cls(
            source_id=d.get("source_id", ""),
            platform=d.get("platform", ""),
            url=d.get("url", ""),
            external_id=d.get("external_id", ""),
            title=d.get("title", ""),
            adapter_name=d.get("adapter_name", ""),
        )

    @classmethod
    def from_stream_message(cls, fields: dict) -> DetailTask:
        """从 Redis Stream 消息解析（字段值为 JSON 字符串）。"""
        import json
        payload = fields.get("payload")
        if payload:
            if isinstance(payload, str):
                data = json.loads(payload)
            else:
                data = payload
            return cls.from_dict(data)
        # 兼容扁平字段
        return cls.from_dict(fields)


@dataclasses.dataclass
class DetailTaskBatch:
    """详情任务批次（对应一条 crawl:detail:queue 消息）。"""

    task_id: str
    source_id: str
    adapter_name: str
    tasks: list[DetailTask] = dataclasses.field(default_factory=list)

    def to_dict(self) -> dict:
        return {
            "task_id": self.task_id,
            "source_id": self.source_id,
            "adapter_name": self.adapter_name,
            "tasks": [t.to_dict() for t in self.tasks],
        }
