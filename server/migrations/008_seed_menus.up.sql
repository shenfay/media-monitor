-- ============================================
-- 初始化默认菜单数据（与前端 menuConfig 一致）
-- ============================================

-- 概览
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_overview', 'overview', '概览', 'DashboardOutlined', '', '', NULL, 0)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_dashboard', 'dashboard', '工作台', 'DashboardOutlined', '/dashboard', 'dashboard:view', 'menu_overview', 0)
ON CONFLICT (id) DO NOTHING;

-- 成长管理
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_growth', 'growth', '成长管理', 'AimOutlined', '', '', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_family', 'family', '家庭管理', 'TeamOutlined', '/family', 'family:manage', 'menu_growth', 0),
    ('menu_goals', 'goals', '目标管理', 'AimOutlined', '/goals', 'goal:manage', 'menu_growth', 1)
ON CONFLICT (id) DO NOTHING;

-- 卡片引擎
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_card_engine', 'card-engine', '卡片引擎', 'FileTextOutlined', '', '', NULL, 2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_card_templates', 'card-templates', '卡片模板', 'FileTextOutlined', '/card-templates', 'card_template:manage', 'menu_card_engine', 0),
    ('menu_card_instances', 'card-instances', '提交记录', 'ProfileOutlined', '/card-instances', 'card_instance:view', 'menu_card_engine', 1)
ON CONFLICT (id) DO NOTHING;

-- 伙伴系统
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_companion', 'companion', '伙伴系统', 'SmileOutlined', '', '', NULL, 3)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_companions', 'companions', '伙伴管理', 'SmileOutlined', '/companions', 'companion:manage', 'menu_companion', 0)
ON CONFLICT (id) DO NOTHING;

-- 验收管理
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_acceptance', 'acceptance', '验收管理', 'CheckCircleOutlined', '', '', NULL, 4)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_acceptance_pending', 'acceptance-pending', '待验收', 'CheckCircleOutlined', '/acceptance', 'acceptance:manage', 'menu_acceptance', 0)
ON CONFLICT (id) DO NOTHING;

-- 积分系统
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_points_system', 'points-system', '积分系统', 'StarOutlined', '', '', NULL, 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_points', 'points', '积分流水', 'StarOutlined', '/points', 'points:view', 'menu_points_system', 0),
    ('menu_shop_items', 'shop-items', '商品管理', 'ShopOutlined', '/shop-items', 'shop_item:manage', 'menu_points_system', 1),
    ('menu_exchange_orders', 'exchange-orders', '兑换订单', 'SwapOutlined', '/exchange-orders', 'exchange_order:manage', 'menu_points_system', 2)
ON CONFLICT (id) DO NOTHING;

-- 用户中心
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_user', 'user', '用户中心', 'UserOutlined', '', '', NULL, 6)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_user_management', 'user-management', '用户管理', 'UserOutlined', '/users', 'user:manage', 'menu_user', 0),
    ('menu_permission_management', 'permission-management', '权限管理', 'LockOutlined', '/permissions', 'permission:manage', 'menu_user', 1),
    ('menu_menu_management', 'menu-management', '菜单管理', 'MenuOutlined', '/menus', 'menu:manage', 'menu_user', 2),
    ('menu_profile', 'profile', '个人中心', 'ProfileOutlined', '/profile', 'profile:view', 'menu_user', 3)
ON CONFLICT (id) DO NOTHING;

-- 系统
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_system', 'system', '系统', 'SettingOutlined', '', '', NULL, 7)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_operation_log', 'operation-log', '操作日志', 'AuditOutlined', '/operation-log', 'operation:log', 'menu_system', 0),
    ('menu_design_system', 'design-system', '设计规范', 'SkinOutlined', '/design-system', 'design:view', 'menu_system', 1),
    ('menu_system_settings', 'system-settings', '系统设置', 'SettingOutlined', '/settings', 'setting:manage', 'menu_system', 2)
ON CONFLICT (id) DO NOTHING;
