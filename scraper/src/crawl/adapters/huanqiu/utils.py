"""环球网适配器共享工具函数。

包含频道树解析、API URL 构造、日期处理等公共功能。
"""
from __future__ import annotations

import datetime
import logging

from crawl.fetchers.http import get_json

log = logging.getLogger(__name__)

# API 常量
RECOMMEND_URL = "https://m.huanqiu.com/api/index/recommend"
CHANNEL_API_URL = "https://m.huanqiu.com/api/channel"
DEFAULT_HEADERS = {"Referer": "https://m.huanqiu.com/"}
PAGE_SIZE = 25

# 类级别缓存：频道树变化很少，进程内复用
_channel_cache: dict | None = None
_channel_map: dict[str, str] | None = None  # channel_name → node_key

# /api/channel 树不完整的硬编码回退映射（channel_name → node_id）
# 这些频道在导航中存在，但 /api/channel 未收录或仅为子节点
_FALLBACK_NODE_MAP: dict[str, str] = {
    "fjxzc": "/e3pmh1nnq/f3pt9mit9",    # 新征程（不在频道树中）
    "editorial": "/e3pmub6h5/e3prafm0g",  # 社评（opinion 子节点）
    "shanrenping": "/e3pmub6h5/e3prcgifj",  # 单仁平（opinion 子节点，已停更）
}


def iso_from_ms(ms) -> str:
    """将毫秒时间戳转换为 ISO 8601 格式字符串。"""
    try:
        return datetime.datetime.fromtimestamp(int(ms) / 1000, datetime.timezone.utc).isoformat()
    except Exception:
        return ""


def within_months(ms, months: int) -> bool:
    """判断毫秒时间戳是否在最近 N 个月内。"""
    if not ms:
        return False
    cutoff = datetime.datetime.now(datetime.timezone.utc).timestamp() - months * 30 * 86400
    return (int(ms) / 1000) >= cutoff


def _collect_child_node_ids(node: dict) -> list[str]:
    """递归收集节点及其所有子节点的 node ID。"""
    ids = []
    node_id = node.get("node", "")
    if node_id:
        ids.append(node_id)
    for child in (node.get("children") or {}).values():
        ids.extend(_collect_child_node_ids(child))
    return ids


def _fetch_channel_tree() -> dict:
    """获取并缓存 /api/channel 频道树。"""
    global _channel_cache
    if _channel_cache is not None:
        return _channel_cache
    try:
        data = get_json(CHANNEL_API_URL, headers=DEFAULT_HEADERS)
        _channel_cache = data
        return data
    except Exception:
        log.warning("Failed to fetch channel tree from %s", CHANNEL_API_URL, exc_info=True)
        return {}


def _build_channel_map() -> dict[str, str]:
    """构建 channel_name → node_key 映射。

    遍历完整频道树（含子节点），匹配规则：
    - 域名前缀：world.huanqiu.com → channel_name="world"
    - 列表路径：m.huanqiu.com/list/english → channel_name="english"
    - URL 尾部路径：china.huanqiu.com/lh → channel_name="lh"
    """
    global _channel_map
    if _channel_map is not None:
        return _channel_map

    tree = _fetch_channel_tree()
    mapping: dict[str, str] = {}

    def _walk(node: dict):
        url = node.get("url", "")
        key = node.get("node", "")
        if url and key:
            # 域名前缀匹配：world.huanqiu.com → "world"
            if "." in url and "huanqiu.com" in url and not url.startswith("m."):
                domain = url.split(".")[0]
                # 仅当 URL 无额外路径时才作为顶层频道映射
                path_after_domain = url.split("huanqiu.com")[-1].strip("/")
                if not path_after_domain:
                    mapping[domain] = key
                else:
                    # 子路径频道：china.huanqiu.com/lh → "lh"
                    ch = path_after_domain.split("/")[0]
                    if ch:
                        mapping[ch] = key
            # 列表路径匹配：m.huanqiu.com/list/english → "english"
            if "m.huanqiu.com/list/" in url:
                ch = url.split("/list/")[-1].strip("/")
                if ch:
                    mapping[ch] = key
        for child in (node.get("children") or {}).values():
            _walk(child)

    for child in (tree.get("children") or {}).values():
        _walk(child)

    _channel_map = mapping
    return mapping


def resolve_channel_node(channel_name: str) -> list[str] | None:
    """根据频道名解析所有子节点 ID 列表。

    优先从频道树查找，回退到硬编码映射。
    """
    mapping = _build_channel_map()
    node_key = mapping.get(channel_name)

    if node_key:
        # 从树中找到节点，递归收集所有子节点
        tree = _fetch_channel_tree()
        target_node = _find_node_by_key(tree, node_key)
        if target_node:
            node_ids = _collect_child_node_ids(target_node)
            if node_ids:
                return node_ids

    # 回退：硬编码映射（单个节点，无子节点）
    fallback = _FALLBACK_NODE_MAP.get(channel_name)
    if fallback:
        return [fallback]

    return None


def _find_node_by_key(root: dict, target_key: str) -> dict | None:
    """在频道树中按 node key 查找节点。"""
    if root.get("node") == target_key:
        return root
    for child in (root.get("children") or {}).values():
        result = _find_node_by_key(child, target_key)
        if result:
            return result
    return None


def build_api_url(list_endpoint: str) -> str:
    """根据 list_endpoint 构造实际可用的 API URL。

    优先使用 node 参数（新 API），失败时回退到 channel 参数（旧 API）。
    """
    if not list_endpoint:
        return RECOMMEND_URL

    channel_name = ""
    if list_endpoint.startswith("/list/"):
        channel_name = list_endpoint[len("/list/"):].strip("/")

    if not channel_name:
        return RECOMMEND_URL

    # 优先：用 node 参数（新 API，返回最新数据）
    node_ids = resolve_channel_node(channel_name)
    if node_ids:
        quoted = ",".join(f'"{nid}"' for nid in node_ids)
        return f"https://m.huanqiu.com/api/list?node={quoted}"

    # 回退：旧 channel 参数（可能返回陈旧数据）
    log.warning("Channel '%s' not found in tree, falling back to channel= param", channel_name)
    return f"https://m.huanqiu.com/api/list?channel={channel_name}"


def item_date_range(items: list[dict]) -> tuple[str, str]:
    """从一批 items 中提取最早和最新的发布时间。"""
    timestamps = []
    for item in items:
        if not item:
            continue
        ms = int(item.get("ctime") or item.get("xtime") or item.get("ext_displaytime") or 0)
        if ms:
            timestamps.append(ms)
    if not timestamps:
        return ("N/A", "N/A")
    oldest = min(timestamps)
    newest = max(timestamps)
    return (iso_from_ms(oldest), iso_from_ms(newest))
