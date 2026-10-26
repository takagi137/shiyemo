-- Posts table
CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL,  -- emotion/insight/case/action
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    images JSONB DEFAULT '[]',
    tags JSONB DEFAULT '[]',
    like_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    status INT DEFAULT 1,  -- 0=草稿 1=已发布 2=隐藏 3=删除
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_type ON posts(type);
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX idx_posts_deleted_at ON posts(deleted_at) WHERE deleted_at IS NOT NULL;

-- Foreign key
ALTER TABLE posts ADD CONSTRAINT fk_posts_user_id FOREIGN KEY (user_id) REFERENCES users(id);

-- Comments
COMMENT ON TABLE posts IS '帖子表';
COMMENT ON COLUMN posts.type IS '帖子类型: emotion-情绪, insight-认知, case-案例, action-行动';
COMMENT ON COLUMN posts.images IS '图片URL数组';
COMMENT ON COLUMN posts.tags IS '标签数组';
COMMENT ON COLUMN posts.status IS '状态: 0-草稿, 1-已发布, 2-隐藏, 3-删除';
