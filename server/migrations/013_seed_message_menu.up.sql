-- ============================================
-- 新增消息管理菜单
-- ============================================

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_message', 'message', '消息管理', 'MessageOutlined', '/messages', 'message:view', 'menu_system', 3)
ON CONFLICT (id) DO NOTHING;
