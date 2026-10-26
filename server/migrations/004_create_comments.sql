-- Comments table
CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    parent_id BIGINT,  -- 父评论 ID，NULL 表示主评论
    content TEXT NOT NULL,
    like_count INT DEFAULT 0,
    status INT DEFAULT 1,  -- 0=草稿 1=已发布 2=隐藏 3=删除
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX idx_comments_status ON comments(status);
CREATE INDEX idx_comments_created_at ON comments(created_at);

-- Foreign keys
ALTER TABLE comments ADD CONSTRAINT fk_comments_post_id FOREIGN KEY (post_id) REFERENCES posts(id);
ALTER TABLE comments ADD CONSTRAINT fk_comments_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE comments ADD CONSTRAINT fk_comments_parent_id FOREIGN KEY (parent_id) REFERENCES comments(id);

-- Comments
COMMENT ON TABLE comments IS '评论表';
COMMENT ON COLUMN comments.parent_id IS '父评论ID，NULL表示主评论';
COMMENT ON COLUMN comments.status IS '状态: 0-草稿, 1-已发布, 2-隐藏, 3-删除';
