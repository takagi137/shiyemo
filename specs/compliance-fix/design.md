# design：合规修复

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| repo 测试方式 | sqlmock + testify | 不依赖真实 DB，与现有 service_test mock 模式一致 |
| 集成测试框架 | httptest + 测试 DB | Go 标准库，无需额外依赖 |
| 前端常量位置 | `lib/constants.ts` | 集中管理，所有页面统一引用 |
| 前端工具函数位置 | `lib/utils.ts` | 与现有 `lib/api.ts` 同级 |
| 分页组件位置 | `app/components/Pagination.tsx` | 与 Nav/Avatar 同级 |
| 测试文件拆分策略 | 按功能域拆分为多文件 | 每文件 ≤ 300 行 |

## 模块影响

| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| server/internal/*/repo.go | 新增测试 | 添加 repo_test.go |
| server/tests/integration/ | 新建 | 集成测试目录 |
| server/internal/auth/ | 重构 | handler 移逻辑到 service |
| web/lib/ | 新增 | constants.ts, utils.ts |
| web/app/components/ | 新增 | Pagination.tsx |
| web/app/*/page.tsx | 修改 | 引用公共模块 |
| 根目录文档 | 修改 | 同步状态 |

## 风险点

- repo_test 引入 sqlmock 需确认与 GORM 兼容
- 前端组件抽取后需确认所有页面样式不变
- 文档同步需确认当前真实进度
