"""环球网适配器（PHASE 1：列表 + PHASE 2：正文）。

复用已逆向验证的 m.huanqiu.com/api/list：
- node 参数 = 逗号拼接的引号字符串（如 "a","b"，注意**不带方括号**），再 URL 编码
- 每条记录含 source.name（来源）/ ctime（毫秒时间戳）/ host + aid（拼 URL）
过滤：来源名精确匹配 source_filter（默认"环球网"），ctime 落在近 months 月内。
PHASE 2 正文：fetch_detail 抓文章页 <textarea class="article-content">。

注意：站点可能限流/反爬（尤其非常规出口 IP），生产请在正常网络环境运行。
"""
from __future__ import annotations

import datetime
import html
import json
import re
import urllib.parse
from pathlib import Path

from scrapers.adapters.base import BaseAdapter
from scrapers.contracts.article import Article
from scrapers.contracts.source import Source
from scrapers.fetchers.http import get_json, get_text_auto

NAV_URL = "https://m.huanqiu.com/api/nav"
LIST_URL = "https://m.huanqiu.com/api/list"
MOBILE_UA = {"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15"}


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


class HuanqiuAdapter(BaseAdapter):
    platform_type = "news"
    name = "huanqiu"

    # 兜底节点清单：站点改版后 nav 不再暴露节点 ID，这里放已验证的主频道节点。
    # 全站覆盖请维护同目录下的 nodes.json（JSON 数组），适配器会自动读取。
    CURATED_NODES = ["e3pmh22ph/e3pmh2398"]  # world 频道
    NODES_FILE = Path(__file__).with_name("nodes.json")

    # ---- 节点发现：nav -> 维护清单 -> 兜底 ----
    def discover_nodes(self, source: Source) -> list:
        if source.nodes:
            return list(source.nodes)
        # 1) nav 接口（站点改版后可能不再含 /e3 节点路径）
        try:
            nav = get_json(NAV_URL)
            found = self._collect_nodes(nav)
            if found:
                return found
        except Exception:
            pass
        # 2) 维护好的节点清单（用户落盘，覆盖全站）
        try:
            if self.NODES_FILE.exists():
                return json.loads(self.NODES_FILE.read_text(encoding="utf-8"))
        except Exception:
            pass
        # 3) 兜底
        return list(self.CURATED_NODES)

    @staticmethod
    def _collect_nodes(obj) -> list:
        found: list[str] = []

        def walk(val):
            if isinstance(val, dict):
                node = val.get("node")
                if isinstance(node, str) and node.startswith("/e3"):
                    found.append(node)
                for v in val.values():
                    walk(v)
            elif isinstance(val, list):
                for v in val:
                    walk(v)

        walk(obj)
        seen = set()
        out = []
        for node in found:
            if node not in seen:
                seen.add(node)
                out.append(node)
        return out

    # ---- 列表抓取 ----
    def fetch_list(self, source: Source, limit: int = 200) -> list[Article]:
        nodes = self.discover_nodes(source)
        if not nodes:
            return []
        source_filter = source.source_filter or "环球网"
        months = source.months or 6

        articles: list[Article] = []
        offset = 0
        page_size = 20
        max_iterations = 0
        headers = {"Referer": "https://m.huanqiu.com/"}
        # node 参数格式：逗号拼接的引号字符串，不带方括号（站点约定）
        node_str = ",".join(f'"{n}"' for n in nodes)
        node_param = urllib.parse.quote(node_str)
        while len(articles) < limit and max_iterations < 5000:
            max_iterations += 1
            url = f"{LIST_URL}?node={node_param}&offset={offset}&limit={page_size}"
            try:
                data = get_json(url, headers=headers)
            except Exception:
                break
            items = (data.get("data") or {}).get("list") or data.get("list") or []
            if not items:
                break
            for item in items:
                if not item or not item.get("aid"):
                    continue
                source_name = (item.get("source") or {}).get("name", "")
                # 仅当来源名非空且与过滤值不符时才跳过；空来源无法判定，保留交由下游
                if source_filter and source_name and source_name != source_filter:
                    continue
                ctime = int(item.get("ctime") or item.get("xtime") or 0)
                if not _within_months(ctime, months):
                    continue
                aid = item.get("aid")
                host = item.get("host") or "m.huanqiu.com"
                if not aid:
                    continue
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
                        images=[item["cover"]] if item.get("cover") else [],
                    )
                )
            offset += page_size
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
