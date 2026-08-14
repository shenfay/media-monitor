"""适配器注册表 + 自动发现。

适配器子包在 __init__.py 中调用 register() 完成注册。
adapters/__init__.py 调用 auto_discover() 自动导入所有子包。
"""
from __future__ import annotations

import importlib
import logging
import pkgutil
from typing import Type

from scrapers.adapters.base import BaseAdapter

logger = logging.getLogger(__name__)

_REGISTRY: dict[str, Type[BaseAdapter]] = {}


def register(name: str, adapter_cls: Type[BaseAdapter]) -> None:
    """注册适配器。可多次调用，同一 adapter_cls 映射多个名称。"""
    _REGISTRY[name] = adapter_cls
    logger.debug("Registered adapter: %s -> %s", name, adapter_cls.__name__)


def get_adapter(platform_type: str, source_name: str) -> BaseAdapter | None:
    """按 platform_type 或 source_name 查找适配器实例。

    优先精确匹配 source_name，其次匹配 platform_type。
    """
    cls = _REGISTRY.get(source_name) or _REGISTRY.get(platform_type)
    return cls() if cls else None


def auto_discover() -> None:
    """自动导入 scrapers.adapters 下所有子包，触发其 register() 调用。"""
    package_name = "scrapers.adapters"
    try:
        package = importlib.import_module(package_name)
    except ImportError:
        logger.warning("Cannot import %s for auto-discovery", package_name)
        return

    for _, module_name, _ in pkgutil.iter_modules(package.__path__):
        if module_name in ("base",):
            continue
        try:
            importlib.import_module(f"{package_name}.{module_name}")
            logger.debug("Auto-discovered adapter module: %s", module_name)
        except Exception:
            logger.exception("Failed to import adapter module: %s", module_name)


def list_registered() -> dict[str, str]:
    """返回当前注册表 {name: class_name}，用于调试/健康检查。"""
    return {name: cls.__name__ for name, cls in _REGISTRY.items()}
