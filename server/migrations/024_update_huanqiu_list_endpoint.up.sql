-- 将环球网频道的 list_endpoint 从 API 路径改为网页路径
-- /api/list?channel=xxx → /list/xxx
UPDATE crawl_sources
SET list_endpoint = '/list/' || regexp_replace(list_endpoint, '^/api/list\?channel=', ''),
    updated_at = NOW()
WHERE list_endpoint LIKE '/api/list?channel=%';
