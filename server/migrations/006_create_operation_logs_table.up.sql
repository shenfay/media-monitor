-- ============================================
-- 统一操作日志表
-- ============================================
-- 合并审计日志（audit_logs）与活动日志（activity_logs）为单表
-- 通过 category 字段区分：AUTH / USER / SYSTEM / BIZ
-- 通过 action 字段标准化事件类型：AUTH.LOGIN.SUCCESS / USER.PROFILE.UPDATED / ...

CREATE TABLE IF NOT EXISTS operation_logs (
    id          VARCHAR(50) PRIMARY KEY,
    user_id     VARCHAR(50) NOT NULL,
    email       VARCHAR(255),
    action      VARCHAR(80) NOT NULL,
    category    VARCHAR(30) NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'SUCCESS',
    ip          VARCHAR(45),
    user_agent  VARCHAR(500),
    device      VARCHAR(100),
    browser     VARCHAR(50),
    os          VARCHAR(50),
    metadata    JSONB DEFAULT '{}'::jsonb,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_operation_logs_user_created ON operation_logs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_operation_logs_category_created ON operation_logs(category, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_operation_logs_action ON operation_logs(action);
CREATE INDEX IF NOT EXISTS idx_operation_logs_action_created ON operation_logs(action, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_operation_logs_status ON operation_logs(status);
CREATE INDEX IF NOT EXISTS idx_operation_logs_failed ON operation_logs(status, created_at DESC) WHERE status = 'FAILED';

-- 添加注释
COMMENT ON TABLE operation_logs IS '统一操作日志表 - 记录安全事件与业务操作';
COMMENT ON COLUMN operation_logs.id IS '日志 ID（ULID 格式）';
COMMENT ON COLUMN operation_logs.user_id IS '用户 ID';
COMMENT ON COLUMN operation_logs.email IS '用户邮箱（冗余字段，便于查询）';
COMMENT ON COLUMN operation_logs.action IS '操作类型（如：AUTH.LOGIN.SUCCESS, USER.PROFILE.UPDATED）';
COMMENT ON COLUMN operation_logs.category IS '操作分类（AUTH/USER/SYSTEM/BIZ）';
COMMENT ON COLUMN operation_logs.status IS '操作状态（SUCCESS/FAILED）';
COMMENT ON COLUMN operation_logs.ip IS 'IP 地址（支持 IPv6）';
COMMENT ON COLUMN operation_logs.user_agent IS 'User-Agent 原始字符串';
COMMENT ON COLUMN operation_logs.device IS '设备类型（mobile/tablet/desktop）';
COMMENT ON COLUMN operation_logs.browser IS '浏览器名称';
COMMENT ON COLUMN operation_logs.os IS '操作系统';
COMMENT ON COLUMN operation_logs.metadata IS '元数据（JSONB 格式）';
COMMENT ON COLUMN operation_logs.created_at IS '创建时间（带时区）';
