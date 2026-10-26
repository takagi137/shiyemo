# Feature: 获取当前用户资料成功

## 场景：已登录用户获取自己的资料

### 前置条件
- 用户已注册并登录
- 持有有效的 JWT token

### 测试步骤

1. 发送 GET 请求到 `/api/v1/auth/profile`
2. 请求头携带 `Authorization: Bearer <token>`

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "email": "test@example.com",
    "phone": "13800138000",
    "nickname": "测试用户",
    "avatar": "",
    "bio": "",
    "company": "",
    "occupation": "",
    "status": 1,
    "created_at": "2026-05-19T16:13:07Z",
    "updated_at": "2026-05-19T16:13:07Z"
  }
}
```

### 验证点

- [x] 返回完整的用户信息
- [x] 包含手机号和邮箱
- [x] 包含所有可编辑字段
