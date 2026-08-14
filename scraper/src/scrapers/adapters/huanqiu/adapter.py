"""环球网适配器（PHASE 1：列表 + PHASE 2：正文）。

移动端 API：
- 列表：m.huanqiu.com/api/index/recommend?offset=0&limit=20
  返回全站推荐信息流，每页约 6 条，可通过 offset 翻页。
- 导航：m.huanqiu.com/api/nav → 获取全部频道（ename）。
- 详情页：抓取 <textarea class="article-content"> 获取正文 HTML。

Source 配置说明：
- list_endpoint: 列表 API 路径（默认 /api/index/recommend）
- source_filter: 来源名精确过滤（如 "环球网"），空则不过滤
- months: 时间窗口（近 N 月），默认 6

注意：站点可能限流/反爬（尤其非常规出口 IP），生产请在正常网络环境运行。
"""
from __future__ import annotations

import datetime
import html
import re

from scrapers.adapters.base import BaseAdapter
from scrapers.contracts.article import Article
from scrapers.contracts.source import Source
from scrapers.fetchers.http import get_json, get_text_auto

RECOMMEND_URL = "https://m.huanqiu.com/api/index/recommend"
NAV_URL = "https://m.huanqiu.com/api/nav"
MOBILE_UA = {"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15"}
DEFAULT_HEADERS = {"Referer": "https://m.huanqiu.com/"}
PAGE_SIZE = 20


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

    # ---- 节点发现（保留兼容，新逻辑不再依赖节点）----
    def discover_nodes(self, source: Source) -> list:
        """已废弃：旧 node 参数已失效，保留接口兼容。"""
        return list(source.nodes or [])

    # ---- 列表抓取（PHASE 1）----
    def fetch_list(self, source: Source, limit: int = 200) -> list[Article]:
        source_filter = source.source_filter or ""
        months = source.months or 6
        api_url = f"https://m.huanqiu.com{source.list_endpoint}" if source.list_endpoint else RECOMMEND_URL

        articles: list[Article] = []
        offset = 0
        max_iterations = 0

        while len(articles) < limit and max_iterations < 5000:
            max_iterations += 1
            url = f"{api_url}?offset={offset}&limit={PAGE_SIZE}"
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
