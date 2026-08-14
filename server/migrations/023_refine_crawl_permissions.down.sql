-- 回滚抓取模块细粒度权限
DELETE FROM casbin_rule WHERE ptype = 'p' AND v0 = 'role_admin'
    AND v1 IN ('crawl:view', 'source:view', 'source:manage', 'task:view', 'task:manage', 'article:view');
