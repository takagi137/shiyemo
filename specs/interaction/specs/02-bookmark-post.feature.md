# Feature: 收藏帖子

## 场景：收藏帖子

### 前置条件
- 用户已登录（ID=1）
- 帖子存在（ID=1）

### 测试步骤

1. 发送 POST 请求到 `/api/v1/posts/1/bookmark`
2. 请求头携带 `Authorization: Bearer <token>`

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 0,
  "message": "ok",
  "data": { "bookmarked": true }
}
```

### 验证点

- [x] 收藏成功

## 场景：重复收藏（幂等）

### 测试步骤

1. 再次收藏同一帖子

### 预期结果

- 返回成功

## 场景：取消收藏

### 测试步骤

1. 再次发送 POST 请求

### 预期结果

- 返回 { "bookmarked": false }
