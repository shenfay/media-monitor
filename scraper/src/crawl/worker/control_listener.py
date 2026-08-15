"""ControlListener — 监听 crawl:worker:control Pub/Sub，响应管理指令。

支持指令：
- pause: 暂停消费（处理完当前任务后不再接新任务），心跳继续
- resume: 恢复消费
- shutdown: 优雅停止（处理完当前任务后退出）
- recalculate: 重新计算流列表（配置变更通知）
"""
from __future__ import annotations

import json
import logging
import threading

import redis

from crawl import config

logger = logging.getLogger(__name__)

CONTROL_CHANNEL = "crawl:worker:control"
CONFIG_CHANGED_CHANNEL = "crawl:config:changed"


class ControlListener:
    """监听控制指令和配置变更通知。"""

    def __init__(
        self,
        worker_id: str,
        redis_conn: redis.Redis,
        on_pause: callable = None,
        on_resume: callable = None,
        on_shutdown: callable = None,
        on_recalculate: callable = None,
    ):
        self._worker_id = worker_id
        self._redis = redis_conn
        self._on_pause = on_pause
        self._on_resume = on_resume
        self._on_shutdown = on_shutdown
        self._on_recalculate = on_recalculate
        self._thread: threading.Thread | None = None
        self._stop_event = threading.Event()
        self._pubsub = None

    def start(self) -> None:
        """启动后台监听线程。"""
        if self._thread is not None:
            return
        self._stop_event.clear()
        self._thread = threading.Thread(
            target=self._listen_loop,
            daemon=True,
            name="control-listener",
        )
        self._thread.start()
        logger.info("ControlListener started for worker %s", self._worker_id)

    def stop(self) -> None:
        """停止监听。"""
        self._stop_event.set()
        if self._pubsub:
            try:
                self._pubsub.unsubscribe()
                self._pubsub.close()
            except Exception:
                pass
        logger.info("ControlListener stopped")

    def _listen_loop(self) -> None:
        """持续监听 Pub/Sub 消息。"""
        while not self._stop_event.is_set():
            try:
                self._pubsub = self._redis.pubsub()
                # 订阅控制频道（所有 Worker 都收到，按 worker_id 过滤）
                self._pubsub.subscribe(CONTROL_CHANNEL)
                # 订阅配置变更频道（所有 Worker 都处理）
                self._pubsub.subscribe(CONFIG_CHANGED_CHANNEL)

                for message in self._pubsub.listen():
                    if self._stop_event.is_set():
                        break
                    if message["type"] != "message":
                        continue

                    channel = message["channel"]
                    if isinstance(channel, bytes):
                        channel = channel.decode("utf-8")
                    data = message["data"]
                    if isinstance(data, bytes):
                        data = data.decode("utf-8")

                    try:
                        msg = json.loads(data)
                    except (json.JSONDecodeError, TypeError):
                        continue

                    if channel == CONTROL_CHANNEL:
                        self._handle_control(msg)
                    elif channel == CONFIG_CHANGED_CHANNEL:
                        self._handle_config_change(msg)

            except redis.ConnectionError:
                logger.warning("Control listener: Redis connection lost, reconnecting in 5s...")
                self._stop_event.wait(5)
            except Exception:
                logger.exception("Control listener error")
                self._stop_event.wait(3)

    def _handle_control(self, msg: dict) -> None:
        """处理控制指令。"""
        target_id = msg.get("worker_id", "")
        action = msg.get("action", "")

        # 指令可以是针对特定 Worker，也可以是广播（worker_id 为空或 "*")
        if target_id and target_id != "*" and target_id != self._worker_id:
            return  # 不是发给自己的

        logger.info("Received control command: action=%s", action)

        if action == "pause" and self._on_pause:
            self._on_pause()
        elif action == "resume" and self._on_resume:
            self._on_resume()
        elif action == "shutdown" and self._on_shutdown:
            self._on_shutdown()
        elif action == "recalculate" and self._on_recalculate:
            self._on_recalculate()

    def _handle_config_change(self, msg: dict) -> None:
        """处理配置变更通知。"""
        logger.info("Config changed notification received")
        if self._on_recalculate:
            self._on_recalculate()
