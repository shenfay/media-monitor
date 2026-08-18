"""环球网历史数据回刷适配器（仅 CLI 模式）。

与常规 HuanqiuAdapter 的区别：
- 无文章数量上限，持续翻页直到 API 返回空或数据超出 since 时间窗口
- 每次 API 请求都打印详细信息（URL、返回条数、日期范围）
- 请求间加入延迟，避免触发反爬
- 支持 since 精确日期参数（从 source.extra["since"] 读取，ISO 8601 格式）
- 支持逐子节点查询：每个子节点独立翻页，各自享有 offset=10000 上限
- 仅通过 CLI `history --since 2026-01-01` 使用，不走任务调度系统

用法示例：
    python -m crawl --source source_hq_world history --since 2026-01-01
    python -m crawl --source source_hq_world history --since 2026-01-01 --per-node
    python -m crawl --source source_hq_world history --since 2026-01-01 --delay 0.5
"""
from __future__ import annotations

import datetime
import logging
import time

from crawl.adapters.base import BaseAdapter
from crawl.adapters.huanqiu.adapter import HuanqiuAdapter
from crawl.adapters.huanqiu.utils import (
    DEFAULT_HEADERS,
    PAGE_SIZE,
    build_api_url,
    iso_from_ms,
    item_date_range,
    resolve_channel_node,
)
from crawl.contracts.article import Article
from crawl.contracts.source import Source
from crawl.fetchers.http import get_json

log = logging.getLogger(__name__)


def _parse_since(since: str) -> datetime.datetime | None:
    """解析 since 参数为 UTC datetime。支持 ISO 8601 日期或日期时间。"""
    if not since:
        return None
    for fmt in ("%Y-%m-%d", "%Y-%m-%dT%H:%M:%S", "%Y-%m-%dT%H:%M:%S%z", "%Y-%m-%d %H:%M:%S"):
        try:
            dt = datetime.datetime.strptime(since, fmt)
            if dt.tzinfo is None:
                dt = dt.replace(tzinfo=datetime.timezone.utc)
            return dt
        except ValueError:
            continue
    log.warning("Cannot parse since=%r, ignoring", since)
    return None


