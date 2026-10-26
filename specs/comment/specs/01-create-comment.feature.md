# Feature: 发表评论

## 场景：发表主评论

### 前置条件
- 用户已登录（ID=1）
- 帖子存在（ID=1）

### 测试步骤

1. 发送 POST 请求到 `/api/v1/posts/1/comments`
2. 请求头携带 `Authorization: Bearer <token>`
3. 请求体：
```json
{
  "content": "这个帖子写得很好！"
}
```

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
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
    }
  }
}
```

### 验证点

- [x] 返回完整评论信息
- [x] parent_id 为 null（主评论）
- [x] 包含作者信息

## 场景：回复评论

### 前置条件
- 主评论存在（ID=1）

### 测试步骤

1. 发送 POST 请求到 `/api/v1/posts/1/comments`
2. 请求体：
```json
{
  "content": "同意你的看法",
  "parent_id": 1
}
```

### 预期结果

- 创建成功
- parent_id 为 1

## 场景：未登录评论失败

### 测试步骤

1. 不携带 token 发送 POST 请求

### 预期结果

- HTTP 状态码：401
