DROP INDEX IF EXISTS idx_crawl_articles_status;
ALTER TABLE crawl_articles DROP COLUMN IF EXISTS status;
