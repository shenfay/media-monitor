-- 删除创始人 Casbin 规则
DELETE FROM casbin_rule WHERE (ptype = 'p' AND v0 = 'role_founder') OR (ptype = 'g' AND v0 = 'user_founder');

-- 删除用户角色关联
DELETE FROM user_roles WHERE user_id = 'user_founder' AND role_id = 'role_founder';

-- 删除创始人账号
DELETE FROM users WHERE id = 'user_founder';

-- 删除创始人角色
DELETE FROM roles WHERE id = 'role_founder';
