# feature: 内容审核队列

## 场景 1：查询待审核帖子

### Given 管理员已登录，存在待审核帖子
### When 管理员请求 GET /api/v1/admin/posts/pending
### Then 返回待审核帖子列表

```json
// 请求
GET /api/v1/admin/posts/pending?page=1&page_size=10
Authorization: Bearer <admin_token>

// 预期响应
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "待审核帖子",
        "review_status": "pending"
      }
    ],
    "total": 1
  }
}
```

## 场景 2：审核通过帖子

### Given 管理员已登录，存在待审核帖子 ID=1
### When 管理员请求 POST /api/v1/admin/posts/1/approve
### Then 帖子状态变为 approved

```json
// 请求
POST /api/v1/admin/posts/1/approve
Authorization: Bearer <admin_token>

// 预期响应
{
  "code": 0,
  "message": "ok"
}

// 验证：查询帖子状态
GET /api/v1/posts/1
// 预期 review_status = "approved"
```

## 场景 3：审核拒绝帖子

### Given 管理员已登录，存在待审核帖子 ID=1
### When 管理员请求 POST /api/v1/admin/posts/1/reject
### Then 帖子状态变为 rejected，reject_reason 为 "违规"

```json
// 请求
POST /api/v1/admin/posts/1/reject
Authorization: Bearer <admin_token>
{
  "reason": "违规"
}

// 预期响应
{
  "code": 0,
  "message": "ok"
}

// 验证：查询帖子状态
GET /api/v1/posts/1
// 预期 review_status = "rejected", reject_reason = "违规"
```
