"""TaskEventEmitter — XADD crawl:task:event。

向 Go 侧回传任务状态事件，驱动 TaskRun 状态机流转。
"""
from __future__ import annotations

import json
import logging

import redis

from scrapers import config
from scrapers.contracts.task_event import TaskEvent

logger = logging.getLogger(__name__)


class TaskEventEmitter:
    """发布任务状态事件到 crawl:task:event Stream。"""

    def __init__(
        self,
        redis_url: str | None = None,
        stream: str | None = None,
    ):
        self._redis_url = redis_url or config.settings.redis_url
        self._stream = stream or config.settings.event_stream
        self._conn: redis.Redis | None = None

    @property
    def conn(self) -> redis.Redis:
        if self._conn is None:
            self._conn = redis.Redis.from_url(self._redis_url, decode_responses=True)
        return self._conn

    def emit(self, event: TaskEvent) -> None:
        """发送一条任务事件。"""
        payload = json.dumps(event.to_dict(), ensure_ascii=False)
        self.conn.xadd(self._stream, {"payload": payload})
        logger.debug("Event emitted: task=%s type=%s", event.task_id, event.type)

    # ── 便捷方法 ──────────────────────────────────────────────

    def status_running(self, task_id: str) -> None:
        """任务开始执行。"""
        self.emit(TaskEvent(task_id=task_id, type="status", status="running"))

    def phase_start(self, task_id: str, phase: str) -> None:
        """阶段开始（list / detail）。"""
        self.emit(TaskEvent(task_id=task_id, type="phase_start", phase=phase))

    def phase_done(self, task_id: str, phase: str, count: int, failed: int = 0) -> None:
        """阶段完成。"""
        self.emit(TaskEvent(
            task_id=task_id, type="phase_done", phase=phase,
            total=count, detail_failed=failed,
        ))

    def list_synced(self, task_id: str, count: int) -> None:
        """列表数据已同步。"""
        self.emit(TaskEvent(
            task_id=task_id, type="list_synced", phase="list",
            list_count=count,
        ))

    def detail_progress(self, task_id: str, fetched: int, failed: int) -> None:
        """详情抓取进度更新。"""
        self.emit(TaskEvent(
            task_id=task_id, type="detail_progress", phase="detail",
            detail_count=fetched, detail_failed=failed,
        ))

    def task_done(
        self,
        task_id: str,
        status: str,
        total: int,
        list_count: int,
        detail_count: int,
        detail_failed: int,
    ) -> None:
        """任务完成。"""
        self.emit(TaskEvent(
            task_id=task_id, type="task_done", status=status,
            total=total, list_count=list_count,
            detail_count=detail_count, detail_failed=detail_failed,
        ))

    def task_failed(self, task_id: str, error: str) -> None:
        """任务失败。"""
        self.emit(TaskEvent(
            task_id=task_id, type="task_failed", status="failed",
            error=error,
        ))

    def close(self) -> None:
        if self._conn is not None:
            self._conn.close()
            self._conn = None
