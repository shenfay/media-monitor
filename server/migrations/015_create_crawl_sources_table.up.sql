-- 抓取数据源表
CREATE TABLE IF NOT EXISTS crawl_sources (
    id             VARCHAR(26) PRIMARY KEY,
    name           VARCHAR(200) NOT NULL,
    platform_type  VARCHAR(30) NOT NULL DEFAULT 'news',
    base_url       VARCHAR(512) DEFAULT '',
    list_endpoint  VARCHAR(512) DEFAULT '',
    nodes          TEXT,
    source_filter  VARCHAR(200) DEFAULT '',
    months         INTEGER NOT NULL DEFAULT 6,
    schedule       VARCHAR(100) DEFAULT '',
    auth           TEXT,
    enabled        BOOLEAN NOT NULL DEFAULT TRUE,
    owner_id       VARCHAR(50),
    last_crawl_at  TIMESTAMP WITH TIME ZONE,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_crawl_sources_enabled ON crawl_sources(enabled);
