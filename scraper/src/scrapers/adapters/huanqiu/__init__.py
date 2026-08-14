"""环球网适配器注册。"""
from scrapers.adapters.huanqiu.adapter import HuanqiuAdapter
from scrapers.core.registry import register

register("huanqiu", HuanqiuAdapter)
register("huanqiu_news", HuanqiuAdapter)
register("news", HuanqiuAdapter)
