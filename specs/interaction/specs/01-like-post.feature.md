# Feature: 点赞帖子

## 场景：点赞帖子

### 前置条件
- 用户已登录（ID=1）
- 帖子存在（ID=1）

### 测试步骤

1. 发送 POST 请求到 `/api/v1/posts/1/like`
2. 请求头携带 `Authorization: Bearer <token>`

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 0,
  "message": "ok",
  "data": { "liked": true }
}
```

### 验证点

- [x] 点赞成功
- [x] 帖子 like_count 增加

## 场景：重复点赞（幂等）

### 测试步骤

1. 再次点赞同一帖子

### 预期结果

- 返回成功
- like_count 不变

## 场景：取消点赞

### 测试步骤

1. 再次发送 POST 请求

### 预期结果

- 返回 { "liked": false }
- like_count 减少

## 场景：未登录点赞失败

### 测试步骤

1. 不携带 token 发送请求

### 预期结果

- HTTP 状态码：401
