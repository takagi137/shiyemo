-- 创建标签表
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 创建帖子标签关联表
CREATE TABLE IF NOT EXISTS post_tags (
    post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

-- 创建索引
CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX idx_post_tags_tag_id ON post_tags(tag_id);

-- 插入预置标签
INSERT INTO tags (name) VALUES 
    ('Python'), ('Go'), ('JavaScript'), ('TypeScript'),
    ('产品经理'), ('设计师'), ('前端'), ('后端'),
    ('AI'), ('数据分析'), ('运营'), ('营销'),
    ('创业'), ('自由职业'), ('远程工作'), ('转型')
ON CONFLICT (name) DO NOTHING;
