"""环球网适配器注册。"""
from crawl.adapters.huanqiu.adapter import HuanqiuAdapter
from crawl.adapters.registry import register

register("huanqiu", HuanqiuAdapter)
register("huanqiu_news", HuanqiuAdapter)
register("news", HuanqiuAdapter)
