-- ============================================
-- 系统设置表
-- ============================================
-- 存储可动态修改的业务配置和功能开关
-- value 统一使用 JSONB 类型，天然支持 string/int/bool/object
-- 基础设施配置（端口、数据库连接等）保留在 yaml 文件中

CREATE TABLE IF NOT EXISTS system_settings (
    id          BIGSERIAL PRIMARY KEY,
    key         VARCHAR(100) NOT NULL UNIQUE,
    value       JSONB NOT NULL,
    category    VARCHAR(50) NOT NULL,
    label       VARCHAR(200),
    description TEXT,
    updated_by  VARCHAR(50) REFERENCES users(id),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_system_settings_category ON system_settings(category);

-- 添加注释
COMMENT ON TABLE system_settings IS '系统设置表 - 存储可动态修改的业务配置和功能开关';
COMMENT ON COLUMN system_settings.id IS '自增主键';
COMMENT ON COLUMN system_settings.key IS '设置项唯一标识（如 site_name、enable_register）';
COMMENT ON COLUMN system_settings.value IS '设置值（JSONB 格式，支持 string/int/bool/object）';
COMMENT ON COLUMN system_settings.category IS '设置分类（basic/toggle/business/notification）';
COMMENT ON COLUMN system_settings.label IS '前端显示标签';
COMMENT ON COLUMN system_settings.description IS '设置说明';
COMMENT ON COLUMN system_settings.updated_by IS '最后修改人 ID';
COMMENT ON COLUMN system_settings.created_at IS '创建时间';
COMMENT ON COLUMN system_settings.updated_at IS '更新时间';
