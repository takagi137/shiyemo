# spec：邮箱注册 - 成功场景

> 所属模块：认证 (auth)
> 优先级：P0

## 场景描述

新用户使用邮箱 + 密码成功注册账户，注册后自动登录。

## 规格

### Given 用户满足注册条件
- 用户未登录（无有效 token）
- 邮箱 test@example.com 未被注册
- 昵称 "测试用户" 未被使用

### When 用户提交注册请求
- 请求：POST /api/v1/auth/register/email
- 请求体：
  ```json
  {
    "email": "test@example.com",
    "password": "a123456",
    "nickname": "测试用户"
  }
  ```

### Then 注册成功
- 响应状态码：200
- 响应体：
  ```json
  {
    "code": 0,
    "message": "ok",
    "data": {
      "token": "<JWT token>",
      "user": {
        "id": 1,
        "email": "test@example.com",
        "nickname": "测试用户"
      }
    }
  }
  ```
- 数据库验证：
  - `users` 表新增 1 条记录
  - `email` = "test@example.com"
  - `nickname` = "测试用户"
  - `password_hash` 不等于明文 "a123456"
