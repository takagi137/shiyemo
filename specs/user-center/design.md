# 设计：个人中心

## API 设计

### 获取我的帖子列表

```
GET /api/v1/users/:id/posts?page=1&page_size=20
Authorization: Bearer <token>
```

响应：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "total": 50,
    "page": 1,
    "page_size": 20,
    "items": [...]
  }
}
```

### 获取个人统计数据

```
GET /api/v1/users/:id/stats
Authorization: Bearer <token>
```

响应：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "post_count": 50,
    "like_received": 120,
    "follower_count": 30,
    "following_count": 25
  }
}
```

## 依赖关系

- 复用已有的收藏列表接口：`GET /api/v1/users/:id/bookmarks`
- 复用已有的关注/粉丝列表接口：`GET /api/v1/users/:id/following` / `followers`
