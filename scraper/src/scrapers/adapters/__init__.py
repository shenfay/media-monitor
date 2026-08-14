"""数据源适配器。每个站点一个子包，通过 register() 自动注册。"""
from scrapers.core.registry import auto_discover

auto_discover()
