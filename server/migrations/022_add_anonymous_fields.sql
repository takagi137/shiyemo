-- 补齐匿名发帖/匿名评论字段
--
-- 背景：internal/post/model.go 与 internal/comment/model.go 均定义了
-- IsAnonymous / AnonymousName，前端 feed 页也有「🎭 匿名」标记逻辑
-- （web/src/app/feed/page.tsx），但从未有对应迁移。
-- 结果是全新部署上发帖直接失败：
--   ERROR: column "is_anonymous" of relation "posts" does not exist
--
-- 字段定义对应各 model.go 的 gorm tag。

ALTER TABLE posts    ADD COLUMN IF NOT EXISTS is_anonymous   BOOLEAN     DEFAULT FALSE;
ALTER TABLE posts    ADD COLUMN IF NOT EXISTS anonymous_name VARCHAR(50);

ALTER TABLE comments ADD COLUMN IF NOT EXISTS is_anonymous   BOOLEAN     DEFAULT FALSE;
ALTER TABLE comments ADD COLUMN IF NOT EXISTS anonymous_name VARCHAR(50);

COMMENT ON COLUMN posts.is_anonymous      IS '是否匿名发布';
COMMENT ON COLUMN posts.anonymous_name    IS '匿名展示名称';
COMMENT ON COLUMN comments.is_anonymous   IS '是否匿名评论';
COMMENT ON COLUMN comments.anonymous_name IS '匿名展示名称';
