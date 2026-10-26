# 单元测试覆盖率标准

> 持久化规范，所有覆盖率评估按此标准执行。

## 核心原则

1. **只测核心业务逻辑**：service 层（含 skill_service 等业务变体）是测试重点
2. **不测的文件类型**：
   - `handler.go` — HTTP 薄封装，参数校验 + 调 service，适合 E2E/集成测试
   - `repo.go` — 纯 DB 操作（GORM 调用），适合集成测试
   - `model.go` — 结构体定义、TableName()，无业务逻辑
   - `*_repo.go` — 仓储接口定义
   - 配置文件、路由注册等
3. **不修改断言来通过测试**：测试应验证设计意图，如果代码返回值不符合预期，修改代码而非断言
4. **不做假测试**：不允许 mock-skip、t.Skip 等方式跳过验证
5. **需要 DB 的模块 mock DB 返回**：通过接口注入 mock 仓储，不依赖真实数据库

## 覆盖率计算方法

```bash
# 1. 运行测试生成覆盖率文件
go test ./internal/<module>/... -count=1 -coverprofile=coverage.out

# 2. 只统计 service.go 相关的语句覆盖率
awk '/service\.go/{split($2, a, " "); stmts=a[1]; count=$3; total+=stmts; if(count>0) covered+=stmts}
END{if(total>0) printf "%d/%d = %.1f%%\n", covered, total, covered*100/total}' coverage.out
```

## 覆盖率目标

| 层级 | 目标 | 说明 |
|------|------|------|
| service 层（核心业务） | ≥ 80% | 每个公开方法的主路径 + 异常路径 |
| handler 层 | 不要求 | 薄封装，测试价值低 |
| repo 层 | 不要求 | 纯 DB 操作，需集成测试 |
| model 层 | 不要求 | 结构体定义 |

## 当前覆盖率（2026-05-21）

| 模块 | service 覆盖率 | 状态 | 备注 |
|------|---------------|------|------|
| moderation | 97.7% | ✅ | admin_service + core service |
| tag | 93.3% | ✅ | core service + skill_service |
| upload | 88.5% | ✅ | Upload 94%, getContentType 83% |
| export | 87.5% | ✅ | ExportUsers 88%, ExportPosts 86% |
| comment | 84.6% | ✅ | Create 79%, Delete 92% |
| report | 84.6% | ✅ | Create 89%, List 100%, Handle 100% |
| auth | 83.9% | ✅ | Register 83%, Login 85% |
| interaction | 83.3% | ✅ | Like 90%, Bookmark 80%, ListBookmarks 80% |
| follow | 81.8% | ✅ | Follow 84%, ListFollowing 80%, ListFollowers 80% |
| search | 67.0% | ⚠️ | InitIndexes 73%, GetSuggestions 94% |
| websocket | 70.7% | ⚠️ | hub.run 70.7% |
| admin | N/A | 🔒 | 需 DB |
| notification | 95%+ | ✅ | Create 100%, MarkAsRead 87.5%, 余 100% |
| post | N/A | 🔒 | 需 DB |
| user | N/A | 🔒 | 需 DB |

## 测试执行命令

```bash
# 运行所有非 DB 模块测试
go test -p 1 ./internal/auth/... ./internal/comment/... ./internal/common/cache/... \
  ./internal/common/storage/... ./internal/export/... ./internal/follow/... \
  ./internal/health/... ./internal/interaction/... ./internal/log/... \
  ./internal/moderation/... ./internal/report/... ./internal/search/... \
  ./internal/tag/... ./internal/upload/... ./internal/websocket/... -count=1
```
