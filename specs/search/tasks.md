# tasks：搜索功能

## Phase 1：数据层

- [x] 任务 1.1：初始化 Meilisearch 索引 ✅ (2026-05-29 同步)
  - 验证：索引创建成功
  - Spec：无

## Phase 2：业务层

- [x] 任务 2.1：实现 search service ✅ (2026-05-29 同步)
  - 文件：`server/internal/search/service.go`
  - 验证：SearchPosts, SearchUsers 方法
  - Spec：
    - 01-search-posts.feature.md
    - 02-search-users.feature.md

## Phase 3：接口层

- [x] 任务 3.1：实现 search handler ✅ (2026-05-29 同步)
  - 文件：`server/internal/search/handler.go`
  - 验证：
    - GET /api/v1/search/posts 可调用
    - GET /api/v1/search/users 可调用
  - Spec：全部 spec 文件

- [x] 任务 3.2：更新 main.go 路由注册 ✅ (2026-05-29 同步)
  - 文件：`server/cmd/main.go`
  - 验证：所有接口可访问

## Phase 4：验证

- [x] 任务 4.1：运行全部 BDD 场景 ✅ (2026-05-29 同步)
  - 验证：通过率 100%

- [x] 任务 任务 4.2：人工验收测试 ✅ (2026-05-29 同步)

## 依赖关系

```
Phase 1 → Phase 2 → Phase 3 → Phase 4
```

## 完成时间

- 2026-05-29 同步（代码已实现，spec 追补）
