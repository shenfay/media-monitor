-- Worker 管理菜单
INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_worker', 'worker', 'Worker 管理', 'BuildOutlined', '/crawl/workers', 'worker:view', 'menu_system', 6)
ON CONFLICT (id) DO NOTHING;
