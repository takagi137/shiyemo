# Feature: 获取帖子详情

## 场景：获取存在的帖子详情

### 前置条件
- 帖子存在（ID=1）

### 测试步骤

1. 发送 GET 请求到 `/api/v1/posts/1`

### 预期结果

- HTTP 状态码：200
- 响应体包含完整帖子信息和作者信息

### 验证点

- [x] 返回帖子所有字段
- [x] 包含作者信息（nickname, avatar）
- [x] 包含 like_count 和 comment_count

## 场景：获取不存在的帖子

### 测试步骤

1. 发送 GET 请求到 `/api/v1/posts/99999`

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 40301,
  "message": "帖子不存在",
  "data": null
}
```
