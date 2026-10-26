# spec：重置密码 - token 过期

> 所属模块：认证 (auth)
> 优先级：P1

## 场景描述
用户使用过期的重置 token 尝试重置密码。

## 规格

### Given 重置 token 已过期
- Redis 中不存在该 token（已过期被删除）

### When 提交重置密码请求
- POST /api/v1/auth/reset-password
- 请求体：{ "token": "expired_token", "new_password": "newpass123" }

### Then 重置失败
- 响应状态码：400
- 错误信息："重置链接已过期，请重新申请"
