-- 回滚：将网页路径恢复为 API 路径
-- /list/xxx → /api/list?channel=xxx
UPDATE crawl_sources
SET list_endpoint = '/api/list?channel=' || regexp_replace(list_endpoint, '^/list/', ''),
    updated_at = NOW()
WHERE list_endpoint LIKE '/list/%'
  AND id LIKE 'source_hq_%';
