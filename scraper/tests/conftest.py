"""pytest 全局 fixtures。"""
import pytest


@pytest.fixture
def sample_source():
    """构造测试用 Source。"""
    from scrapers.contracts.source import Source
    return Source(
        id="test-huanqiu",
        name="环球网",
        platform_type="news",
        source_filter="环球网",
        months=6,
    )


@pytest.fixture
def sample_article():
    """构造测试用 Article。"""
    from scrapers.contracts.article import Article
    return Article(
        source_id="test-huanqiu",
        platform="news",
        title="测试文章",
        url="https://m.huanqiu.com/article/12345",
        external_id="12345",
        source_name="环球网",
        published_at="2026-08-14T00:00:00+00:00",
    )
