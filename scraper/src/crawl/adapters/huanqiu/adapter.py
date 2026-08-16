"""环球网适配器（PHASE 1：列表 + PHASE 2：正文）。

移动端 API：
- 全站推荐：m.huanqiu.com/api/index/recommend?offset=0&limit=20
  返回全站推荐信息流，可通过 offset 翻页。
- 频道列表：m.huanqiu.com/api/list?node={nodeIds}&offset=0&limit=25
  按节点 ID 列表获取文章列表（node 参数从 /api/channel 动态解析）。
- 频道树：m.huanqiu.com/api/channel → 获取全部频道及其子节点 ID。
- 详情页：抓取 <textarea class="article-content"> 获取正文 HTML。

Source 配置说明：
- list_endpoint: 列表路径
  - 空 → 全站推荐流（/api/index/recommend）
  - /list/{channel} → 指定频道（自动从频道树解析 node ID）
- source_filter: 来源名精确过滤（如 "环球网"），空则不过滤
- months: 时间窗口（近 N 月），默认 6

注意：站点可能限流/反爬（尤其非常规出口 IP），生产请在正常网络环境运行。
"""
from __future__ import annotations

import datetime
import html
import logging
import re

from crawl.adapters.base import BaseAdapter
from crawl.contracts.article import Article
from crawl.contracts.source import Source
from crawl.fetchers.http import get_json, get_text_auto

log = logging.getLogger(__name__)

RECOMMEND_URL = "https://m.huanqiu.com/api/index/recommend"
CHANNEL_API_URL = "https://m.huanqiu.com/api/channel"
MOBILE_UA = {"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15"}
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


def _iso_from_ms(ms) -> str:
    try:
        return datetime.datetime.fromtimestamp(int(ms) / 1000, datetime.timezone.utc).isoformat()
    except Exception:
        return ""


def _within_months(ms, months: int) -> bool:
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


def _resolve_channel_node(channel_name: str) -> list[str] | None:
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


def _build_api_url(list_endpoint: str) -> str:
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
    node_ids = _resolve_channel_node(channel_name)
    if node_ids:
        quoted = ",".join(f'"{nid}"' for nid in node_ids)
        return f"https://m.huanqiu.com/api/list?node={quoted}"

    # 回退：旧 channel 参数（可能返回陈旧数据）
    log.warning("Channel '%s' not found in tree, falling back to channel= param", channel_name)
    return f"https://m.huanqiu.com/api/list?channel={channel_name}"


class HuanqiuAdapter(BaseAdapter):
    platform_type = "news"
    name = "huanqiu"
    required_tags = ["huanqiu"]

    # ---- 节点发现 ----
    def discover_nodes(self, source: Source) -> list:
        """返回频道对应的节点 ID 列表（从频道树动态解析）。"""
        if source.nodes:
            return list(source.nodes)
        if source.list_endpoint and source.list_endpoint.startswith("/list/"):
            ch = source.list_endpoint[len("/list/"):].strip("/")
            return _resolve_channel_node(ch) or []
        return []

    # ---- 列表抓取（PHASE 1）----
    def fetch_list(self, source: Source, limit: int = 200) -> list[Article]:
        source_filter = source.source_filter or ""
        months = source.months or 6
        api_url = _build_api_url(source.list_endpoint)

        articles: list[Article] = []
        offset = 0
        max_iterations = 0

        while len(articles) < limit and max_iterations < 5000:
            max_iterations += 1
            sep = "&" if "?" in api_url else "?"
            url = f"{api_url}{sep}offset={offset}&limit={PAGE_SIZE}"
            try:
                data = get_json(url, headers=DEFAULT_HEADERS)
            except Exception:
                break
            items = (data.get("data") or {}).get("list") or data.get("list") or []
            if not items:
                break

            for item in items:
                if not item or not item.get("aid"):
                    continue
                source_name = (item.get("source") or {}).get("name", "")
                # 来源过滤：仅当过滤值非空且来源名不匹配时跳过
                if source_filter and source_name and source_name != source_filter:
                    continue
                ctime = int(item.get("ctime") or item.get("xtime") or item.get("ext_displaytime") or 0)
                if not _within_months(ctime, months):
                    continue
                aid = item.get("aid")
                host = item.get("host") or "m.huanqiu.com"
                cover = item.get("cover") or ""
                if cover and cover.startswith("//"):
                    cover = "https:" + cover
                articles.append(
                    Article(
                        source_id=source.id,
                        platform="news",
                        title=item.get("title", ""),
                        url=f"https://{host}/article/{aid}",
                        external_id=str(aid),
                        author=item.get("author", "") or "",
                        source_name=source_name,
                        summary=item.get("summary", "") or "",
                        published_at=_iso_from_ms(ctime),
                        images=[cover] if cover else [],
                    )
                )
            offset += PAGE_SIZE
        return articles[:limit]

    # ---- 正文抓取（PHASE 2）----
    def fetch_detail(self, source: Source, article: Article) -> Article:
        try:
            page = get_text_auto(article.url, headers=MOBILE_UA)
        except Exception:
            return article
        match = re.search(
            r'<textarea[^>]*class="article-content"[^>]*>(.*?)</textarea>', page, re.S
        )
        if match:
            article.content = html.unescape(match.group(1))
        return article
