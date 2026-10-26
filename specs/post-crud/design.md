# design：发帖功能

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 帖子类型 | 枚举 | emotion/insight/case/action |
| 图片存储 | JSONB 数组 | MVP 阶段存 URL，后续迁移到 OSS |
| 标签存储 | JSONB 数组 | 灵活查询 |
| 删除方式 | 软删除 | 保留数据，支持恢复 |
| 排序方式 | 时间/热度 | latest/hot |

## 接口设计

### 1. 创建帖子

```
POST /api/v1/posts
Authorization: Bearer <token>
Content-Type: application/json

请求：
{
  "type": "emotion",
  "title": "失业第30天，我开始焦虑了",
  "content": "今天是失业的第30天...",
  "images": ["https://example.com/img1.jpg"],
  "tags": ["失业", "焦虑", "心情"]
}

响应：
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "user_id": 1,
    "type": "emotion",
    "title": "失业第30天，我开始焦虑了",
    "content": "今天是失业的第30天...",
    "images": ["https://example.com/img1.jpg"],
    "tags": ["失业", "焦虑", "心情"],
    "like_count": 0,
    "comment_count": 0,
    "status": 1,
    "created_at": "2026-05-19T16:30:00Z",
    "updated_at": "2026-05-19T16:30:00Z"
  }
}
```

### 2. 获取帖子详情

```
GET /api/v1/posts/:id

响应：
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "user_id": 1,
    "type": "emotion",
    "title": "失业第30天，我开始焦虑了",
    "content": "今天是失业的第30天...",
    "images": ["https://example.com/img1.jpg"],
    "tags": ["失业", "焦虑", "心情"],
    "like_count": 5,
    "comment_count": 3,
    "status": 1,
    "created_at": "2026-05-19T16:30:00Z",
    "updated_at": "2026-05-19T16:30:00Z",
    "author": {
      "id": 1,
      "nickname": "测试用户",
      "avatar": "/uploads/avatars/1.jpg"
    }
  }
}
```

### 3. 获取帖子列表

```
GET /api/v1/posts?type=emotion&sort=latest&page=1&size=20

响应：
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [...],
    "page": 1,
    "size": 20,
    "has_more": true
  }
}
```

### 4. 更新帖子

```
PUT /api/v1/posts/:id
Authorization: Bearer <token>

请求：
{
  "title": "新标题",
  "content": "新内容"
}

响应：
{
  "code": 0,
  "message": "ok",
  "data": null
}
```

### 5. 删除帖子

```
DELETE /api/v1/posts/:id
Authorization: Bearer <token>

响应：
{
  "code": 0,
  "message": "ok",
  "data": null
}
```

## 数据变更

| 操作 | 对象 | 说明 |
|------|------|------|
| 新增 | posts 表 | 创建帖子表 |

## 依赖

- 已有 PostgreSQL：帖子数据存储
- 已有 User 模型：关联作者信息

## 风险点

1. 帖子内容长度限制（建议 10000 字符）
2. 图片 URL 格式校验
3. 标签数量限制（建议 10 个）
