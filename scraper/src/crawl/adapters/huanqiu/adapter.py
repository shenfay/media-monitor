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

import html
import logging
import re

from crawl.adapters.base import BaseAdapter
from crawl.adapters.huanqiu.utils import (
    DEFAULT_HEADERS,
    PAGE_SIZE,
    build_api_url,
    iso_from_ms,
    resolve_channel_node,
    within_months,
)
from crawl.contracts.article import Article
from crawl.contracts.source import Source
from crawl.fetchers.http import get_json, get_text_auto

log = logging.getLogger(__name__)

MOBILE_UA = {"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15"}


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
            return resolve_channel_node(ch) or []
        return []

    # ---- 列表抓取（PHASE 1）----
    def fetch_list(self, source: Source, limit: int = 200) -> list[Article]:
        source_filter = source.source_filter or ""
        months = source.months or 6
        api_url = build_api_url(source.list_endpoint)

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
                if not within_months(ctime, months):
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
                        published_at=iso_from_ms(ctime),
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

        # 辅助：从 <textarea class="xxx"> 提取值
        def _textarea(cls: str) -> str:
            m = re.search(
                rf'<textarea[^>]*class="{cls}"[^>]*>(.*?)</textarea>',
                page, re.S,
            )
            return m.group(1).strip() if m else ""

        # 正文
        content = _textarea("article-content")
        if content:
            article.content = html.unescape(content)

        # 元数据（详情页数据更权威，覆盖列表阶段）
        title = _textarea("article-title")
        if title:
            article.title = title

        subtitle = _textarea("article-subtitle")
        if subtitle:
            article.subtitle = subtitle

        author = _textarea("article-author")
        if author:
            # 清理 "作者：" / "作者:" 前缀
            author = re.sub(r'^作者[：:]\s*', '', author).strip()
            article.author = author

        source_name = _textarea("article-source-name")
        if source_name:
            article.source_name = source_name

        time_str = _textarea("article-time")
        if time_str and time_str.isdigit():
            article.published_at = iso_from_ms(int(time_str))

        cover = _textarea("article-cover")
        if cover:
            if cover.startswith("//"):
                cover = "https:" + cover
            article.images = [cover]

        return article
