# design：评论功能

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 评论层级 | 最多两级 | 简化设计，避免过深嵌套 |
| 删除方式 | 软删除 | 保留数据，支持恢复 |
| 排序方式 | 时间正序 | 符合阅读习惯 |

## 接口设计

### 1. 发表评论

```
POST /api/v1/posts/:id/comments
Authorization: Bearer <token>
Content-Type: application/json

请求：
{
  "content": "这个帖子写得很好！",
  "parent_id": null  // 主评论，不传或传 null
}

或回复评论：
{
  "content": "同意你的看法",
  "parent_id": 1  // 回复 ID=1 的评论
}
```

### 2. 获取帖子评论列表

```
GET /api/v1/posts/:id/comments?page=1&size=20

响应：
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "post_id": 1,
        "user_id": 1,
        "parent_id": null,
        "content": "这个帖子写得很好！",
        "like_count": 0,
        "created_at": "2026-05-19T16:30:00Z",
        "author": {
          "id": 1,
          "nickname": "测试用户",
          "avatar": "/uploads/avatars/1.jpg"
        },
        "replies": [
          {
            "id": 2,
            "post_id": 1,
            "user_id": 2,
            "parent_id": 1,
            "content": "同意你的看法",
            "like_count": 0,
            "created_at": "2026-05-19T16:35:00Z",
            "author": {
              "id": 2,
              "nickname": "用户2",
              "avatar": ""
            }
          }
        ]
      }
    ],
    "page": 1,
    "size": 20,
    "has_more": false
  }
}
```

### 3. 删除评论

```
DELETE /api/v1/comments/:id
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
| 新增 | comments 表 | 创建评论表 |

## 依赖

- 已有 PostgreSQL：评论数据存储
- 已有 Post 模型：关联帖子
- 已有 User 模型：关联作者

## 风险点

1. 评论内容长度限制（建议 1000 字符）
2. 评论层级限制（最多两级）
3. 评论计数更新（帖子的 comment_count）
