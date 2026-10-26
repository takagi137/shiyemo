# design：他人主页

## 技术决策
| 决策点 | 选择 | 理由 |
|--------|------|------|
| 路由 | /users/[id] | Next.js 动态路由 |
| 数据获取 | 客户端请求 | 复用已有 API，无需 SSR |
| 组件复用 | 帖子卡片 | 复用 Feed 页的帖子展示组件 |

## 模块影响
| 模块 | 改动类型 | 说明 |
|------|---------|------|
| 前端 | 新增 | web/src/app/users/[id]/page.tsx |
| 前端 | 修改 | 帖子卡片的头像/昵称加链接 |

## 接口设计（已有，无需新增）
```
GET /api/v1/users/:id           → 用户信息
GET /api/v1/users/:id/stats     → 统计数据 {post_count, like_received, follower_count, following_count}
GET /api/v1/users/:id/posts?page=1&page_size=20 → 用户帖子列表
POST /api/v1/users/:id/follow   → 关注
DELETE /api/v1/users/:id/follow  → 取消关注
```

## 数据变更
无，全部使用已有 API

## 风险点
- 无
