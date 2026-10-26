-- Add following_count and follower_count to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS following_count INT DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS follower_count INT DEFAULT 0;

-- Comments
COMMENT ON COLUMN users.following_count IS '关注数';
COMMENT ON COLUMN users.follower_count IS '粉丝数';
