# tasks：合规修复（修订版 — 2026-05-21 验证后）

## 审计修正说明

- repo_test.go：11/11 全部存在 ✅
- tests/integration/：已存在 4 个文件 ✅
- auth handler：已通过 service 调用 ✅
- Skeleton/LazyImage/swl/jose/react-hot-toast：均已清理 ✅

## Phase 1：前端 DRY 重构（3 任务）

- [x] Task 1.1：创建 lib/constants.ts + lib/utils.ts
- [x] Task 1.2：迁移 feed + search + profile 页面引用
- [x] Task 1.3：迁移 posts/[id] + notifications + admin 页面引用

## Phase 2：前端超大文件拆分（1 任务）

- [x] Task 2.1：拆分 profile/page.tsx（497 行）

## Phase 3：后端超大测试文件拆分（2 任务）

- [x] Task 3.1：拆分 post/service_test.go（515 行）
- [x] Task 3.2：拆分 user/service_test.go（438 行）

## Phase 4：验证（1 任务）

- [x] Task 4.1：全量验证

> 状态同步：2026-05-29，代码已全部实现
