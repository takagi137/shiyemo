# Feature: 获取指定用户公开资料

## 场景：获取其他用户的公开资料

### 前置条件
- 目标用户存在（ID=1）

### 测试步骤

1. 发送 GET 请求到 `/api/v1/users/1`

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "nickname": "测试用户",
    "avatar": "",
    "bio": "",
    "company": "",
    "occupation": "",
    "created_at": "2026-05-19T16:13:07Z"
  }
}
```

### 验证点

- [x] 不返回 email 和 phone（隐私保护）
- [x] 不返回 status 字段
- [x] 返回公开信息

## 场景：获取不存在的用户

### 测试步骤

1. 发送 GET 请求到 `/api/v1/users/99999`

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 40203,
  "message": "用户不存在",
  "data": null
}
```
