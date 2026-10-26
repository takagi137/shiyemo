-- Likes table (点赞)
CREATE TABLE IF NOT EXISTS likes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, post_id)
);

CREATE INDEX idx_likes_user_id ON likes(user_id);
CREATE INDEX idx_likes_post_id ON likes(post_id);

-- Bookmarks table (收藏)
CREATE TABLE IF NOT EXISTS bookmarks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, post_id)
);

CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_post_id ON bookmarks(post_id);

-- Foreign keys
ALTER TABLE likes ADD CONSTRAINT fk_likes_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE likes ADD CONSTRAINT fk_likes_post_id FOREIGN KEY (post_id) REFERENCES posts(id);
ALTER TABLE bookmarks ADD CONSTRAINT fk_bookmarks_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE bookmarks ADD CONSTRAINT fk_bookmarks_post_id FOREIGN KEY (post_id) REFERENCES posts(id);

-- Comments
COMMENT ON TABLE likes IS '点赞表';
COMMENT ON TABLE bookmarks IS '收藏表';
