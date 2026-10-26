# spec：邮箱注册 - 邮箱已存在

> 所属模块：认证 (auth)
> 优先级：P0

## 场景描述

用户尝试使用已注册的邮箱注册，系统应拒绝并返回明确错误信息。

## 规格

### Given 邮箱已被注册
- 数据库中已存在用户：email = "test@example.com"

### When 用户提交注册请求
- 请求：POST /api/v1/auth/register/email
- 请求体：
  ```json
  {
    "email": "test@example.com",
    "password": "a123456",
    "nickname": "新用户"
  }
  ```

### Then 注册失败
- 响应状态码：200
- 响应体：
  ```json
  {
    "code": 40201,
    "message": "邮箱已注册",
    "data": null
  }
  ```
- 数据库验证：
  - `users` 表记录数不变
