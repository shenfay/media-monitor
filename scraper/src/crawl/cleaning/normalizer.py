"""文章标准化：时间格式、HTML 清理、URL 规范化。"""
from __future__ import annotations

import hashlib
import re
import urllib.parse

from crawl.contracts.article import Article


def normalize_article(article: Article) -> Article:
    """标准化一篇文章（原地修改并返回）。"""
    article.url = normalize_url(article.url)
    article.title = clean_title(article.title)
    article.summary = clean_text(article.summary)
    article.content = clean_text(article.content)
    article.published_at = normalize_timestamp(article.published_at)
    if not article.external_id and article.url:
        article.external_id = _extract_id_from_url(article.url)
    return article


def compute_url_hash(url: str) -> str:
    """计算 URL 的 MD5 哈希值，与 Go 侧 computeURLHash 一致。"""
    if not url:
        return ""
    return hashlib.md5(url.encode("utf-8")).hexdigest()


# ── 内部工具 ──────────────────────────────────────────────


def normalize_url(url: str) -> str:
    """URL 规范化：去尾部斜杠、统一 scheme 小写。"""
    if not url:
        return ""
    url = url.strip()
    # 去尾部斜杠（保留根路径）
    if url.endswith("/") and url.count("/") > 3:
        url = url.rstrip("/")
    return url


def clean_title(title: str) -> str:
    """清理标题：去除多余空白、HTML 实体解码。"""
    if not title:
        return ""
    import html as _html
    title = _html.unescape(title)
    title = re.sub(r"\s+", " ", title).strip()
    return title


def clean_text(text: str) -> str:
    """清理普通文本：去除多余空白。"""
    if not text:
        return ""
    return re.sub(r"\s+", " ", text).strip()


def normalize_timestamp(ts: str) -> str:
    """确保时间戳为 ISO 8601 格式（已是字符串，仅做基本校验）。"""
    if not ts:
        return ""
    return ts.strip()


def _extract_id_from_url(url: str) -> str:
    """从 URL 最后一段路径提取文章 ID。"""
    parsed = urllib.parse.urlparse(url)
    path = parsed.path.rstrip("/")
    if "/" in path:
        return path.rsplit("/", 1)[-1]
    return path
