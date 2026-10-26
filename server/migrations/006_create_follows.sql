-- Follows table (关注关系)
CREATE TABLE IF NOT EXISTS follows (
    id BIGSERIAL PRIMARY KEY,
    follower_id BIGINT NOT NULL,
    following_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(follower_id, following_id)
);

CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_following_id ON follows(following_id);

-- Foreign keys
ALTER TABLE follows ADD CONSTRAINT fk_follows_follower_id FOREIGN KEY (follower_id) REFERENCES users(id);
ALTER TABLE follows ADD CONSTRAINT fk_follows_following_id FOREIGN KEY (following_id) REFERENCES users(id);

-- Comments
COMMENT ON TABLE follows IS '关注关系表';
COMMENT ON COLUMN follows.follower_id IS '关注者ID';
COMMENT ON COLUMN follows.following_id IS '被关注者ID';
