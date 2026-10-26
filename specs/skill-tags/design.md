# 技能标签功能设计

## 数据模型

### tags 表
```sql
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### post_tags 表
```sql
CREATE TABLE post_tags (
    post_id INT REFERENCES posts(id),
    tag_id INT REFERENCES tags(id),
    PRIMARY KEY (post_id, tag_id)
);
```

### users 表扩展
```sql
ALTER TABLE users ADD COLUMN tags TEXT;  -- JSON数组存储
```

## API 设计

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v1/tags` | GET | 获取热门标签 |
| `/api/v1/posts/:id/tags` | PUT | 更新帖子标签 |
| `/api/v1/posts?tag=xxx` | GET | 按标签筛选 |

## 实现方案

1. 标签预置 + 用户自定义
2. 帖子标签通过关联表存储
3. 用户标签存储在 users.tags 字段（JSON）
