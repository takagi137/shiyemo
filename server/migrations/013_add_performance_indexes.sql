-- 性能优化：为高频查询添加复合索引
-- 分析了 repo 层所有查询模式后确定的优化点

-- 1. likes 表：优化 IsLiked/Unlike 查询
-- 查询模式：WHERE user_id = ? AND post_id = ?
-- 现有：单列索引 idx_likes_user_id, idx_likes_post_id
-- 优化：复合索引可直接定位到唯一行，避免两次索引扫描
CREATE INDEX IF NOT EXISTS idx_likes_user_post
ON likes(user_id, post_id);

COMMENT ON INDEX idx_likes_user_post IS '点赞查询复合索引：优化 IsLiked/Unlike';

-- 2. bookmarks 表：优化 IsBookmarked/Unbookmark/收藏列表查询
-- 查询模式：WHERE user_id = ? AND post_id = ? 或 WHERE user_id = ? ORDER BY created_at DESC
-- 现有：单列索引 idx_bookmarks_user_id, idx_bookmarks_post_id
-- 优化：复合索引同时支持精确查询和用户收藏列表分页
CREATE INDEX IF NOT EXISTS idx_bookmarks_user_post
ON bookmarks(user_id, post_id);

CREATE INDEX IF NOT EXISTS idx_bookmarks_user_created
ON bookmarks(user_id, created_at DESC);

COMMENT ON INDEX idx_bookmarks_user_post IS '收藏精确查询复合索引：优化 IsBookmarked/Unbookmark';
COMMENT ON INDEX idx_bookmarks_user_created IS '收藏列表索引：优化 ListBookmarks 分页';

-- 3. follows 表：优化 IsFollowing/Unfollow 查询
-- 查询模式：WHERE follower_id = ? AND following_id = ?
-- 现有：单列索引 idx_follows_follower, idx_follows_following
-- 优化：复合索引可直接定位到唯一行
CREATE INDEX IF NOT EXISTS idx_follows_follower_following
ON follows(follower_id, following_id);

COMMENT ON INDEX idx_follows_follower_following IS '关注关系精确查询索引：优化 IsFollowing/Unfollow';
