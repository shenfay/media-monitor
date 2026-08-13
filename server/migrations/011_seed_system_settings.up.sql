-- ============================================
-- 初始化系统设置默认值
-- ============================================

-- Tab 1：基础配置
INSERT INTO system_settings (key, value, category, label, description) VALUES
    ('site_name', '"管理系统"', 'basic', '站点名称', '站点显示名称'),
    ('logo_url', '""', 'basic', '站点 Logo', 'Logo 图片 URL'),
    ('default_language', '"zh-CN"', 'basic', '默认语言', '系统默认语言'),
    ('timezone', '"Asia/Shanghai"', 'basic', '时区', '系统默认时区'),
    ('session_timeout', '30', 'basic', '会话超时', '会话超时时间（分钟）')
ON CONFLICT (key) DO NOTHING;

-- Tab 2：功能开关
INSERT INTO system_settings (key, value, category, label, description) VALUES
    ('enable_register', 'true', 'toggle', '开放注册', '是否开放用户注册'),
    ('enable_audit_log', 'true', 'toggle', '审计日志', '是否启用审计日志'),
    ('enable_notification', 'true', 'toggle', '消息通知', '是否启用消息通知')
ON CONFLICT (key) DO NOTHING;

-- Tab 3：业务规则
INSERT INTO system_settings (key, value, category, label, description) VALUES
    ('max_children_per_family', '3', 'business', '家庭最大孩子数', '每个家庭最多可添加的孩子数量'),
    ('max_daily_cards', '3', 'business', '每日最大卡片数', '每天最多分配的卡片数量'),
    ('max_goal_cards', '5', 'business', '目标最大关联卡片数', '每个目标最多关联的卡片类型数量'),
    ('streak_bonus_rules', '{"1-3":1.0,"4-7":1.2,"8-14":1.5,"15-21":2.0,"22+":2.5}', 'business', '积分连续加成规则', '连续打卡天数对应的积分倍率'),
    ('xp_level_divisor', '100', 'business', 'XP等级公式除数', '伙伴等级 = floor(总XP / 除数) + 1'),
    ('default_companion_name', '"波奇"', 'business', '伙伴初始形态', '新创建伙伴的默认名称'),
    ('default_goal_reward', '50', 'business', '积分达成默认奖金', '目标达成时默认发放的积分奖励')
ON CONFLICT (key) DO NOTHING;

-- Tab 4：通知设置 - 事件开关
INSERT INTO system_settings (key, value, category, label, description) VALUES
    ('notify_acceptance', 'true', 'notification', '验收提醒', '孩子提交卡片后是否通知家长验收'),
    ('notify_goal_progress', 'true', 'notification', '目标进度通知', '目标进度变化时是否推送通知'),
    ('notify_companion_status', 'true', 'notification', '伙伴状态通知', '伙伴状态变化时是否推送通知'),
    ('notify_streak_inactive', 'true', 'notification', '连续未完成提醒', '连续多天未完成任务时是否推送提醒')
ON CONFLICT (key) DO NOTHING;

-- Tab 4：通知设置 - 渠道配置
INSERT INTO system_settings (key, value, category, label, description) VALUES
    ('channel_email', '{"host":"","port":0,"user":"","password":"","from":""}', 'notification', '邮件渠道', 'SMTP 邮件发送配置（密码 AES 加密存储）'),
    ('channel_webhook', '{"url":"","secret":""}', 'notification', 'Webhook 渠道', 'Webhook 推送配置（secret AES 加密存储）')
ON CONFLICT (key) DO NOTHING;
