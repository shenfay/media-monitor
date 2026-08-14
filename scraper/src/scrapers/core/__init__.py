"""业务编排层。"""
from scrapers.core.registry import register, get_adapter, auto_discover, list_registered
from scrapers.core.executor import TaskExecutor

__all__ = ["register", "get_adapter", "auto_discover", "list_registered", "TaskExecutor"]
