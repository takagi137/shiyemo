# tasks：首页温馨升级

## Phase 1：数据层

- [x] 任务 1.1：创建 tips 表
  - 验收：`go run cmd/migrate/main.go up` 成功
  - Spec：02-daily-api.md

- [x] 任务 1.2：实现 daily repo（CRUD + 随机选取）
  - 验收：单元测试通过
  - Spec：02-daily-api.md

## Phase 2：业务层

- [x] 任务 2.1：实现 daily service（获取今日内容 + 缓存逻辑）
  - 验收：单元测试通过
  - Spec：02-daily-api.md

- [x] 任务 2.2：集成 Unsplash API（获取风景图）
  - 验收：能返回有效的图片 URL
  - Spec：01-homepage-scene.md

## Phase 3：接口层

- [x] 任务 3.1：实现 GET /api/v1/daily handler
  - 验收：curl 请求返回正确数据
  - Spec：02-daily-api.md

## Phase 4：前端

- [x] 任务 4.1：创建 DailyScene 组件
  - 验收：组件能展示图片 + 文案
  - Spec：01-homepage-scene.md

- [x] 任务 4.2：重写首页 page.tsx
  - 验收：页面加载展示风景图 + 文案 + 入口按钮
  - Spec：01-homepage-scene.md

## Phase 5：验证

- [x] 任务 5.1：部署到验证环境
  - 验收：首页正常展示

- [x] 任务 5.2：BDD 场景验证
  - 验收：所有 spec 通过
