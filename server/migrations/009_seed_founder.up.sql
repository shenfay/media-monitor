-- ============================================
-- 初始化创始人账号
-- ============================================

-- 1. 插入创始人角色
INSERT INTO roles (id, name, code, description, status) VALUES
    ('role_founder', '创始人', 'founder', '项目创始人，拥有所有权限', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 2. 插入创始人账号（密码：Founder@2026）
INSERT INTO users (id, email, name, password, email_verified, locked, failed_attempts) VALUES
    ('user_founder', 'admin@example.com', '管理员',
     '$2b$12$/gqsOy.VaUZN5IuUBjtZcOb1y4Cl/DWcVc6UzcYNzy9w0u1/piXcm',
     TRUE, FALSE, 0)
ON CONFLICT (id) DO NOTHING;

-- 3. 关联创始人角色
INSERT INTO user_roles (user_id, role_id) VALUES
    ('user_founder', 'role_founder')
ON CONFLICT (user_id, role_id) DO NOTHING;

-- 4. Casbin g 规则：用户→角色
INSERT INTO casbin_rule (ptype, v0, v1) VALUES
    ('g', 'user_founder', 'role_founder')
ON CONFLICT DO NOTHING;

-- 5. Casbin p 规则：创始人拥有所有权限
INSERT INTO casbin_rule (ptype, v0, v1) VALUES
    ('p', 'role_founder', 'dashboard:view'),
    ('p', 'role_founder', 'family:manage'),
    ('p', 'role_founder', 'goal:manage'),
    ('p', 'role_founder', 'card_template:manage'),
    ('p', 'role_founder', 'card_instance:view'),
    ('p', 'role_founder', 'companion:manage'),
    ('p', 'role_founder', 'acceptance:manage'),
    ('p', 'role_founder', 'points:view'),
    ('p', 'role_founder', 'shop_item:manage'),
    ('p', 'role_founder', 'exchange_order:manage'),
    ('p', 'role_founder', 'user:manage'),
    ('p', 'role_founder', 'permission:manage'),
    ('p', 'role_founder', 'menu:manage'),
    ('p', 'role_founder', 'profile:view'),
    ('p', 'role_founder', 'operation:log'),
    ('p', 'role_founder', 'design:view'),
    ('p', 'role_founder', 'setting:manage')
ON CONFLICT DO NOTHING;
