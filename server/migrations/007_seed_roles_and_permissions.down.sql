-- 清除默认角色及其权限
DELETE FROM casbin_rule WHERE ptype = 'p' AND v0 IN ('role_admin', 'role_operator', 'role_viewer');
DELETE FROM roles WHERE id IN ('role_admin', 'role_operator', 'role_viewer');
