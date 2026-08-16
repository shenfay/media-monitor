-- ============================================
-- 批量添加环球网各频道数据源
-- 基于 m.huanqiu.com 导航栏频道列表
-- ============================================

-- 领航
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_lh', '环球网-领航', 'news', 'https://m.huanqiu.com', '/list/lh', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 新征程
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_fjxzc', '环球网-新征程', 'news', 'https://m.huanqiu.com', '/list/fjxzc', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 国际
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_world', '环球网-国际', 'news', 'https://m.huanqiu.com', '/list/world', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 台海
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_taiwan', '环球网-台海', 'news', 'https://m.huanqiu.com', '/list/taiwan', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 社会
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_society', '环球网-社会', 'news', 'https://m.huanqiu.com', '/list/society', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 军事
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_mil', '环球网-军事', 'news', 'https://m.huanqiu.com', '/list/mil', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 社评
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_editorial', '环球网-社评', 'news', 'https://m.huanqiu.com', '/list/editorial', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 国内
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_china', '环球网-国内', 'news', 'https://m.huanqiu.com', '/list/china', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 评论
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_opinion', '环球网-评论', 'news', 'https://m.huanqiu.com', '/list/opinion', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 海外看中国
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_oversea', '环球网-海外看中国', 'news', 'https://m.huanqiu.com', '/list/oversea', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 视频
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_v', '环球网-视频', 'news', 'https://m.huanqiu.com', '/list/v', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 财经
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_finance', '环球网-财经', 'news', 'https://m.huanqiu.com', '/list/finance', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 中部
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_zy', '环球网-中部', 'news', 'https://m.huanqiu.com', '/list/zy', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 汽车
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_auto', '环球网-汽车', 'news', 'https://m.huanqiu.com', '/list/auto', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 科技
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_tech', '环球网-科技', 'news', 'https://m.huanqiu.com', '/list/tech', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 单仁平
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_shanrenping', '环球网-单仁平', 'news', 'https://m.huanqiu.com', '/list/shanrenping', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 低空经济
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_uav', '环球网-低空经济', 'news', 'https://m.huanqiu.com', '/list/uav', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 文旅
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_go', '环球网-文旅', 'news', 'https://m.huanqiu.com', '/list/go', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 亲子
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_qinzi', '环球网-亲子', 'news', 'https://m.huanqiu.com', '/list/qinzi', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 健康
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_health', '环球网-健康', 'news', 'https://m.huanqiu.com', '/list/health', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 大文娱
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_ent', '环球网-大文娱', 'news', 'https://m.huanqiu.com', '/list/ent', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 艺术·文博
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_art', '环球网-艺术文博', 'news', 'https://m.huanqiu.com', '/list/art', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 时尚
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_fashion', '环球网-时尚', 'news', 'https://m.huanqiu.com', '/list/fashion', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 女性
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_women', '环球网-女性', 'news', 'https://m.huanqiu.com', '/list/women', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 体育
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_sports', '环球网-体育', 'news', 'https://m.huanqiu.com', '/list/sports', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 丝路
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_silkroad', '环球网-丝路', 'news', 'https://m.huanqiu.com', '/list/silkroad', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 教育
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_lx', '环球网-教育', 'news', 'https://m.huanqiu.com', '/list/lx', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 房产
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_house', '环球网-房产', 'news', 'https://m.huanqiu.com', '/list/house', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 公益
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_hope', '环球网-公益', 'news', 'https://m.huanqiu.com', '/list/hope', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 城市
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_city', '环球网-城市', 'news', 'https://m.huanqiu.com', '/list/city', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 商业
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_biz', '环球网-商业', 'news', 'https://m.huanqiu.com', '/list/biz', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 听书
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_book', '环球网-听书', 'news', 'https://m.huanqiu.com', '/list/book', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 农业
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_xy', '环球网-农业', 'news', 'https://m.huanqiu.com', '/list/xy', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 文化
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_cul', '环球网-文化', 'news', 'https://m.huanqiu.com', '/list/cul', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 消费
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_quality', '环球网-消费', 'news', 'https://m.huanqiu.com', '/list/quality', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 能源
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_energy', '环球网-能源', 'news', 'https://m.huanqiu.com', '/list/energy', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 长三角
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_yrd', '环球网-长三角', 'news', 'https://m.huanqiu.com', '/list/yrd', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 融媒联播
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_media', '环球网-融媒联播', 'news', 'https://m.huanqiu.com', '/list/media', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 产业新闻
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_capital', '环球网-产业新闻', 'news', 'https://m.huanqiu.com', '/list/capital', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 消防
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_anquan', '环球网-消防', 'news', 'https://m.huanqiu.com', '/list/anquan', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 英语
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_english', '环球网-英语', 'news', 'https://m.huanqiu.com', '/list/english', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 法语
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_french', '环球网-法语', 'news', 'https://m.huanqiu.com', '/list/french', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 俄语
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_russian', '环球网-俄语', 'news', 'https://m.huanqiu.com', '/list/russian', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 西语
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_spanish', '环球网-西语', 'news', 'https://m.huanqiu.com', '/list/spanish', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 阿语
INSERT INTO crawl_sources (id, name, platform_type, base_url, list_endpoint, source_filter, months, tags, enabled)
VALUES ('source_hq_arabic', '环球网-阿语', 'news', 'https://m.huanqiu.com', '/list/arabic', '环球网', 6, '["huanqiu"]', TRUE)
ON CONFLICT (id) DO NOTHING;
