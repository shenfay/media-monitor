-- 为菜单表添加 permissions 字段（JSON 数组，支持多个权限）
-- 保留原有 permission 字段用于向后兼容（菜单可见性判断）
ALTER TABLE menus ADD COLUMN IF NOT EXISTS permissions TEXT NOT NULL DEFAULT '[]';

-- 更新现有菜单的 permissions 字段
-- 每个菜单关联其 view + manage 权限，权限管理页面勾选菜单即分配全部关联权限

-- 抓取管理
UPDATE menus SET permissions = '["crawl:view","crawl:manage"]' WHERE key = 'crawl';
UPDATE menus SET permissions = '["source:view","source:manage"]' WHERE key = 'sources';
UPDATE menus SET permissions = '["task:view","task:manage"]' WHERE key = 'tasks';
UPDATE menus SET permissions = '["article:view"]' WHERE key = 'articles';
UPDATE menus SET permissions = '["worker:view","worker:manage"]' WHERE key = 'worker';

-- 系统管理
UPDATE menus SET permissions = '["user:manage"]' WHERE key = 'user-management';
UPDATE menus SET permissions = '["permission:manage"]' WHERE key = 'permission-management';
UPDATE menus SET permissions = '["menu:manage"]' WHERE key = 'menu-management';
UPDATE menus SET permissions = '["operation:log"]' WHERE key = 'operation-log';
UPDATE menus SET permissions = '["setting:manage"]' WHERE key = 'system-settings';

-- 消息管理
UPDATE menus SET permissions = '["message:view","message:manage"]' WHERE key = 'message';

-- 仪表盘
UPDATE menus SET permissions = '["dashboard:view"]' WHERE key = 'dashboard';

-- 设计规范
UPDATE menus SET permissions = '["design:view"]' WHERE key = 'design-system';

-- WebSocket 测试
UPDATE menus SET permissions = '["ws_test:view"]' WHERE key = 'ws_test';
