-- ============================================
-- 回滚菜单重构：恢复原始菜单结构
-- ============================================

-- 恢复占位菜单组
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_growth', 'growth', '成长管理', 'AimOutlined', '', '', NULL, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_family', 'family', '家庭管理', 'TeamOutlined', '/family', 'family:manage', 'menu_growth', 0),
    ('menu_goals', 'goals', '目标管理', 'AimOutlined', '/goals', 'goal:manage', 'menu_growth', 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_card_engine', 'card-engine', '卡片引擎', 'FileTextOutlined', '', '', NULL, 2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_card_templates', 'card-templates', '卡片模板', 'FileTextOutlined', '/card-templates', 'card_template:manage', 'menu_card_engine', 0),
    ('menu_card_instances', 'card-instances', '提交记录', 'ProfileOutlined', '/card-instances', 'card_instance:view', 'menu_card_engine', 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_companion', 'companion', '伙伴系统', 'SmileOutlined', '', '', NULL, 3)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_companions', 'companions', '伙伴管理', 'SmileOutlined', '/companions', 'companion:manage', 'menu_companion', 0)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_acceptance', 'acceptance', '验收管理', 'CheckCircleOutlined', '', '', NULL, 4)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_acceptance_pending', 'acceptance-pending', '待验收', 'CheckCircleOutlined', '/acceptance', 'acceptance:manage', 'menu_acceptance', 0)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_points_system', 'points-system', '积分系统', 'StarOutlined', '', '', NULL, 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_points', 'points', '积分流水', 'StarOutlined', '/points', 'points:view', 'menu_points_system', 0),
    ('menu_shop_items', 'shop-items', '商品管理', 'ShopOutlined', '/shop-items', 'shop_item:manage', 'menu_points_system', 1),
    ('menu_exchange_orders', 'exchange-orders', '兑换订单', 'SwapOutlined', '/exchange-orders', 'exchange_order:manage', 'menu_points_system', 2)
ON CONFLICT (id) DO NOTHING;

-- 恢复个人中心到用户中心
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_profile', 'profile', '个人中心', 'ProfileOutlined', '/profile', 'profile:view', 'menu_user', 3)
ON CONFLICT (id) DO NOTHING;

-- Worker 管理移回系统
UPDATE menus SET parent_id = 'menu_system', sort_order = 6 WHERE id = 'menu_worker';

-- 移除抓取管理菜单组及子菜单
DELETE FROM menus WHERE id IN ('menu_crawl', 'menu_sources', 'menu_tasks', 'menu_articles');

-- 恢复系统菜单排序
UPDATE menus SET sort_order = 0 WHERE id = 'menu_operation_log';
UPDATE menus SET sort_order = 1 WHERE id = 'menu_design_system';
UPDATE menus SET sort_order = 2 WHERE id = 'menu_system_settings';
