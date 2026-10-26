# design：Phase 2 Sprint 1 — 邮箱验证 + 自由天数 + 生存贴士

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 邮件发送方式 | SMTP（net/smtp 标准库） | Phase1 不引入三方 SDK，可用免费 SMTP 服务 |
| 验证 Token 生成 | crypto/rand 生成 32 字节 hex | 安全且无依赖 |
| Token 存储 | Redis（TTL 24h） | 简单、自动过期、无需新表 |
| 自由天数计算 | 前端实时计算（不存服务端） | 减少后端开销，字段只存 last_work_date |
| 贴士存储 | PostgreSQL 表 + Redis 缓存 | 支持后台管理 + 高频读取 |
| 贴士随机策略 | 按用户 ID + 日期做种子 | 同一用户同一天看到同一条，刷新换一条 |

## 模块影响

| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| auth | 修改 | model.go 新增 EmailVerified 字段；注册后触发验证邮件 |
| user | 修改 | service.go 支持 last_work_date 字段更新与返回 |
| verify（新） | 新增 | 邮箱验证完整模块 |
| tips（新） | 新增 | 生存贴士完整模块 |
| common/email（新） | 新增 | SMTP 邮件发送工具 |
| common/config | 修改 | 新增 SMTP 配置项 |

## 接口设计

### 邮箱验证

```
POST /api/v1/auth/verify/send
Headers: Authorization: Bearer <token>
请求：{}
响应：{ "code": 0, "message": "验证邮件已发送" }
错误：40101 未登录 / 40201 邮箱已验证 / 50002 发送失败

GET /api/v1/auth/verify/:token
请求：无（URL 参数）
响应：{ "code": 0, "message": "邮箱验证成功" }
错误：40001 Token 无效或已过期
```

### 用户资料（扩展）

```
PUT /api/v1/user/profile
Headers: Authorization: Bearer <token>
请求（新增字段）：{ "last_work_date": "2026-03-15" }
响应：{ "code": 0, "data": { ...user, "last_work_date": "2026-03-15" } }

GET /api/v1/users/:id
响应新增字段：{ "last_work_date": "2026-03-15", "email_verified": true }
```

### 生存贴士

```
GET /api/v1/tips/random
请求：无
响应：{ "code": 0, "data": { "id": 1, "content": "...", "category": "finance" } }

GET /api/v1/tips  (管理员)
请求：?page=1&size=20
响应：{ "code": 0, "data": { "list": [...], "page": 1, "size": 20, "has_more": false } }

POST /api/v1/tips  (管理员)
请求：{ "content": "...", "category": "finance" }
响应：{ "code": 0, "data": { "id": 1 } }
```

## 数据变更

| 操作 | 对象 | 说明 |
|------|------|------|
| ALTER TABLE | users | 新增 `email_verified BOOLEAN DEFAULT FALSE` |
| ALTER TABLE | users | 新增 `last_work_date DATE` |
| CREATE TABLE | survival_tips | 贴士表（id, content, category, is_active, created_at） |
| SEED DATA | survival_tips | 初始 20 条贴士 |

## 风险点

| 风险 | 影响 | 缓解 |
|------|------|------|
| SMTP 服务不稳定 | 验证邮件发不出 | 加重试机制 + 支持重新发送 |
| 邮件进垃圾箱 | 用户收不到 | 前端提示"检查垃圾箱" |
| last_work_date 不填 | 自由天数无法计算 | 前端判断无数据时不展示 |
| 贴士数量少时重复 | 用户体验差 | 初始准备 30+ 条 |

## 前端页面变更

| 页面 | 改动 |
|------|------|
| 注册成功后 | 跳转"请验证邮箱"提示页 |
| 首页 | 新增"自由第 X 天"卡片 + "今日贴士"卡片 |
| 个人资料编辑 | 新增"最后工作日"日期选择 |
| 发帖页 | 未验证邮箱时弹窗提示（非阻断） |
