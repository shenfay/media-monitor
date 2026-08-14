-- 抓取文章表（幂等键 source_id + url_hash）
CREATE TABLE IF NOT EXISTS crawl_articles (
    id            VARCHAR(26) PRIMARY KEY,
    source_id     VARCHAR(26) NOT NULL,
    platform      VARCHAR(50) DEFAULT '',
    external_id   VARCHAR(100) DEFAULT '',
    url           TEXT NOT NULL,
    url_hash      VARCHAR(64) DEFAULT '',
    title         VARCHAR(512) DEFAULT '',
    subtitle      TEXT,
    summary       TEXT,
    body          TEXT,
    body_format   VARCHAR(20) DEFAULT 'html',
    channel       VARCHAR(100) DEFAULT '',
    author        VARCHAR(200) DEFAULT '',
    source_name   VARCHAR(200) DEFAULT '',
    published_at  TIMESTAMP WITH TIME ZONE,
    language      VARCHAR(20) DEFAULT '',
    interactions  TEXT,
    media         TEXT,
    thread_id     VARCHAR(100),
    raw_payload   TEXT,
    fetched_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_crawl_articles_source_id ON crawl_articles(source_id);
CREATE INDEX IF NOT EXISTS idx_crawl_articles_external_id ON crawl_articles(external_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_crawl_articles_source_url_hash ON crawl_articles(source_id, url_hash);
