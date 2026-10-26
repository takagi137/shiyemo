# feature: 按类型筛选通知

## 场景 1：按类型筛选通知

### Given 用户已登录，存在点赞类型通知
### When 用户请求 GET /api/v1/notifications?type=like
### Then 只返回点赞类型通知

```json
// 请求
GET /api/v1/notifications?type=like&page=1&page_size=10
Authorization: Bearer <token>

// 预期响应
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": 1,
        "type": "like",
        "content": "张三 点赞了你的帖子"
      }
    ],
    "total": 1
  }
}
```

## 场景 2：批量标记已读

### Given 用户已登录，存在未读通知 [1, 2, 3]
### When 用户请求 POST /api/v1/notifications/batch/read { ids: [1, 2, 3] }
### Then 通知 1, 2, 3 标记为已读

```json
// 请求
POST /api/v1/notifications/batch/read
Authorization: Bearer <token>
{
  "ids": [1, 2, 3]
}

// 预期响应
{
  "code": 0,
  "message": "ok"
}
```

## 场景 3：批量删除通知

### Given 用户已登录，存在通知 [1, 2, 3]
### When 用户请求 POST /api/v1/notifications/batch/delete { ids: [1, 2, 3] }
### Then 通知 1, 2, 3 被删除

```json
// 请求
POST /api/v1/notifications/batch/delete
Authorization: Bearer <token>
{
  "ids": [1, 2, 3]
}

// 预期响应
{
  "code": 0,
  "message": "ok"
}
```
