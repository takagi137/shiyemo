# proposal：合规修复 — 补全验证体系 & 消除架构违规

## 背景

2026-05-21 合规审计发现项目在验证体系和代码规范上存在多处违规：
- 验证三层中 L1 repo 层和 L2 集成测试完全缺失
- 1 处红线违规（handler 绕过 service）
- 前端大量重复代码违反 DRY 原则
- 项目管理文档状态不同步

## 目标用户

- 开发者（保证代码质量、可维护性）
- AI 开发助手（依赖文档正确性进行开发循环）

## 核心功能

1. 补全 repo_test.go（11 个模块）
2. 建立 L2 集成测试框架 + 核心模块集成测试
3. 修复 auth handler 架构违规
4. 前端公共组件/常量抽取（消除重复代码）
5. 同步项目管理文档
6. 拆分超大测试文件
7. 清理死代码和未使用依赖

## 不做范围

- 不重构 health/websocket 分层（P2，特殊模块可接受）
- 不引入新的状态管理库
- 不修改现有功能逻辑
- 不改动已通过的测试

## 成功标准

- [ ] 11 个模块 repo_test.go 全部存在，覆盖率 ≥ 70%
- [ ] tests/integration/ 目录存在，auth + post 模块集成测试通过
- [ ] auth handler GetProfile 不再直接访问 repo
- [ ] POST_TYPES / timeAgo / Pagination 各只定义一处
- [ ] TODO.md / PROJECT_STATUS.md / NEXT_PLAN.md 状态一致
- [ ] 无 >300 行的测试文件
- [ ] package.json 无未使用依赖
