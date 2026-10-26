# spec：修复 auth handler 架构违规

## 场景描述

消除 auth handler 中直接访问 repo 的红线违规。

### Given 前置条件
- `auth/handler.go` 的 `GetProfile` 方法直接调用 `h.service.repo.GetByID()`
- 这违反红线 #6「不在 handler 层写业务逻辑」

### When 操作
- 在 `auth/service.go` 中新增 `GetProfile(userID int64) (*User, error)` 方法
- 修改 `auth/handler.go` 的 `GetProfile` 改为调用 `h.service.GetProfile(userID)`

### Then 预期结果
- `grep "h.service.repo" auth/handler.go` 返回 0 结果
- `go build ./...` 编译通过
- `go test ./internal/auth/...` 全部通过
- BDD 场景：登录后 GET /api/v1/user/profile 返回 200 + 正确用户信息
