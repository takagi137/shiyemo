-- 复合索引优化查询性能

-- 帖子列表查询：WHERE status = 1 AND type = ? ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_posts_status_type_created 
ON posts(status, type, created_at DESC);

-- 用户帖子查询：WHERE user_id = ? AND status = 1 ORDER BY created_at DESC
CREATE INDEX IF NOT EXISTS idx_posts_user_status_created 
ON posts(user_id, status, created_at DESC);

-- 评论列表查询：WHERE post_id = ? AND status = 1 ORDER BY created_at ASC
CREATE INDEX IF NOT EXISTS idx_comments_post_status_created 
ON comments(post_id, status, created_at ASC);

-- 通知列表查询：WHERE user_id = ? AND is_read = ? ORDER BY created_at DESC
-- 注意：列名是 is_read（见 008_create_notifications.sql），不是 read。
-- 原先写作 read 会导致本迁移在任何全新部署上失败（column "read" does not exist）。
CREATE INDEX IF NOT EXISTS idx_notifications_user_read_created
ON notifications(user_id, is_read, created_at DESC);

-- 关注关系查询：WHERE follower_id = ? OR following_id = ?
CREATE INDEX IF NOT EXISTS idx_follows_follower 
ON follows(follower_id);

CREATE INDEX IF NOT EXISTS idx_follows_following 
ON follows(following_id);

COMMENT ON INDEX idx_posts_status_type_created IS '帖子列表查询复合索引';
COMMENT ON INDEX idx_posts_user_status_created IS '用户帖子查询复合索引';
COMMENT ON INDEX idx_comments_post_status_created IS '评论列表查询复合索引';
COMMENT ON INDEX idx_notifications_user_read_created IS '通知列表查询复合索引';
