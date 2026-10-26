# design：关注动态 Feed

## 技术决策
| 决策点 | 选择 | 理由 |
|--------|------|------|
| 接口设计 | GET /posts/following | RESTful 规范 |
| 数据查询 | 关联 follows 表 | 获取关注用户 ID 列表 |

## 模块影响
| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| post/handler | 新增 | FollowingFeed 接口 |
| post/service | 新增 | FollowingFeed 方法 |
| post/repo | 新增 | GetByUserIDs 方法 |
| web/feed | 修改 | 添加"关注"Tab |

## 接口设计
```
GET /api/v1/posts/following?page=1&page_size=10
响应：{ items: Post[], total: number }
```
