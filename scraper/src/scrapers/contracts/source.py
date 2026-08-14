"""Source 数据源配置契约。

Go 通过 crawl:task:dispatch 消息内嵌下发（含明文 auth）。
"""
from __future__ import annotations

import dataclasses


@dataclasses.dataclass
class Source:
    """数据源配置。"""

    id: str
    name: str
    platform_type: str = "news"           # news | social | social_overseas
    base_url: str = ""
    list_endpoint: str = ""
    nodes: list[str] = dataclasses.field(default_factory=list)
    source_filter: str = ""               # 来源名精确过滤，如 "环球网"
    months: int = 6                       # 时间窗（近 N 月）
    schedule: str = ""                    # cron 表达式
    auth: dict = dataclasses.field(default_factory=dict)  # 社媒登录态/cookie/token
    enabled: bool = True
    extra: dict = dataclasses.field(default_factory=dict)

    @classmethod
    def from_dict(cls, d: dict) -> Source:
        """从 Go 下发的 JSON 字典构造。"""
        return cls(
            id=str(d.get("id") or d.get("source_id") or ""),
            name=d.get("name", ""),
            platform_type=d.get("platform_type") or d.get("platform") or "news",
            base_url=d.get("base_url", ""),
            list_endpoint=d.get("list_endpoint", ""),
            nodes=d.get("nodes") or [],
            source_filter=d.get("source_filter", ""),
            months=int(d.get("months") or 6),
            schedule=d.get("schedule", ""),
            auth=d.get("auth") or {},
            enabled=bool(d.get("enabled", True)),
            extra=d.get("extra") or {},
        )

    def to_dict(self) -> dict:
        return {
            "id": self.id,
            "name": self.name,
            "platform_type": self.platform_type,
            "base_url": self.base_url,
            "list_endpoint": self.list_endpoint,
            "nodes": self.nodes,
            "source_filter": self.source_filter,
            "months": self.months,
            "schedule": self.schedule,
            "auth": self.auth,
            "enabled": self.enabled,
            "extra": self.extra,
        }
