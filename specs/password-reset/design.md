# design：密码重置

## 技术决策
| 决策点 | 选择 | 理由 |
|--------|------|------|
| 重置方式 | 邮件链接+token | 复用已有邮件服务 |
| token 存储 | Redis，15分钟过期 | 与邮箱验证码一致 |
| 密码加密 | bcrypt | 与注册一致 |

## 模块影响
| 模块 | 改动类型 | 说明 |
|------|---------|------|
| verify/service.go | 新增 | SendResetEmail + ResetPassword |
| verify/handler.go | 新增 | 2 个 handler |
| auth/repo.go | 新增 | UpdatePassword 方法 |
| cmd/main.go | 新增 | 注册 2 个路由 |
| 前端 | 新增 | forgot-password + reset-password 页面 |

## 接口设计
```
POST /api/v1/auth/forgot-password
请求：{ "email": "user@example.com" }
响应：200 { "message": "重置邮件已发送" }

POST /api/v1/auth/reset-password
请求：{ "token": "xxx", "new_password": "newpass123" }
响应：200 { "message": "密码重置成功" }
错误：400 token无效/过期
```

## 数据变更
无新增表，token 存 Redis

## 风险点
- 邮件发送服务依赖 verify 模块已有逻辑
