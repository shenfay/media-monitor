-- ============================================
-- 新增 WebSocket 测试菜单
-- ============================================

INSERT INTO menus (id, key, label, icon, path, permission, parent_id, sort_order) VALUES
    ('menu_ws_test', 'ws_test', 'WebSocket 测试', 'ApiOutlined', '/ws-test', '', 'menu_system', 99)
ON CONFLICT (id) DO NOTHING;
