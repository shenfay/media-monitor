-- 移除 Worker 管理权限
DELETE FROM casbin_rule WHERE ptype = 'p' AND v0 = 'role_admin' AND v1 IN ('worker:view', 'worker:manage');
