-- ============================================
-- RBAC 权限模型表（Casbin 引擎）
-- ============================================

-- 1. 角色表
CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    status BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE roles IS '角色表';
COMMENT ON COLUMN roles.id IS '角色 ID (ULID)';
COMMENT ON COLUMN roles.name IS '角色名称';
COMMENT ON COLUMN roles.code IS '角色编码（唯一标识）';
COMMENT ON COLUMN roles.description IS '角色描述';
COMMENT ON COLUMN roles.status IS '是否启用';

-- 2. 用户角色关联表（多对多）
CREATE TABLE IF NOT EXISTS user_roles (
    user_id VARCHAR(50) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id VARCHAR(50) NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);

COMMENT ON TABLE user_roles IS '用户角色关联表';

-- 3. Casbin 策略规则表
CREATE TABLE IF NOT EXISTS casbin_rule (
    id    SERIAL PRIMARY KEY,
    ptype VARCHAR(100) NOT NULL,
    v0    VARCHAR(100) NOT NULL DEFAULT '',
    v1    VARCHAR(100) NOT NULL DEFAULT '',
    v2    VARCHAR(100) NOT NULL DEFAULT '',
    v3    VARCHAR(100) NOT NULL DEFAULT '',
    v4    VARCHAR(100) NOT NULL DEFAULT '',
    v5    VARCHAR(100) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_casbin_rule_ptype ON casbin_rule(ptype);

COMMENT ON TABLE casbin_rule IS 'Casbin 策略规则表';
