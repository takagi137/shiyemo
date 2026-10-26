> 开源副本：部署地址已替换为示例；历史验证记录不代表当前部署状态。

# tasks：密码重置

## Phase 1：后端

- [x] 1.1 auth/repo.go 新增 UpdatePassword(userID, hashedPassword)
  - 验收：单元测试通过

- [x] 1.2 verify/service.go 新增 SendResetEmail + ResetPassword
  - 复用已有邮件发送逻辑
  - token 存 Redis，15分钟过期
  - 验收：单元测试通过

- [x] 1.3 verify/handler.go 新增 2 个 handler
  - POST /auth/forgot-password
  - POST /auth/reset-password
  - 验收：go build 通过

- [x] 1.4 cmd/main.go 注册路由
  - 验收：go build 通过

## Phase 2：前端

- [x] 2.1 创建 web/src/app/forgot-password/page.tsx
  - 输入邮箱表单
  - 提交后提示"请查收邮件"
  - 验收：npm run build 通过

- [x] 2.2 创建 web/src/app/reset-password/page.tsx
  - 输入新密码表单
  - 带 token 参数（从邮件链接获取）
  - 重置成功后跳转登录页
  - 验收：npm run build 通过

- [x] 2.3 登录页加"忘记密码"链接
  - 验收：点击跳转 /forgot-password

## Phase 3：验证

- [x] 3.1 后端单测通过
- [x] 3.2 前端 build 通过
- [x] 3.3 手动验证完整流程

---

**完成时间**：2026-05-28
**验证结果**：代码已部署到生产环境 http://app.example.com
