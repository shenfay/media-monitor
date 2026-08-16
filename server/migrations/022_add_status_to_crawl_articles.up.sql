-- 文章抓取状态：pending(列表阶段) | completed(正文已抓取) | failed(抓取失败)
ALTER TABLE crawl_articles ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'pending';
CREATE INDEX IF NOT EXISTS idx_crawl_articles_status ON crawl_articles(status);
