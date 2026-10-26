# 设计：通知功能

## 架构

```
┌─────────────────┐
│   handler.go    │  路由注册、参数解析
├─────────────────┤
│   service.go    │  业务逻辑
├─────────────────┤
│     repo.go     │  数据库 CRUD
├─────────────────┤
│   model.go      │  数据模型
└─────────────────┘
```

## 数据模型

### notifications 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL | 主键 |
| user_id | BIGINT | 接收通知的用户 |
| actor_id | BIGINT | 触发通知的用户 |
| type | VARCHAR(20) | like/comment/follow/reply |
| target_id | BIGINT | 目标ID（帖子ID/评论ID） |
| target_type | VARCHAR(20) | post/comment |
| content | TEXT | 通知内容摘要 |
| is_read | BOOLEAN | 是否已读 |
| created_at | TIMESTAMP | 创建时间 |

### 索引

- idx_notifications_user_id (user_id)
- idx_notifications_user_read (user_id, is_read)
- idx_notifications_created_at (created_at DESC)

## API 设计

### 获取通知列表

```
GET /api/v1/notifications?page=1&page_size=20
Authorization: Bearer <token>
```

响应：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": 1,
        "type": "like",
        "actor": {"id": 2, "nickname": "张三", "avatar": "..."},
        "target": {"id": 10, "type": "post", "title": "..."},
        "content": "赞了你的帖子",
        "is_read": false,
        "created_at": "2026-05-19T10:00:00Z"
      }
    ]
  }
}
```

### 获取未读数量

```
GET /api/v1/notifications/unread-count
Authorization: Bearer <token>
```

响应：
```json
{
  "code": 0,
  "message": "ok",
  "data": {"count": 5}
}
```

### 标记单条已读

```
PUT /api/v1/notifications/:id/read
Authorization: Bearer <token>
```

### 标记全部已读

```
PUT /api/v1/notifications/read-all
Authorization: Bearer <token>
```

## 触发时机

通知由其他业务模块在以下时机创建：
- **点赞**：interaction service 点赞成功后
- **评论**：comment service 创建评论后
- **回复**：comment service 创建回复后
- **关注**：follow service 关注成功后
