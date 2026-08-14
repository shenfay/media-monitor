-- 抓取模块权限种子（Casbin 策略）
INSERT INTO casbin_rule (ptype, v0, v1)
SELECT 'p', 'role_admin', 'crawl:manage'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'role_admin' AND v1 = 'crawl:manage'
);

INSERT INTO casbin_rule (ptype, v0, v1)
SELECT 'p', 'role_operator', 'crawl:manage'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'role_operator' AND v1 = 'crawl:manage'
);
