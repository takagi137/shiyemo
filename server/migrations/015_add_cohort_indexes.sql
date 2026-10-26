-- 为同温层查询创建复合索引
-- 加速按城市+行业+状态查询用户

-- 用户表复合索引（城市+行业+状态）
CREATE INDEX IF NOT EXISTS idx_users_city_industry_status ON users(city, industry, status);

-- 帖子表复合索引（状态+置顶+精选+创建时间）
CREATE INDEX IF NOT EXISTS idx_posts_status_top_featured_created ON posts(status, is_top, is_featured, created_at DESC);
