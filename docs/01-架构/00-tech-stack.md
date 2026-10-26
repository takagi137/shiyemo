# 技术栈选型

> 本文档定义"失业么"项目的技术栈选型及版本锁定。
> 所有 AI 生成的代码必须遵循本文档的技术栈约束。

---

## 技术栈总览

| 层级 | 选型 | 版本 | 选型理由 |
|------|------|------|---------|
| **前端** | Next.js | 14+ | SSR/SEO 友好，AI Coding 效率高 |
| **前端语言** | TypeScript | 5.x | 类型安全，AI 更擅长 |
| **前端样式** | TailwindCSS | 3.x | 实用优先，快速开发 |
| **后端语言** | Go | 1.21+ | 性能优秀，编译快，单人友好 |
| **后端框架** | Gin | latest | 轻量高性能，生态成熟 |
| **ORM** | GORM | latest | Go 生态主流，文档完善 |
| **数据库** | PostgreSQL | 15+ | JSONB 支持，全文搜索，稳定可靠 |
| **缓存** | Redis | 7+ | 轻量高性能，功能丰富 |
| **搜索** | Meilisearch | latest | 轻量全文搜索，中文分词好 |
| **容器化** | Docker | 24+ | 标准化部署 |
| **编排** | Docker Compose | latest | Phase1 本地开发 & 部署 |
| **反向代理** | Caddy | latest | 自动 HTTPS，配置简单 |
| **对象存储** | MinIO / Aliyun OSS | - | 图片存储，Phase1 可用文件系统 |

---

## 前端技术栈详情

### 核心依赖

| 包名 | 用途 | 版本锁定 |
|------|------|---------|
| next | 框架 | ^14.0.0 |
| react / react-dom | UI 库 | ^18.2.0 |
| typescript | 类型系统 | ^5.3.0 |
| tailwindcss | CSS 框架 | ^3.4.0 |

### 认证相关

| 包名 | 用途 | 版本锁定 |
|------|------|---------|
| next-auth（或自实现） | 认证 | - |
| jose | JWT 处理 | ^5.0.0 |

### UI 相关

| 包名 | 用途 | 版本锁定 |
|------|------|---------|
| @headlessui/react | 无样式组件 | ^2.0.0 |
| react-hook-form | 表单管理 | ^7.50.0 |
| zod | 表单校验 | ^3.22.0 |
| react-hot-toast | 提示 | ^2.4.0 |

### 数据请求

| 包名 | 用途 | 版本锁定 |
|------|------|---------|
| swr / react-query | 数据获取 | - |

### 编辑器（可选）

| 包名 | 用途 | 版本锁定 |
|------|------|---------|
| @tiptap/react | 富文本编辑器 | ^2.0.0 |

---

## 后端技术栈详情

### go.mod 示例

```go
module github.com/shiyemo/server

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.4
    github.com/redis/go-redis/v9 v9.3.0
    github.com/meilisearch/meilisearch-go v0.26.0
)
```

### 核心依赖

| 包名 | 用途 | 说明 |
|------|------|------|
| github.com/gin-gonic/gin | HTTP 框架 | 轻量高性能 |
| github.com/golang-jwt/jwt/v5 | JWT | Token 生成与验证 |
| gorm.io/gorm | ORM | 数据库操作 |
| gorm.io/driver/postgres | PostgreSQL 驱动 | GORM 配套 |
| github.com/redis/go-redis/v9 | Redis 客户端 | 缓存操作 |
| github.com/meilisearch/meilisearch-go | Meilisearch 客户端 | 搜索 |
| github.com/google/uuid | UUID | ID 生成 |
| golang.org/x/crypto | 加密 | 密码哈希 |

### 密码处理

```go
import "golang.org/x/crypto/bcrypt"

// 哈希密码
hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 验证密码
err := bcrypt.CompareHashAndPassword(hashed, []byte(password))
```

---

## 工具链

### 开发工具

| 工具 | 用途 | 安装方式 |
|------|------|---------|
| gofmt | 代码格式化 | Go 自带 |
| golangci-lint | 代码检查 | `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` |
| migrate | 数据库迁移 | `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest` |
| air | 热重载 | `go install github.com/air-verse/air@latest` |

### 前端工具

| 工具 | 用途 | 安装方式 |
|------|------|---------|
| pnpm | 包管理 | `npm install -g pnpm` |
| prettier | 代码格式化 | 项目内安装 |
| eslint | 代码检查 | 项目内安装 |

---

## 版本锁定原则

| 原则 | 说明 |
|------|------|
|锁定大版本 | 使用 `^14.0.0` 而非 `14.0.0`，允许 patch 更新 |
| 关键依赖锁定 | 认证、数据库相关锁定 minor 版本 |
| 前端使用 pnpm-lock.yaml | 确保团队和 CI 环境一致 |
| 后端使用 go.sum | Go 原生机制 |

---

## 为什么不选其他技术栈

### 后端替代方案

| 方案 | 不选原因 |
|------|---------|
| Java + Spring Boot | 重量级，单人开发效率低，部署包大 |
| Python + FastAPI | 动态类型，AI 生成代码质量不如静态类型语言 |
| Node.js + Express | 单线程，CPU 密集型场景性能差 |
| Rust + Actix | 学习曲线陡，编译慢，生态不如 Go 成熟 |

### 前端替代方案

| 方案 | 不选原因 |
|------|---------|
| Vue 3 + Nuxt | SSR 配置复杂，React Native 复用性差 |
| Svelte + SvelteKit | 生态尚不成熟，组件库少 |
| React + Vite（纯 SPA） | SEO 不友好，内容社区需要 SEO |

### ORM 替代方案

| 方案 | 不选原因 |
|------|---------|
| sqlx | 需要手写 SQL，维护成本高 |
| ent | 代码生成复杂，单人项目过度设计 |
| Beego ORM | 社区活跃度下降 |

### 数据库替代方案

| 方案 | 不选原因 |
|------|---------|
| MySQL | JSONB 支持弱于 PostgreSQL，全文搜索能力差 |
| MongoDB | 不适合关系型数据，事务支持弱 |
| SQLite | 并发性能差，不适合 40 万用户场景 |

---

## 部署环境要求

### 开发环境

| 工具 | 版本要求 |
|------|---------|
| Go | 1.21+ |
| Node.js | 20+ |
| Docker Desktop | 4.x |
| Git | 2.x |

### 生产环境（Phase1）

| 组件 | 最低配置 |
|------|---------|
| 云服务器 | 2 核 4G |
| 云数据库 | PostgreSQL 15, 2 核 4G（可与应用同机） |
| 云 Redis | 256MB |
| 域名 + HTTPS | Caddy 自动证书 |

---

## 相关文档

| 文档 | 路径 |
|------|------|
| 整体架构 | `docs/01-架构/01-architecture.md` |
| 数据模型 | `docs/01-架构/02-data-model.md` |
| 部署指南 | `docs/01-架构/04-deployment.md` |
| AI 开发约束 | `docs/02-开发规范/02-ai-development-guide.md` |
