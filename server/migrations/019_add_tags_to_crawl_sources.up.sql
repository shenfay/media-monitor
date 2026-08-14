-- 为 crawl_sources 表添加 tags 列（能力标签，用于 Stream 路由）
ALTER TABLE crawl_sources ADD COLUMN IF NOT EXISTS tags TEXT DEFAULT '[]';
