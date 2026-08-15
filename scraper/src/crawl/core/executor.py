"""TaskExecutor — 编排完整任务生命周期。

流程：
1. 查找适配器
2. Phase 1: 列表抓取 → 标准化 → 校验 → 分批回传
3. Phase 2: 详情抓取（可选）→ 分批回传
4. 发送完成/失败事件
"""
from __future__ import annotations

import logging

from crawl.cleaning.normalizer import normalize_article
from crawl.cleaning.validator import validate_article
from crawl.contracts.article import Article
from crawl.contracts.source import Source
from crawl.contracts.task_event import CrawlTaskMessage
from crawl.adapters.registry import get_adapter
from crawl.streams.event_emitter import TaskEventEmitter
from crawl.streams.publisher import ArticlePublisher

logger = logging.getLogger(__name__)


class TaskExecutor:
    """任务编排器：列表抓取 → 详情抓取 → 事件回报。"""

    def __init__(
        self,
        publisher: ArticlePublisher,
        emitter: TaskEventEmitter,
    ):
        self._publisher = publisher
        self._emitter = emitter

    def execute(self, task: CrawlTaskMessage) -> None:
        """执行一次完整抓取任务。"""
        source = task.source
        adapter = get_adapter(source.platform_type, source.name)
        if adapter is None:
            msg = f"No adapter for platform={source.platform_type} name={source.name}"
            logger.error(msg)
            self._emitter.task_failed(task.task_id, msg)
            return

        logger.info(
            "Executing task=%s source=%s adapter=%s",
            task.task_id, source.name, adapter.__class__.__name__,
        )

        self._emitter.status_running(task.task_id)

        # ── Phase 1: 列表抓取 ──────────────────────────────
        self._emitter.phase_start(task.task_id, "list")
        try:
            articles = adapter.fetch_list(source, limit=task.params.get("limit", 200))
        except Exception as e:
            logger.exception("List fetch failed for task=%s", task.task_id)
            self._emitter.task_failed(task.task_id, f"list fetch error: {e}")
            return

        # 标准化 + 校验
        articles = [normalize_article(a) for a in articles]
        valid_articles = []
        for a in articles:
            errs = validate_article(a)
            if errs:
                logger.warning("Article validation failed: url=%s errors=%s", a.url, errs)
                continue
            valid_articles.append(a)

        if not valid_articles:
            logger.warning("No valid articles after validation, task=%s", task.task_id)
            self._emitter.task_done(task.task_id, "success", 0, 0, 0, 0)
            return

        # 回传列表数据
        self._publisher.publish_list(task.task_id, source.id, valid_articles)
        self._emitter.list_synced(task.task_id, len(valid_articles))
        logger.info("List phase done: task=%s count=%d", task.task_id, len(valid_articles))

        # ── Phase 2: 详情抓取（可选）────────────────────────
        detail_count = 0
        detail_failed = 0
        with_body = task.params.get("with_body", False)

        if with_body:
            self._emitter.phase_start(task.task_id, "detail")
            detail_count, detail_failed = self._fetch_details(adapter, source, task, valid_articles)
            self._emitter.phase_done(task.task_id, "detail", detail_count, detail_failed)
            logger.info(
                "Detail phase done: task=%s fetched=%d failed=%d",
                task.task_id, detail_count, detail_failed,
            )

        # ── 完成 ──────────────────────────────────────────
        status = "success"
        if with_body and detail_failed > 0:
            status = "partial" if detail_count > 0 else "failed"

        self._emitter.task_done(
            task.task_id, status,
            total=len(valid_articles),
            list_count=len(valid_articles),
            detail_count=detail_count,
            detail_failed=detail_failed,
        )
        logger.info("Task completed: task=%s status=%s", task.task_id, status)

    def _fetch_details(
        self,
        adapter,
        source: Source,
        task: CrawlTaskMessage,
        articles: list[Article],
    ) -> tuple[int, int]:
        """逐篇抓正文，分批回传。返回 (success_count, failed_count)。"""
        success = 0
        failed = 0
        batch: list[Article] = []
        batch_size = 10  # 每 10 篇回传一批

        for i, article in enumerate(articles):
            try:
                updated = adapter.fetch_detail(source, article)
                updated = normalize_article(updated)
                errs = validate_article(updated)
                if errs:
                    logger.warning("Detail validation failed: url=%s errors=%s", article.url, errs)
                    failed += 1
                else:
                    batch.append(updated)
                    success += 1
            except Exception:
                logger.exception("Detail fetch failed: url=%s", article.url)
                failed += 1

            # 分批回传
            if len(batch) >= batch_size:
                self._publisher.publish_detail(
                    task.task_id, source.id, batch,
                    batch_seq_offset=success - len(batch),
                )
                self._emitter.detail_progress(task.task_id, success, failed)
                batch = []

        # 剩余的回传
        if batch:
            self._publisher.publish_detail(
                task.task_id, source.id, batch,
                batch_seq_offset=success - len(batch),
            )

        return success, failed
