# tasks：Sprint 1 — 邮箱验证 + 自由天数 + 生存贴士

## Phase 1：数据层 + 公共设施

### Task 1.1：User 模型新增字段
- 文件：`server/internal/auth/model.go`
- 改动：新增 `EmailVerified bool` + `LastWorkDate *time.Time`
- 验收：GORM AutoMigrate 执行成功，表结构更新

### Task 1.2：创建 SMTP 邮件发送工具
- 文件：`server/common/email/sender.go`
- 改动：实现 `SendVerificationEmail(to, token string) error`
- 验收：单元测试通过，可发送邮件

### Task 1.3：Config 新增 SMTP 配置
- 文件：`server/common/config/config.go`
- 改动：新增 Email 配置段（Host/Port/User/Pass/From）
- 验收：配置加载成功

---

## Phase 2：生存贴士模块（独立，无依赖）

### Task 2.1：Tips model + repo
- 文件：`server/internal/tips/model.go`、`server/internal/tips/repo.go`
- 改动：SurvivalTip 模型 + GetRandom / GetList / Create
- 验收：repo_test.go 通过

### Task 2.2：Tips service
- 文件：`server/internal/tips/service.go`
- 改动：GetRandomTip / ListTips / CreateTip
- 验收：service_test.go 通过（≥80% 覆盖率）

### Task 2.3：Tips handler + 路由注册
- 文件：`server/internal/tips/handler.go`、`server/cmd/main.go`
- 改动：handler 实现 + 路由 `GET /tips/random`、`GET /tips`、`POST /tips`
- 验收：curl 测试通过

### Task 2.4：Tips 种子数据
- 文件：`server/internal/tips/seed.go`
- 改动：初始 30 条中文贴士（finance/life/mindset/policy 四类）
- 验收：启动后 `GET /tips/random` 返回数据

---

## Phase 3：邮箱验证模块

### Task 3.1：Verify service
- 文件：`server/internal/verify/service.go`
- 改动：SendVerification / ConfirmVerification
- 验收：service_test.go 通过

### Task 3.2：Verify handler + 路由
- 文件：`server/internal/verify/handler.go`、`server/cmd/main.go`
- 改动：`POST /auth/verify/send`（需登录）、`GET /auth/verify/:token`（公开）
- 验收：BDD 场景 01-email-verify 通过

### Task 3.3：注册流程集成验证邮件
- 文件：`server/internal/auth/service.go`
- 改动：RegisterByEmail 成功后异步发送验证邮件
- 验收：注册后 Redis 中有验证 Token

---

## Phase 4：自由天数（后端）

### Task 4.1：User Profile 接口扩展
- 文件：`server/internal/user/service.go`、`server/internal/user/handler.go`
- 改动：UpdateProfile 支持 last_work_date 字段
- 验收：PUT /user/profile 可更新 last_work_date

### Task 4.2：User 详情返回新字段
- 文件：`server/internal/user/service.go`
- 改动：GetUserProfile 返回 last_work_date + email_verified
- 验收：GET /users/:id 响应包含新字段

---

## Phase 5：前端 — 生存贴士 + 自由天数

### Task 5.1：DailyTip 组件
- 文件：`web/src/app/components/DailyTip.tsx`
- 改动：调用 `/tips/random`，展示卡片，支持刷新
- 验收：首页可见贴士卡片

### Task 5.2：FreedomDayCard 组件
- 文件：`web/src/app/components/FreedomDayCard.tsx`
- 改动：读取 user.last_work_date，计算并展示"自由第 X 天"
- 验收：有数据时展示，无数据时隐藏

### Task 5.3：首页集成两个组件
- 文件：`web/src/app/page.tsx`
- 改动：首页布局中加入 DailyTip + FreedomDayCard
- 验收：首页完整展示

---

## Phase 6：前端 — 邮箱验证

### Task 6.1：验证邮箱提示页
- 文件：`web/src/app/verify-email/page.tsx`
- 改动：提示用户去邮箱验证 + 重新发送按钮
- 验收：页面渲染正常、重发按钮调用接口

### Task 6.2：验证确认页
- 文件：`web/src/app/verify-email/confirm/page.tsx`
- 改动：从 URL 读取 token，调用验证接口，展示结果
- 验收：有效 token 显示成功，无效显示失败

### Task 6.3：个人资料编辑增加"最后工作日"
- 文件：`web/src/app/profile/page.tsx`（或 ProfileEditForm）
- 改动：新增日期选择器字段
- 验收：可选择日期并保存

---

## Phase 7：部署 + 集成验证

### Task 7.1：docker-compose 新增 SMTP 环境变量
- 文件：`deploy/docker-compose.yml`
- 改动：API 服务新增 EMAIL_* 环境变量
- 验收：容器启动正常

### Task 7.2：部署到服务器
- 操作：scp + docker-compose build + up
- 验收：远程 `GET /tips/random` 返回数据

### Task 7.3：L3 全流程验证
- 操作：按 specs/01-03 逐条验证
- 验收：所有 BDD 场景通过

---

## 执行约束提醒

- ⚠️ 单次任务最多改 3 个文件
- ⚠️ 单个函数不超过 50 行
- ⚠️ 单个文件不超过 300 行
- ⚠️ handler 不写业务逻辑
- ⚠️ 每个 service 必须有 service_test.go
- ⚠️ 编码前回答三问：改什么？能否更小？怎么验证？

---

*状态：待执行*
*预计工作量：Task 1.1 ~ 7.3 共 17 个子任务*
