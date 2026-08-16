-- 回滚：移除 menus 表的 permissions 字段
ALTER TABLE menus DROP COLUMN IF EXISTS permissions;