class HuanqiuHistoryAdapter(BaseAdapter):
    """环球网历史数据回刷适配器。

    与常规适配器的核心区别：
    - 不设文章数量上限（忽略 limit 参数）
    - 使用 since 精确时间过滤（而非 months 近似计算）
    - 详细的请求日志，方便调试和监控
    - 请求间加延迟
    """

    platform_type = "news"
    name = "huanqiu_history"
    required_tags = ["huanqiu"]

    def fetch_detail(self, source: Source, article: Article) -> Article:
        """正文抓取委托给常规适配器（逻辑相同）。"""
        return HuanqiuAdapter().fetch_detail(source, article)

    def fetch_list(self, source: Source, limit: int = 200, on_page=None, start_offset: int = 0) -> list[Article]:
        """翻页抓取。

        Args:
            on_page: 可选回调 ``callback(page_articles: list[Article], page: int)``，
                     每页处理完后立即调用，方便上层增量发布。
            start_offset: 起始 offset，用于断点续传（仅非 per-node 模式）。
        """
        # limit 参数在此适配器中忽略，一直翻页直到 API 返回空或超出时间窗口
        source_filter = source.source_filter or ""
        since_str = (source.extra or {}).get("since", "")
        until_str = (source.extra or {}).get("until", "")
        delay = float((source.extra or {}).get("delay", 1.0))
        per_node = (source.extra or {}).get("per_node", False)
        since_dt = _parse_since(since_str)
        until_dt = _parse_since(until_str)

        if since_dt:
            since_cutoff_ts = since_dt.timestamp()
        else:
            since_cutoff_ts = 0

        if until_dt:
            until_cutoff_ts = until_dt.timestamp()
        else:
            until_cutoff_ts = float('inf')

        log.info(
            "=== 历史回刷开始 | since=%s | until=%s | source_filter=%r | delay=%.1fs | per_node=%s ===",
            since_str or "无", until_str or "无", source_filter, delay, per_node,
        )

        if per_node:
            return self._fetch_per_node(
                source, source_filter, since_str, since_cutoff_ts,
                until_cutoff_ts, delay, on_page,
            )

        api_url = build_api_url(source.list_endpoint)
        log.info("API base URL: %s", api_url)

        articles, _ = self._paginate_url(
            api_url, source, source_filter, since_str, since_cutoff_ts,
            until_cutoff_ts, delay, on_page, start_offset,
        )
        return articles

    def _fetch_per_node(
        self, source, source_filter, since_str, since_cutoff_ts,
        until_cutoff_ts, delay, on_page,
    ) -> list[Article]:
        """逐子节点查询：每个子节点独立翻页，各自享有 offset=10000 上限。"""
        # 解析频道名 → 子节点 ID 列表
        channel_name = ""
        if source.list_endpoint.startswith("/list/"):
            channel_name = source.list_endpoint[len("/list/"):].strip("/")

        node_ids = resolve_channel_node(channel_name) if channel_name else None
        if not node_ids:
            log.warning("无法解析子节点，回退到合并查询")
            api_url = build_api_url(source.list_endpoint)
            articles, _ = self._paginate_url(
                api_url, source, source_filter, since_str, since_cutoff_ts,
                until_cutoff_ts, delay, on_page, 0,
            )
            return articles

        log.info("逐子节点查询: %d 个子节点", len(node_ids))
        all_articles: list[Article] = []
        seen_aids: set[str] = set()
        total_pages = 0

        for node_idx, node_id in enumerate(node_ids):
            node_url = f'https://m.huanqiu.com/api/list?node="{node_id}"'
            log.info(
                "── 子节点 [%d/%d] %s ──",
                node_idx + 1, len(node_ids), node_id,
            )

            node_articles, node_pages = self._paginate_url(
                node_url, source, source_filter, since_str, since_cutoff_ts,
                until_cutoff_ts, delay, on_page, 0,
                page_prefix=f"N{node_idx+1}",
            )

            # 去重：跨子节点可能重复
            new_articles = []
            for a in node_articles:
                if a.external_id not in seen_aids:
                    seen_aids.add(a.external_id)
                    new_articles.append(a)

            if new_articles:
                all_articles.extend(new_articles)
                log.info(
                    "  ↳ 子节点 %s 新增 %d 条（去重后），累计 %d 条",
                    node_id, len(new_articles), len(all_articles),
                )
            total_pages += node_pages

        log.info("=" * 60)
        log.info("=== 逐子节点回刷完成 ===")
        log.info("  子节点数: %d", len(node_ids))
        log.info("  总页数:   %d", total_pages)
        log.info("  总条数:   %d（去重后）", len(all_articles))
        if all_articles:
            dates = [a.published_at for a in all_articles if a.published_at]
            if dates:
                log.info("  日期范围: %s ~ %s", min(dates), max(dates))
        log.info("=" * 60)
        return all_articles

    def _paginate_url(
        self, api_url, source, source_filter, since_str, since_cutoff_ts,
        until_cutoff_ts, delay, on_page, start_offset=0,
        page_prefix="",
    ) -> tuple[list[Article], int]:
        """对单个 API URL 进行翻页抓取。返回 (articles, page_count)。"""
        prefix = f"[{page_prefix}]" if page_prefix else ""
        articles: list[Article] = []
        offset = start_offset
        page = start_offset // PAGE_SIZE if start_offset else 0
        empty_pages = 0
        max_empty_pages = 3

        while True:
            page += 1
            sep = "&" if "?" in api_url else "?"
            url = f"{api_url}{sep}offset={offset}&limit={PAGE_SIZE}"

            log.info(
                "%s[Page %d] GET %s  (offset=%d, collected=%d)",
                prefix, page, url, offset, len(articles),
            )

            # ── 发起请求（带重试） ──
            max_retries = 3
            data = None
            for attempt in range(max_retries):
                try:
                    data = get_json(url, headers=DEFAULT_HEADERS)
                    break
                except Exception as e:
                    err_str = str(e)
                    if "404" in err_str or "Not Found" in err_str:
                        log.warning("%s[Page %d] HTTP 404: API offset 上限，停止翻页", prefix, page)
                        data = None
                        break
                    if attempt < max_retries - 1:
                        log.warning("%s[Page %d] 请求失败 (尝试 %d/%d): %s", prefix, page, attempt + 1, max_retries, e)
                        time.sleep(delay * 2)
                    else:
                        log.warning("%s[Page %d] 请求失败，已重试 %d 次: %s", prefix, page, max_retries, e)
                        data = None

            if data is None:
                break

            items = (data.get("data") or {}).get("list") or data.get("list") or []

            if not items:
                empty_pages += 1
                log.info("%s[Page %d] 返回空列表 (%d/%d)", prefix, page, empty_pages, max_empty_pages)
                if empty_pages >= max_empty_pages:
                    log.info("%s连续 %d 页为空，API 数据已耗尽，停止翻页", prefix, max_empty_pages)
                    break
                offset += PAGE_SIZE
                time.sleep(delay)
                continue

            empty_pages = 0

            oldest_date, newest_date = item_date_range(items)
            log.info(
                "%s[Page %d] 返回 %d 条 | 日期范围: %s ~ %s",
                prefix, page, len(items), oldest_date, newest_date,
            )

            page_added = 0
            page_articles: list[Article] = []
            reached_cutoff = False
            for item in items:
                if not item or not item.get("aid"):
                    continue

                ctime = int(item.get("ctime") or item.get("xtime") or item.get("ext_displaytime") or 0)

                if ctime:
                    item_ts = ctime / 1000
                    if since_cutoff_ts and item_ts < since_cutoff_ts:
                        log.info(
                            "%s  └─ 文章 %s 发布时间 %s 早于 since %s，后续数据更早，停止",
                            prefix, item.get("aid"), iso_from_ms(ctime), since_str,
                        )
                        reached_cutoff = True
                        break
                    if until_cutoff_ts != float('inf') and item_ts > until_cutoff_ts:
                        continue

                source_name = (item.get("source") or {}).get("name", "")
                if source_filter and source_name and source_name != source_filter:
                    continue

                aid = item.get("aid")
                host = item.get("host") or "m.huanqiu.com"
                cover = item.get("cover") or ""
                if cover and cover.startswith("//"):
                    cover = "https:" + cover

                page_articles.append(
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
                page_added += 1

            articles.extend(page_articles)

            log.info(
                "%s[Page %d] 本页新增 %d 条 | 累计 %d 条",
                prefix, page, page_added, len(articles),
            )

            if page % 10 == 0:
                log.info(
                    "%s── 进度: %d 页, %d 条, 最新=%s, 最早=%s ──",
                    prefix, page, len(articles), newest_date, oldest_date,
                )

            if on_page and page_articles:
                try:
                    on_page(page_articles, page)
                except Exception:
                    log.exception("on_page 回调失败")

            if reached_cutoff:
                log.info("%s已到达 since=%s 时间边界，停止翻页", prefix, since_str)
                break

            offset += PAGE_SIZE
            time.sleep(delay)

        return articles, page
