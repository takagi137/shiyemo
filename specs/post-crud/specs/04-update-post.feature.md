# Feature: 更新帖子

## 场景：作者更新自己的帖子

### 前置条件
- 用户 A（ID=1）已登录
- 用户 A 的帖子存在（ID=1）

### 测试步骤

1. 发送 PUT 请求到 `/api/v1/posts/1`
2. 请求头携带 `Authorization: Bearer <token>`
3. 请求体：
```json
{
  "title": "新标题",
  "content": "新内容"
}
```

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

- [x] 更新成功
- [x] 再次获取时字段已更新

## 场景：更新他人帖子失败

### 前置条件
- 用户 A（ID=1）已登录
- 用户 B 的帖子存在（ID=2）

### 测试步骤

1. 用户 A 发送 PUT 请求到 `/api/v1/posts/2`

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
