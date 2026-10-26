# design：我的收藏页

## 技术决策
| 决策点 | 选择 | 理由 |
|--------|------|------|
| 路由 | /bookmarks | 独立页面 |
| 数据获取 | 客户端请求 | 复用已有 API |

## 模块影响
| 模块 | 改动类型 | 说明 |
|------|---------|------|
| 前端 | 新增 | web/src/app/bookmarks/page.tsx |
| 前端 | 修改 | profile/page.tsx 加"我的收藏"入口 |

## 接口设计（已有）
```
GET /api/v1/users/:id/bookmarks?page=1&size=20 → 收藏列表
```

## 数据变更
无
