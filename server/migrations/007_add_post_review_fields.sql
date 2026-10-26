-- 007_add_post_review_fields.sql
-- 添加帖子审核相关字段

ALTER TABLE posts ADD COLUMN IF NOT EXISTS review_status VARCHAR(20) DEFAULT 'approved';
ALTER TABLE posts ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMPTZ;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS reviewed_by BIGINT;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS reject_reason TEXT;

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_posts_review_status ON posts(review_status);

-- 更新现有帖子的审核状态为 approved
UPDATE posts SET review_status = 'approved' WHERE review_status IS NULL;
