-- 抓取模块细粒度权限种子
INSERT INTO casbin_rule (ptype, v0, v1)
SELECT 'p', 'role_admin', perm
FROM (VALUES ('crawl:view'), ('source:view'), ('source:manage'),
             ('task:view'), ('task:manage'), ('article:view')) AS t(perm)
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'role_admin' AND v1 = perm
);
