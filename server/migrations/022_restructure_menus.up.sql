-- ============================================
-- 管理后台菜单重构
-- 1. 移除占位菜单组（成长管理/卡片引擎/伙伴系统/验收管理/积分系统）
-- 2. 个人中心从菜单移入顶栏头像下拉
-- 3. 新增“抓取管理”顶级菜单组
-- ============================================

-- 1. 移除占位菜单组及其子菜单
DELETE FROM menus WHERE id IN (
    'menu_growth', 'menu_family', 'menu_goals',
    'menu_card_engine', 'menu_card_templates', 'menu_card_instances',
    'menu_companion', 'menu_companions',
    'menu_acceptance', 'menu_acceptance_pending',
    'menu_points_system', 'menu_points', 'menu_shop_items', 'menu_exchange_orders'
);

-- 2. 个人中心从菜单移除（改由顶栏头像下拉访问）
DELETE FROM menus WHERE id = 'menu_profile';

-- 3. 新增"抓取管理"顶级菜单组
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_crawl', 'crawl', '抓取管理', 'GlobalOutlined', '', 'crawl:view', NULL, 1)
ON CONFLICT (id) DO NOTHING;

-- 4. Worker 管理移入抓取管理
UPDATE menus SET parent_id = 'menu_crawl', sort_order = 3 WHERE id = 'menu_worker';

-- 5. 新增抓取子菜单
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_sources', 'sources', '数据源', 'DatabaseOutlined', '/crawl/sources', 'source:view', 'menu_crawl', 0),
    ('menu_tasks', 'tasks', '任务记录', 'UnorderedListOutlined', '/crawl/tasks', 'task:view', 'menu_crawl', 1),
    ('menu_articles', 'articles', '文章管理', 'ReadOutlined', '/crawl/articles', 'article:view', 'menu_crawl', 2)
ON CONFLICT (id) DO NOTHING;

-- 6. 重新排列顶级菜单 sort_order
UPDATE menus SET sort_order = 0 WHERE id = 'menu_overview';
UPDATE menus SET sort_order = 2 WHERE id = 'menu_user';
UPDATE menus SET sort_order = 3 WHERE id = 'menu_system';
