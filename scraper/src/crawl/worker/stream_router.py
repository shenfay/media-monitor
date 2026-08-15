"""StreamRouter — 根据 Worker capabilities 计算需要消费的 Stream 列表。

Worker 启动时和收到 crawl:config:changed 通知时重新计算。
逻辑：
1. 从 Redis 读取所有活跃 Source 的 tags
2. 对每个 Source，合并 adapter.required_tags ∪ source.tags → 排序 → 计算 stream name
3. 筛选：stream name 包含的所有 tag 都在 worker capabilities 中
4. 去重返回
"""
from __future__ import annotations

import json
import logging

import redis

from crawl import config

logger = logging.getLogger(__name__)

CONFIG_CHANGED_CHANNEL = "crawl:config:changed"


def compute_stream_name(base_stream: str, tags: list[str]) -> str:
    """根据 tags 计算 Stream 名称（与 Go 侧 ComputeStreamName 对齐）。"""
    if not tags:
        return base_stream
    unique = sorted(set(t.strip() for t in tags if t.strip()))
    if not unique:
        return base_stream
    return base_stream + ":" + ":".join(unique)


class StreamRouter:
    """根据 Worker capabilities 计算需要消费的 Stream 列表。"""

    def __init__(
        self,
        redis_conn: redis.Redis,
        base_stream: str | None = None,
        capabilities: list[str] | None = None,
        adapter_required_tags: dict[str, list[str]] | None = None,
    ):
        self._redis = redis_conn
        self._base_stream = base_stream or config.settings.dispatch_stream
        self._capabilities = set(capabilities or [])
        # adapter_name → required_tags 映射
        self._adapter_required_tags = adapter_required_tags or {}

    def compute_streams(self) -> list[str]:
        """计算当前 Worker 需要消费的所有 Stream。

        遍历所有活跃 Source，计算每个 Source 对应的 stream name，
        筛选出当前 Worker 能处理的（所有 tag 都在 capabilities 中）。
        """
        streams: set[str] = set()

        # 始终消费默认流（如果没有特殊 tag）
        if not self._capabilities:
            streams.add(self._base_stream)
            return list(streams)

        # 扫描所有 Source 的 tags
        try:
            source_keys = self._redis.keys("crawl:source:*")  # 备用方案
        except Exception:
            source_keys = []

        # 直接从数据库的 Redis 缓存或从 Source 列表获取
        # 这里使用简化方案：扫描所有可能的 tag 组合
        # 实际生产中 Source 数量有限，直接遍历即可
        all_tags = self._collect_all_source_tags()

        for tags_combo in all_tags:
            stream_name = compute_stream_name(self._base_stream, tags_combo)
            # 检查这个 stream 的所有 tag 是否都在 capabilities 中
            if all(t in self._capabilities for t in tags_combo):
                streams.add(stream_name)

        # 如果没有找到任何匹配的 stream，至少消费默认流
        if not streams:
            streams.add(self._base_stream)

        result = sorted(streams)
        logger.info("Computed streams for worker: %s", result)
        return result

    def _collect_all_source_tags(self) -> list[list[str]]:
        """收集所有 Source 的 tags 组合。

        从 Redis 中的 scraper:adapters:* 和 Source 数据推断。
        简化实现：基于 adapter 的 required_tags 生成可能的组合。
        """
        combos: set[tuple[str, ...]] = set()

        # 空 tags → 默认流
        combos.add(())

        # 单个 adapter 的 required_tags
        for adapter_name, req_tags in self._adapter_required_tags.items():
            if req_tags:
                combos.add(tuple(sorted(req_tags)))

        # 扫描 Redis 中的 Source tags（如果有的话）
        # Source tags 存储在数据库中，Go 侧 dispatch 时已合并
        # 这里从 Redis 的 Source 缓存或直接从已知 tags 推断
        try:
            # 尝试从 Redis 中读取已知的 source tags
            # 实际中 Source 数据在 PostgreSQL，这里用简化方案
            source_tags_raw = self._redis.get("scraper:known_source_tags")
            if source_tags_raw:
                tags_list = json.loads(source_tags_raw)
                for tags in tags_list:
                    combos.add(tuple(sorted(tags)))
        except Exception:
            pass

        return [list(c) for c in combos]

    def update_adapter_tags(self, adapter_name: str, required_tags: list[str]) -> None:
        """更新 adapter 的 required_tags（动态注册用）。"""
        self._adapter_required_tags[adapter_name] = required_tags
