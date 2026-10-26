# Feature: 关注用户

## 场景：关注用户

### 前置条件
- 用户 A（ID=1）已登录
- 用户 B（ID=2）存在

### 测试步骤

1. 发送 POST 请求到 `/api/v1/users/2/follow`
2. 请求头携带 `Authorization: Bearer <token>`

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 0,
  "message": "ok",
  "data": { "following": true }
}
```

### 验证点

- [x] 关注成功
- [x] 用户 A 的 following_count 增加
- [x] 用户 B 的 follower_count 增加

## 场景：重复关注（幂等）

### 测试步骤

1. 再次关注同一用户

### 预期结果

- 返回成功

## 场景：取消关注

### 测试步骤

1. 再次发送 POST 请求

### 预期结果

- 返回 { "following": false }

## 场景：不能关注自己

### 测试步骤

1. 关注自己（ID=1）

### 预期结果

- 返回错误：不能关注自己
