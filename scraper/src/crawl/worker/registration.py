"""Worker 实例注册：心跳保活、adapter 元数据上报。

Worker 启动时：
1. 将 adapter 元数据写入 Redis（scraper:adapters:{name}），不设 TTL
2. 将自身实例信息写入 Redis（scraper:workers:{id}），TTL 60s
3. 后台线程定期续期 TTL（心跳）

Worker 退出时：
- 删除 worker key，从 ids 集合中移除
"""
from __future__ import annotations

import json
import logging
import socket
import threading
import time
import uuid
from datetime import datetime, timezone

import redis

from crawl import config

logger = logging.getLogger(__name__)

ADAPTER_KEY_PREFIX = "scraper:adapters:"
WORKER_KEY_PREFIX = "scraper:workers:"
WORKER_IDS_KEY = "scraper:worker:ids"


class WorkerRegistration:
    """Worker 实例注册与心跳管理。"""

    def __init__(
        self,
        redis_url: str | None = None,
        worker_id: str | None = None,
        worker_name: str | None = None,
        adapters_metadata: list[dict] | None = None,
        capabilities: list[str] | None = None,
        concurrency: int = 1,
    ):
        self._redis_url = redis_url or config.settings.redis_url
        self._worker_id = worker_id or config.settings.worker_id or str(uuid.uuid4())
        self._worker_name = worker_name or config.settings.worker_name or self._default_name()
        self._adapters_metadata = adapters_metadata or []
        self._capabilities = capabilities or []
        self._concurrency = concurrency
        self._conn: redis.Redis | None = None
        self._heartbeat_thread: threading.Thread | None = None
        self._stop_event = threading.Event()
        self._started_at = datetime.now(timezone.utc).isoformat()
        self._processed_count = 0
        self._current_task = ""

    @property
    def conn(self) -> redis.Redis:
        if self._conn is None:
            self._conn = redis.Redis.from_url(self._redis_url, decode_responses=True)
        return self._conn

    @property
    def worker_id(self) -> str:
        return self._worker_id

    @staticmethod
    def _default_name() -> str:
        try:
            return f"worker-{socket.gethostname()}"
        except Exception:
            return f"worker-{uuid.uuid4().hex[:8]}"

    def register(self) -> None:
        """注册 adapter 元数据 + Worker 实例信息到 Redis。"""
        # 1. 注册 adapter 元数据（不设 TTL，长期保留）
        for meta in self._adapters_metadata:
            key = f"{ADAPTER_KEY_PREFIX}{meta['name']}"
            self.conn.hset(key, mapping={
                "required_tags": json.dumps(meta.get("required_tags", [])),
                "platform_type": meta.get("platform_type", "news"),
                "class_name": meta.get("class_name", ""),
                "first_seen": datetime.now(timezone.utc).isoformat(),
            })
            logger.info("Registered adapter metadata: %s", meta["name"])

        # 2. 注册 Worker 实例信息
        self._write_worker_key()

        # 3. 将 worker_id 加入全局集合
        self.conn.sadd(WORKER_IDS_KEY, self._worker_id)

        logger.info(
            "Worker registered: id=%s name=%s adapters=%s capabilities=%s",
            self._worker_id, self._worker_name,
            [m["name"] for m in self._adapters_metadata],
            self._capabilities,
        )

    def start_heartbeat(self) -> None:
        """启动后台心跳线程。"""
        if self._heartbeat_thread is not None:
            return
        self._stop_event.clear()
        self._heartbeat_thread = threading.Thread(
            target=self._heartbeat_loop,
            daemon=True,
            name="worker-heartbeat",
        )
        self._heartbeat_thread.start()
        logger.info(
            "Heartbeat started: interval=%ds ttl=%ds",
            config.settings.heartbeat_interval,
            config.settings.heartbeat_ttl,
        )

    def deregister(self) -> None:
        """优雅退出：停止心跳、删除 Redis key。"""
        self._stop_event.set()
        try:
            self.conn.delete(f"{WORKER_KEY_PREFIX}{self._worker_id}")
            self.conn.srem(WORKER_IDS_KEY, self._worker_id)
            logger.info("Worker deregistered: %s", self._worker_id)
        except Exception:
            logger.exception("Failed to deregister worker")

    def update_status(self, current_task: str = "", processed_count: int | None = None) -> None:
        """更新 Worker 状态（当前任务、已处理数）。"""
        if current_task:
            self._current_task = current_task
        if processed_count is not None:
            self._processed_count = processed_count
        try:
            updates = {
                "current_task": self._current_task,
                "processed_count": str(self._processed_count),
                "last_heartbeat": str(int(time.time())),
            }
            self.conn.hset(f"{WORKER_KEY_PREFIX}{self._worker_id}", mapping=updates)
        except Exception:
            logger.warning("Failed to update worker status")

    def _write_worker_key(self) -> None:
        """写入 Worker 实例信息到 Redis，设置 TTL。"""
        key = f"{WORKER_KEY_PREFIX}{self._worker_id}"
        self.conn.hset(key, mapping={
            "name": self._worker_name,
            "adapters": json.dumps([m["name"] for m in self._adapters_metadata]),
            "capabilities": json.dumps(self._capabilities),
            "status": "online",
            "concurrency": str(self._concurrency),
            "current_task": self._current_task,
            "processed_count": str(self._processed_count),
            "last_heartbeat": str(int(time.time())),
            "started_at": self._started_at,
        })
        self.conn.expire(key, config.settings.heartbeat_ttl)

    def _heartbeat_loop(self) -> None:
        """定期续期 Worker key 的 TTL。"""
        interval = config.settings.heartbeat_interval
        ttl = config.settings.heartbeat_ttl
        key = f"{WORKER_KEY_PREFIX}{self._worker_id}"

        while not self._stop_event.is_set():
            try:
                self.conn.hset(key, "last_heartbeat", str(int(time.time())))
                self.conn.expire(key, ttl)
            except Exception:
                logger.warning("Heartbeat failed, will retry")
            self._stop_event.wait(interval)
