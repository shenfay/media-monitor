"""HTTP 获取层。

设计：纯标准库 urllib 即可工作（无第三方依赖即可 run 模式）；
若环境装了 httpx 则自动优先使用（更快、连接复用更好）。
带超时与指数退避重试。后续 Crawlee/Playwright 作为"渲染+登录"后端在此层扩展。
"""
from __future__ import annotations

import json
import ssl
import time
import urllib.error
import urllib.request

from crawl.config import settings

DEFAULT_USER_AGENT = "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15"

try:
    import httpx  # type: ignore

    _HAS_HTTPX = True
except Exception:  # pragma: no cover
    _HAS_HTTPX = False


def _headers(headers):
    hdrs = {"User-Agent": DEFAULT_USER_AGENT}
    if headers:
        hdrs.update(headers)
    return hdrs


def get_text(url: str, headers=None, timeout=None, retries=None) -> str:
    timeout = timeout or settings.http_timeout
    retries = retries or settings.http_retries
    hdrs = _headers(headers)
    last = None
    for i in range(max(1, retries)):
        try:
            req = urllib.request.Request(url, headers=hdrs)
            with urllib.request.urlopen(req, timeout=timeout) as r:
                return r.read().decode("utf-8", "replace")
        except Exception as e:  # noqa: BLE001
            last = e
            time.sleep(1 + i)
    raise last or RuntimeError("request failed")


def get_json(url: str, headers=None, timeout=None, retries=None):
    return json.loads(get_text(url, headers=headers, timeout=timeout, retries=retries))


def get_text_httpx(url: str, headers=None, timeout=None, retries=None) -> str:
    timeout = timeout or settings.http_timeout
    retries = retries or settings.http_retries
    hdrs = _headers(headers)
    last = None
    for i in range(max(1, retries)):
        try:
            with httpx.Client(timeout=timeout, follow_redirects=True, headers=hdrs) as c:
                r = c.get(url)
                r.raise_for_status()
                return r.text
        except Exception as e:  # noqa: BLE001
            last = e
            time.sleep(1 + i)
    raise last or RuntimeError("request failed")


def get_text_auto(url: str, headers=None, timeout=None, retries=None) -> str:
    """优先 httpx，失败回退 urllib。"""
    if _HAS_HTTPX:
        try:
            return get_text_httpx(url, headers=headers, timeout=timeout, retries=retries)
        except Exception:  # noqa: BLE001
            pass
    return get_text(url, headers=headers, timeout=timeout, retries=retries)
