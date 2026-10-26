-- 添加 is_top 和 is_featured 字段到 posts 表
ALTER TABLE posts ADD COLUMN IF NOT EXISTS is_top BOOLEAN DEFAULT FALSE;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS is_featured BOOLEAN DEFAULT FALSE;

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_posts_is_top ON posts(is_top);
CREATE INDEX IF NOT EXISTS idx_posts_is_featured ON posts(is_featured);
