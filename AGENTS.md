# AGENTS.md - AI 开发助手必读

## 项目概述

**失业么** — AI 时代个体的"能力重建与收入重建"社区。

技术栈：Go (Gin) + Next.js + PostgreSQL + Redis + Meilisearch

## 🔴 红线（绝对禁止）

1. 不得引入 Java/Python 作为后端
2. 不得使用除 GORM 外的其他 ORM
3. 不得在 Phase1 拆分微服务
4. 不得在 MVP 阶段实现即时通讯/推荐算法/支付系统
5. 不得修改已有的数据库迁移文件
6. 不得在 handler 层写业务逻辑
7. 不得直接拼接 SQL 字符串

## 📋 编码前三问

每次编码前必须回答：

1. **改什么？** 明确功能模块 + 涉及文件列表
2. **能否更小？** 单次任务不超过 3 个文件
3. **怎么验证？** 至少有一个场景可通过

## 🏗️ 分层规范

```
handler → service → repo
```

| 层 | 职责 | 禁止 |
|----|------|------|
| handler | 参数校验、调用 service、返回响应 | 写业务逻辑 |
| service | 业务逻辑、事务编排 | 依赖 HTTP 框架 |
| repo | 数据库 CRUD | 写业务规则 |

## 📦 内容模型

- 帖子类型：emotion / insight / case / action / achievement
  （achievement「成果展示」于 2026-07-20 补入后端 enum：前端 constants.ts 早已提供该选项，
   但后端校验未包含，用户选「成果展示」发帖会被拒。前后端 enum 需保持同步。）
- Phase1 不实现 action 扩展
- 评论最多两级
- 排序：latest / hot

## 🔧 API 规范

统一响应格式：
```json
{ "code": 0, "message": "ok", "data": {} }
```

错误码：400xx 参数错误，401xx 认证错误，402xx 用户错误，403xx 帖子错误，500xx 服务错误

## 🎯 验证优先原则（强制）

**核心理念：先写验证，再写代码；验证不过，代码不入库。**

### 验证三层架构

| 层级 | 名称 | 覆盖范围 | 工具 | 时机 |
|------|------|---------|------|------|
| L1 | 单元测试 | 每个函数/方法 | `go test` + `testify` | 开发中 |
| L2 | 集成测试 | API 端点 | `httptest` + 真实 DB | 开发后 |
| L3 | 部署验证 | 完整流程 | curl + BDD 场景 | 部署后 |

### 验证检查清单（每个功能完成前必须确认）

- [ ] L1：单元测试通过（`go test ./...`）
- [ ] L1：测试覆盖率达标（service ≥ 80%）
- [ ] L2：集成测试通过
- [ ] L3：部署到验证环境成功
- [ ] L3：BDD 场景全部通过
- [ ] L3：关键路径手动验证通过

### 文件要求

```
internal/
├── auth/
│   ├── service.go
│   ├── service_test.go  ← 必须
│   ├── repo.go
│   └── repo_test.go     ← 必须
```

## 📚 参考文档

| 文档 | 路径 |
|------|------|
| 架构设计 | docs/01-架构/01-architecture.md |
| 技术栈 | docs/01-架构/00-tech-stack.md |
| OpenSpec 流程 | docs/04-开发流程/00-openspec-guide.md |
| BDD 场景写法 | docs/04-开发流程/01-bdd-guide.md |
| **验证指南** | docs/04-开发流程/02-verification-guide.md |
| AI 开发约束 | docs/02-开发规范/02-ai-development-guide.md |
