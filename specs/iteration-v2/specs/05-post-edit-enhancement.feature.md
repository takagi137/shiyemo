# feature: 帖子编辑增强

## 场景 1：更新帖子标题和内容

### Given 用户已登录，是帖子作者
### When 用户请求 PUT /api/v1/posts/:id 更新标题和内容
### Then 帖子更新成功

```json
// 请求
PUT /api/v1/posts/1
Authorization: Bearer <token>
{
  "title": "更新后的标题",
  "content": "更新后的内容"
}

// 预期响应
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "title": "更新后的标题",
    "content": "更新后的内容"
  }
}
```

## 场景 2：更新帖子标签

### Given 用户已登录，是帖子作者
### When 用户请求 PUT /api/v1/posts/:id 更新标签
### Then 帖子标签更新成功

```json
// 请求
PUT /api/v1/posts/1
Authorization: Bearer <token>
{
  "tags": ["Go", "后端"]
}

// 预期响应
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "tags": ["Go", "后端"]
  }
}
```
