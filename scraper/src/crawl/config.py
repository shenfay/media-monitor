"""配置：环境变量优先，带合理默认值。

所有键以 MM_ 前缀，避免与 Go 侧冲突。
"""
from __future__ import annotations

import os

DEFAULTS = {
    "MM_REDIS_URL": "redis://localhost:6379/0",
    "MM_DISPATCH_STREAM": "crawl:task:dispatch",
    "MM_ARTICLE_STREAM": "crawl:article:ingest",
    "MM_EVENT_STREAM": "crawl:task:event",
    "MM_DETAIL_QUEUE": "crawl:detail:queue",
    "MM_CONSUMER_GROUP": "crawl:worker",
    "MM_CONSUMER_NAME": "scraper-1",
    "MM_HTTP_TIMEOUT": "20",
    "MM_HTTP_RETRIES": "3",
    "MM_BATCH_SIZE": "20",
    "MM_MAX_STREAM_LEN": "500",
    "MM_WORKER_ID": "",
    "MM_WORKER_NAME": "",
    "MM_HEARTBEAT_INTERVAL": "20",
    "MM_HEARTBEAT_TTL": "60",
    "MM_DATABASE_DSN": "",
}


def get(key: str, default=None):
    return os.environ.get(key, DEFAULTS.get(key, default))


class Settings:
    redis_url = get("MM_REDIS_URL")
    dispatch_stream = get("MM_DISPATCH_STREAM")
    article_stream = get("MM_ARTICLE_STREAM")
    event_stream = get("MM_EVENT_STREAM")
    detail_queue = get("MM_DETAIL_QUEUE")
    consumer_group = get("MM_CONSUMER_GROUP")
    consumer_name = get("MM_CONSUMER_NAME")
    http_timeout = float(get("MM_HTTP_TIMEOUT") or 20)
    http_retries = int(get("MM_HTTP_RETRIES") or 3)
    batch_size = int(get("MM_BATCH_SIZE") or 20)
    max_stream_len = int(get("MM_MAX_STREAM_LEN") or 500)
    worker_id = get("MM_WORKER_ID") or ""
    worker_name = get("MM_WORKER_NAME") or ""
    heartbeat_interval = int(get("MM_HEARTBEAT_INTERVAL") or 20)
    heartbeat_ttl = int(get("MM_HEARTBEAT_TTL") or 60)
    database_dsn = get("MM_DATABASE_DSN") or ""


settings = Settings()
