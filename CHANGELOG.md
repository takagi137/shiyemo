# Changelog

## 2026-05-19 - 验证标准建立

### 新增
- 验证三层架构文档 (docs/04-开发流程/02-verification-guide.md)
- BDD 场景自动化验证脚本 (scripts/verify-bdd.sh)
- 测试工具函数库 (server/testutil/testutil.go)
- Auth 模块单元测试 (server/internal/auth/service_test.go)
- Post 模块单元测试 (server/internal/post/service_test.go)
- 验证标准 OpenSpec 文档 (specs/verification-standard/)

### 更新
- AGENTS.md - 新增验证优先原则
- docs/02-开发规范/02-ai-development-guide.md - 新增验证层次模型

### 验证三层架构

| 层级 | 名称 | 覆盖范围 | 工具 |
|------|------|---------|------|
| L1 | 单元测试 | 每个函数/方法 | `go test` + `testify` |
| L2 | 集成测试 | API 端点 | `httptest` + 真实 DB |
| L3 | 部署验证 | 完整流程 | curl + BDD 场景 |

### 验证检查清单

每个功能完成前必须确认：
- [ ] L1：单元测试通过（`go test ./...`）
- [ ] L1：测试覆盖率达标（service ≥ 80%）
- [ ] L2：集成测试通过
- [ ] L3：部署到验证环境成功
- [ ] L3：BDD 场景全部通过
- [ ] L3：关键路径手动验证通过

---

## 2026-05-19 - 搜索功能完成

### 新增
- 搜索服务 (server/internal/search/service.go)
- 搜索处理器 (server/internal/search/handler.go)
- 搜索 OpenSpec 文档 (specs/search/)

### 验证结果
- 搜索帖子：✅ 通过
- 搜索用户：✅ 通过

---

## 2026-05-19 - 关注功能完成

### 新增
- 关注服务 (server/internal/follow/service.go)
- 关注处理器 (server/internal/follow/handler.go)
- 关注 OpenSpec 文档 (specs/follow/)

### 验证结果
- 关注用户：✅ 通过
- 取消关注：✅ 通过
- 关注列表：✅ 通过
- 粉丝列表：✅ 通过

---

## 2026-05-19 - 点赞/收藏功能完成

### 新增
- 点赞/收藏服务 (server/internal/interaction/service.go)
- 点赞/收藏处理器 (server/internal/interaction/handler.go)
- 点赞/收藏 OpenSpec 文档 (specs/interaction/)

### 验证结果
- 点赞帖子：✅ 通过
- 取消点赞：✅ 通过
- 收藏帖子：✅ 通过
- 取消收藏：✅ 通过
- 收藏列表：✅ 通过

---

## 2026-05-19 - 评论功能完成

### 新增
- 评论服务 (server/internal/comment/service.go)
- 评论处理器 (server/internal/comment/handler.go)
- 评论 OpenSpec 文档 (specs/comment/)

### 验证结果
- 发表评论：✅ 通过
- 回复评论：✅ 通过
- 评论列表：✅ 通过
- 删除评论：✅ 通过

---

## 2026-05-19 - 发帖功能完成

### 新增
- 帖子服务 (server/internal/post/service.go)
- 帖子处理器 (server/internal/post/handler.go)
- 帖子 OpenSpec 文档 (specs/post-crud/)

### 验证结果
- 创建帖子：✅ 通过
- 帖子详情：✅ 通过
- 帖子列表：✅ 通过
- 更新帖子：✅ 通过
- 删除帖子：✅ 通过

---

## 2026-05-19 - 用户资料功能完成

### 新增
- 用户资料服务 (server/internal/user/service.go)
- 用户资料处理器 (server/internal/user/handler.go)
- 用户资料 OpenSpec 文档 (specs/user-profile/)

### 验证结果
- 获取用户信息：✅ 通过
- 更新用户信息：✅ 通过
- 上传头像：✅ 通过

---

## 2026-05-19 - 用户注册/登录功能完成

### 新增
- 认证服务 (server/internal/auth/service.go)
- 认证处理器 (server/internal/auth/handler.go)
- 认证 OpenSpec 文档 (specs/user-auth/)

### 验证结果
- 邮箱注册：✅ 通过
- 手机注册：✅ 通过
- 邮箱登录：✅ 通过
- 手机登录：✅ 通过
- 发送验证码：✅ 通过
- 获取用户信息：✅ 通过

---

## 2026-05-19 - 项目初始化

### 新增
- 项目脚手架
- Docker Compose 配置
- 数据库迁移
- 基础架构文档
