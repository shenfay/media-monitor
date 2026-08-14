-- Worker 管理权限种子
INSERT INTO casbin_rule (ptype, v0, v1)
SELECT 'p', 'role_admin', 'worker:view'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'role_admin' AND v1 = 'worker:view'
);

INSERT INTO casbin_rule (ptype, v0, v1)
SELECT 'p', 'role_admin', 'worker:manage'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = 'role_admin' AND v1 = 'worker:manage'
);
