# design：帖子编辑功能

## 技术决策
| 决策点 | 选择 | 理由 |
|--------|------|------|
| 编辑方式 | 独立编辑页面 | 复用发帖页面结构，用户体验一致 |
| 接口复用 | PUT /posts/:id | RESTful 规范 |

## 模块影响
| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| post/handler | 确认 | 确认 UpdatePost 接口已存在 |
| web/posts/edit | 新增 | 编辑页面 |
| web/posts/[id] | 修改 | 添加编辑按钮 |

## 接口设计
```
PUT /api/v1/posts/:id
请求：{ title, content, type, tags }
响应：{ post }
```

## 数据变更
无（复用现有 posts 表）
