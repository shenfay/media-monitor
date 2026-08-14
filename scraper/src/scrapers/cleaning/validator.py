"""文章校验：返回错误列表（空列表 = 通过）。"""
from __future__ import annotations

from scrapers.contracts.article import Article


def validate_article(article: Article) -> list[str]:
    """校验文章必填字段，返回错误描述列表。

    空列表表示校验通过。
    """
    errors: list[str] = []

    if not article.source_id:
        errors.append("source_id is required")
    if not article.url:
        errors.append("url is required")
    if not article.title:
        errors.append("title is required")
    if not article.platform:
        errors.append("platform is required")

    # URL 格式基本校验
    if article.url and not article.url.startswith(("http://", "https://")):
        errors.append(f"url must start with http(s)://, got: {article.url[:50]}")

    return errors
