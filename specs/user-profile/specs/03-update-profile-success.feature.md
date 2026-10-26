# Feature: 更新用户资料成功

## 场景：更新所有可编辑字段

### 前置条件
- 用户已登录
- 用户 ID=1

### 测试步骤

1. 发送 PUT 请求到 `/api/v1/users/1`
2. 请求头携带 `Authorization: Bearer <token>`
3. 请求体：
```json
{
  "nickname": "新昵称",
  "bio": "这是我的个人简介",
  "company": "新公司",
  "occupation": "新职业"
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

- [x] 更新成功返回 code=0
- [x] 再次获取资料时字段已更新

## 场景：部分更新（只更新昵称）

### 测试步骤

1. 发送 PUT 请求到 `/api/v1/users/1`
2. 请求体：
```json
{
  "nickname": "只改昵称"
}
```

### 预期结果

- HTTP 状态码：200
- 其他字段保持不变

## 场景：更新他人资料失败

### 前置条件
- 用户 A（ID=1）已登录
- 用户 B（ID=2）存在

### 测试步骤

1. 用户 A 发送 PUT 请求到 `/api/v1/users/2`

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 40302,
  "message": "无权操作",
  "data": null
}
```
