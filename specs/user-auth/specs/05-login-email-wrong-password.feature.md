# spec：邮箱登录 - 密码错误

> 所属模块：认证 (auth)
> 优先级：P0

## 场景描述

用户输入错误的密码登录，系统应拒绝并返回明确错误信息。

## 规格

### Given 用户已注册
- 数据库中存在用户：email = "test@example.com", password_hash = bcrypt("a123456")

### When 用户提交登录请求（密码错误）
- 请求：POST /api/v1/auth/login/email
- 请求体：
  ```json
  {
    "email": "test@example.com",
    "password": "wrong_password"
  }
  ```

### Then 登录失败
- 响应状态码：200
- 响应体：
  ```json
  {
    "code": 40204,
    "message": "密码错误",
    "data": null
  }
  ```
