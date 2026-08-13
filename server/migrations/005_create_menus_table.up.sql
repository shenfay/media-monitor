-- ============================================
-- 菜单管理表
-- ============================================

CREATE TABLE IF NOT EXISTS menus (
    id VARCHAR(50) PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    label VARCHAR(100) NOT NULL,
    icon VARCHAR(100) DEFAULT '',
    path VARCHAR(200) DEFAULT '',
    permission VARCHAR(100) DEFAULT '',
    parent_id VARCHAR(50) REFERENCES menus(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    status BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_menus_parent_id ON menus(parent_id);
CREATE INDEX IF NOT EXISTS idx_menus_sort_order ON menus(sort_order);

COMMENT ON TABLE menus IS '菜单管理表';
COMMENT ON COLUMN menus.id IS '菜单 ID (ULID)';
COMMENT ON COLUMN menus.key IS '菜单标识（唯一）';
COMMENT ON COLUMN menus.label IS '菜单名称';
COMMENT ON COLUMN menus.icon IS '图标名称';
COMMENT ON COLUMN menus.path IS '路由路径';
COMMENT ON COLUMN menus.permission IS '权限标识';
COMMENT ON COLUMN menus.parent_id IS '父菜单 ID，空表示顶级菜单';
COMMENT ON COLUMN menus.sort_order IS '排序序号';
COMMENT ON COLUMN menus.status IS '是否启用';
