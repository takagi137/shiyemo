# 更新帖子标签

## 场景：成功更新帖子标签
- API: `PUT /api/v1/posts/:id/tags`
- 需要登录，且是帖子作者

### 请求
```
PUT /api/v1/posts/1/tags
Authorization: Bearer <token>

{
  "tags": ["Python", "Go", "后端"]
}
```

### 响应 (200)
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "tags": ["Python", "Go", "后端"]
  }
}
```
