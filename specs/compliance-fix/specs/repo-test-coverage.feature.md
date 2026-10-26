# spec：repo 层测试覆盖率达标

## 场景描述

验证所有 11 个模块的 repo_test.go 存在且覆盖率达标。

### Given 前置条件
- 11 个模块各有 repo.go 文件
- 项目可正常编译（`go build ./...` 通过）

### When 操作
- 为每个模块创建 repo_test.go
- 使用 sqlmock 模拟数据库交互
- 每个 repo 方法至少 3 个测试用例（正常/边界/错误）

### Then 预期结果
- 每个模块目录下存在 repo_test.go
- `go test ./internal/<module>/...` 全部通过
- `go test -cover ./internal/<module>/...` 中 repo.go 覆盖率 ≥ 70%
- 不依赖真实数据库连接
