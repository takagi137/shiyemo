# feature: 评论回复通知

## 场景 1：评论回复时发送通知

### Given 用户 A 发表了评论
### When 用户 B 回复该评论
### Then 用户 A 收到通知，类型为 comment

```json
// 前置条件：用户 A 创建评论
POST /api/v1/posts/1/comments
Authorization: Bearer <token_a>
{
  "content": "用户A的评论"
}
// 返回评论 ID = 1

// 执行：用户 B 回复评论
POST /api/v1/posts/1/comments
Authorization: Bearer <token_b>
{
  "content": "用户B的回复",
  "parent_id": 1
}

// 验证：用户 A 收到通知
GET /api/v1/notifications?type=comment
Authorization: Bearer <token_a>

// 预期响应
{
  "code": 0,
  "data": {
    "items": [
      {
        "type": "comment",
        "content": "用户B 回复了你的评论"
      }
    ]
  }
}
```
