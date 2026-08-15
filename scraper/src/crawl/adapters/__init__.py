"""数据源适配器。每个站点一个子包，通过 register() 自动注册。"""
from crawl.adapters.registry import (
    auto_discover,
    get_adapter,
    get_all_required_tags,
    list_metadata,
    list_registered,
    register,
)

auto_discover()

__all__ = [
    "register",
    "get_adapter",
    "auto_discover",
    "list_registered",
    "list_metadata",
    "get_all_required_tags",
]
