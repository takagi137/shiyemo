# Feature: 删除帖子

## 场景：作者删除自己的帖子

### 前置条件
- 用户 A（ID=1）已登录
- 用户 A 的帖子存在（ID=1）

### 测试步骤

1. 发送 DELETE 请求到 `/api/v1/posts/1`
2. 请求头携带 `Authorization: Bearer <token>`

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 0,
  "message": "ok",
  "data": null
}
```

### 验证点

- [x] 删除成功
- [x] 再次获取时返回 40301

## 场景：删除他人帖子失败

### 前置条件
- 用户 A（ID=1）已登录
- 用户 B 的帖子存在（ID=2）

### 测试步骤

1. 用户 A 发送 DELETE 请求到 `/api/v1/posts/2`

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 40302,
  "message": "无权操作此帖子",
  "data": null
}
```
