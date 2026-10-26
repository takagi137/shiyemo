# spec：邮箱登录 - 成功场景

> 所属模块：认证 (auth)
> 优先级：P0

## 场景描述

已注册用户使用邮箱 + 密码成功登录。

## 规格

### Given 用户已注册
- 数据库中存在用户：email = "test@example.com", password_hash = bcrypt("a123456")

### When 用户提交登录请求
- 请求：POST /api/v1/auth/login/email
- 请求体：
  ```json
  {
    "email": "test@example.com",
    "password": "a123456"
  }
  ```

### Then 登录成功
- 响应状态码：200
- 响应体包含 token 和 user 对象
- token 可用于访问需要认证的接口
