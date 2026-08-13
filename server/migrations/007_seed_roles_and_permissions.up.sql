-- ============================================
-- 初始化默认角色和权限
-- ============================================

-- 1. 插入默认角色
INSERT INTO roles (id, name, code, description, status) VALUES
    ('role_admin', '管理员', 'admin', '系统管理员，拥有所有权限', TRUE),
    ('role_operator', '运营', 'operator', '运营人员，管理业务内容', TRUE),
    ('role_viewer', '观察员', 'viewer', '只读角色，仅可查看部分内容', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 2. 管理员权限（全量）
INSERT INTO casbin_rule (ptype, v0, v1) VALUES
    ('p', 'role_admin', 'dashboard:view'),
    ('p', 'role_admin', 'family:manage'),
    ('p', 'role_admin', 'goal:manage'),
    ('p', 'role_admin', 'card_template:manage'),
    ('p', 'role_admin', 'card_instance:view'),
    ('p', 'role_admin', 'companion:manage'),
    ('p', 'role_admin', 'acceptance:manage'),
    ('p', 'role_admin', 'points:view'),
    ('p', 'role_admin', 'shop_item:manage'),
    ('p', 'role_admin', 'exchange_order:manage'),
    ('p', 'role_admin', 'user:manage'),
    ('p', 'role_admin', 'permission:manage'),
    ('p', 'role_admin', 'menu:manage'),
    ('p', 'role_admin', 'profile:view'),
    ('p', 'role_admin', 'operation:log'),
    ('p', 'role_admin', 'design:view'),
    ('p', 'role_admin', 'setting:manage')
ON CONFLICT DO NOTHING;

-- 3. 运营权限
INSERT INTO casbin_rule (ptype, v0, v1) VALUES
    ('p', 'role_operator', 'dashboard:view'),
    ('p', 'role_operator', 'family:manage'),
    ('p', 'role_operator', 'goal:manage'),
    ('p', 'role_operator', 'card_template:manage'),
    ('p', 'role_operator', 'card_instance:view'),
    ('p', 'role_operator', 'companion:manage'),
    ('p', 'role_operator', 'acceptance:manage'),
    ('p', 'role_operator', 'points:view'),
    ('p', 'role_operator', 'shop_item:manage'),
    ('p', 'role_operator', 'exchange_order:manage'),
    ('p', 'role_operator', 'profile:view')
ON CONFLICT DO NOTHING;

-- 4. 观察员权限
INSERT INTO casbin_rule (ptype, v0, v1) VALUES
    ('p', 'role_viewer', 'dashboard:view'),
    ('p', 'role_viewer', 'card_instance:view'),
    ('p', 'role_viewer', 'points:view'),
    ('p', 'role_viewer', 'profile:view')
ON CONFLICT DO NOTHING;
