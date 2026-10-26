# design：静思空间

## 技术决策
| 决策点 | 选择 | 理由 |
|--------|------|------|
| 呼吸动画 | CSS Animation + React 状态 | 纯前端实现，无需后端 |
| 环境音 | HTML5 Audio + 预置音频文件 | 简单可靠，无需流媒体服务 |
| 私密日记 | 后端存储 + 前端展示 | 数据持久化，支持多设备 |
| 今日阅读 | 后端 readings 表 + API | 支持管理后台编辑 |
| 页面框架 | Client Component 为主 | 大量交互（动画/音频/输入） |

## 模块影响
| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| web/app/sanctuary/page.tsx | 新建 | 静思空间主页面 |
| web/components/BreathingCircle.tsx | 新建 | 呼吸动画组件 |
| web/components/AmbientPlayer.tsx | 新建 | 环境音播放器 |
| web/components/JournalEntry.tsx | 新建 | 私密日记输入组件 |
| web/components/ReadingCard.tsx | 新建 | 今日阅读展示组件 |
| server/internal/journal/ | 新建 | 日记 CRUD 模块 |
| server/internal/reading/ | 新建 | 阅读 CRUD 模块 |
| server/cmd/main.go | 修改 | 注册 journal + reading 路由 |

## 接口设计
```
# 日记
POST /api/v1/journal
请求：{ "content": "今天..." }
响应：{ "id": 1, "content": "今天...", "created_at": "..." }

GET /api/v1/journal
响应：{ "items": [...], "total": 10 }

# 阅读
GET /api/v1/readings/today
响应：{
  "id": 1,
  "source": "苏轼《定风波》",
  "author": "苏轼",
  "content": "莫听穿林打叶声...",
  "interpretation": "苏轼被贬黄州时...",
  "image_url": "..."
}

GET /api/v1/readings/random
响应：同上（随机一篇）

POST /api/v1/readings/:id/favorite
响应：{ "message": "ok" }
```

## 数据变更
| 操作 | 对象 | 说明 |
|------|------|------|
| 新建 | journal 表 | id, user_id, content, created_at |
| 新建 | readings 表 | id, source, dynasty, author, title, content, interpretation, tags, image_url, is_active, created_at |

## 依赖
- 环境音文件：预置 6 个音频文件（CC 协议）
- Redis：缓存今日阅读

## 风险点
- 音频文件体积 → 使用压缩格式（MP3/OGG），单文件 < 1MB
- 日记安全 → 严格认证，只能访问自己的日记
