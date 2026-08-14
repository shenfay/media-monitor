"""跨服务数据契约。"""
from scrapers.contracts.source import Source
from scrapers.contracts.article import Article, ArticleBatch, utc_now_iso
from scrapers.contracts.task_event import CrawlTaskMessage, TaskEvent

__all__ = [
    "Source", "Article", "ArticleBatch", "utc_now_iso",
    "CrawlTaskMessage", "TaskEvent",
]
