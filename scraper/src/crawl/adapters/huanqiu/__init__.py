"""环球网适配器注册。"""
from crawl.adapters.huanqiu.adapter import HuanqiuAdapter
from crawl.adapters.huanqiu.history import HuanqiuHistoryAdapter
from crawl.adapters.registry import register

# 常规适配器（近期数据）
register("huanqiu", HuanqiuAdapter)
register("huanqiu_news", HuanqiuAdapter)
register("news", HuanqiuAdapter)

# 历史回刷适配器（深度翻页）
register("huanqiu_history", HuanqiuHistoryAdapter)
