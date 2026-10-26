# design：首页温馨升级

## 技术决策
| 决策点 | 选择 | 理由 |
|--------|------|------|
| 风景图来源 | Unsplash API | 免费高清摄影图库，无需自建图库 |
| 文案存储 | 数据库 tips 表 | 支持管理后台编辑，不硬编码 |
| 页面框架 | Server Component + Client 交互 | 风景图/文案 SSR，交互按钮 CSR |
| 图片加载 | 渐进式 + 占位色 | 防止白屏闪烁 |

## 模块影响
| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| web/app/page.tsx | 重写 | 新首页（风景图 + 文案 + 入口） |
| web/components/DailyScene.tsx | 新建 | 风景图 + 文案展示组件 |
| server/internal/daily/ | 新建 | 每日内容 API（图片 + 文案） |
| server/cmd/main.go | 修改 | 注册 daily 路由 |

## 接口设计
```
GET /api/v1/daily
响应：
{
  "code": 0,
  "data": {
    "image": {
      "url": "https://images.unsplash.com/...",
      "description": "日落海滩",
      "credit": "Photo by xxx on Unsplash"
    },
    "quote": {
      "text": "今天的你，已经比昨天勇敢了一点。",
      "author": ""
    }
  }
}
```

## 数据变更
| 操作 | 对象 | 说明 |
|------|------|------|
| 新建 | tips 表 | id, content, author, category, is_active, created_at |

## 依赖
- Unsplash API（免费，50次/小时）
- Redis（缓存今日内容，避免重复请求）

## 风险点
- Unsplash API 限流 → 缓存策略：每日只请求一次，存 Redis
- 图片加载慢 → 使用 Next.js Image 组件 + 占位背景色
