# tasks：用户注册/登录功能

## Phase 1：数据层

- [x] 任务 1.1：创建 users 表 migration
  - 文件：`server/migrations/001_create_users.sql`
  - 验证：SQL 语法正确，可执行
  - Spec：无

- [x] 任务 1.2：完善 user model
  - 文件：`server/internal/auth/model.go`
  - 验证：结构体字段与 migration 一致
  - Spec：无

- [x] 任务 1.3：实现 auth repo（CreateUser, GetByEmail, GetByPhone）
  - 文件：`server/internal/auth/repo.go`
  - 验证：编译通过
  - Spec：无

## Phase 2：业务层

- [x] 任务 2.1：完善 JWT 工具函数
  - 文件：`server/internal/common/jwt/jwt.go`
  - 验证：GenerateToken, ParseToken 函数可用
  - Spec：无

- [x] 任务 2.2：实现 auth service
  - 文件：`server/internal/auth/service.go`
  - 验证：Register, Login, SendSMS 方法实现
  - Spec：
    - 01-register-email-success.feature.md
    - 03-register-phone-success.feature.md
    - 04-login-email-success.feature.md

## Phase 3：接口层

- [x] 任务 3.1：完善 auth handler
  - 文件：`server/internal/auth/handler.go`
  - 验证：
    - POST /api/v1/auth/register/email 可调用
    - POST /api/v1/auth/register/phone 可调用
    - POST /api/v1/auth/login/email 可调用
    - POST /api/v1/auth/login/phone 可调用
    - POST /api/v1/auth/sms/send 可调用
  - Spec：全部 spec 文件

- [x] 任务 3.2：完善 JWT 中间件
  - 文件：`server/internal/common/middleware/jwt.go`
  - 验证：protected 路由需要 token 才能访问
  - Spec：无

- [x] 任务 3.3：完善 main.go 路由注册
  - 文件：`server/cmd/main.go`
  - 验证：所有注册接口可访问
  - Spec：无

## Phase 4：验证

- [x] 任务 4.1：运行全部 BDD 场景
  - 验证：通过率 100%
  - 方式：curl 或 Postman 手动测试
  - 状态：✅ 2026-05-20 生产环境验证通过

- [x] 任务 4.2：人工验收测试
  - 验证：
    - 邮箱注册 → 登录 → 访问受保护接口
    - 手机号注册 → 登录 → 访问受保护接口
    - 重复注册返回正确错误
    - 密码错误返回正确错误
  - 状态：✅ 2026-05-20 生产环境验收通过

## 依赖关系

```
Phase 1 (数据层) 
    ↓
Phase 2 (业务层) 
    ↓
Phase 3 (接口层) 
    ↓
Phase 4 (验证)
```

每个 Phase 内的任务可以并行，但 Phase 之间有依赖。

## 完成时间

- Phase 1-3 完成时间：2026-05-19
- Phase 4 完成时间：2026-05-20
- 编译状态：✅ 通过
- 生产环境验证：✅ 通过
