"""Source 适配器抽象基类。

每个数据源（环球网、微博、X…）实现一个子类：
- fetch_list: 拉列表并产出标准 Article（必须）
- fetch_detail: 抓正文（可选，PHASE 2）
- discover_nodes: 返回平台节点路径（新闻站用；社媒可为空）
"""
from __future__ import annotations

from abc import ABC, abstractmethod

from crawl.contracts.article import Article
from crawl.contracts.source import Source


class BaseAdapter(ABC):
    platform_type: str = "news"
    name: str = "base"
    required_tags: list[str] = []  # 适配器声明所需能力标签（如 ["overseas"]）

    @abstractmethod
    def fetch_list(self, source: Source, limit: int = 200) -> list[Article]:
        ...

    def fetch_detail(self, source: Source, article: Article) -> Article:
        # 默认不抓正文
        return article

    def discover_nodes(self, source: Source) -> list:
        return list(source.nodes or [])
