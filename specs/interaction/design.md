# design：点赞/收藏功能

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 存储方式 | 独立表 | likes 和 bookmarks 分开，便于查询 |
| 唯一约束 | (user_id, post_id) | 防止重复点赞/收藏 |
| 计数更新 | 同步更新 posts 表 | 便于列表展示 |
| 取消方式 | 物理删除 | 简化逻辑 |

## 接口设计

### 1. 点赞/取消点赞

```
POST /api/v1/posts/:id/like
Authorization: Bearer <token>

响应（点赞）：
{
  "code": 0,
  "message": "ok",
  "data": { "liked": true }
}

响应（取消点赞）：
{
  "code": 0,
  "message": "ok",
  "data": { "liked": false }
}
```

### 2. 收藏/取消收藏

```
POST /api/v1/posts/:id/bookmark
Authorization: Bearer <token>

响应（收藏）：
{
  "code": 0,
  "message": "ok",
  "data": { "bookmarked": true }
}

响应（取消收藏）：
{
  "code": 0,
  "message": "ok",
  "data": { "bookmarked": false }
}
```

### 3. 查询用户收藏列表

```
GET /api/v1/users/:id/bookmarks?page=1&size=20
Authorization: Bearer <token>

响应：
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [...],
    "page": 1,
    "size": 20,
    "has_more": false
  }
}
```

## 数据变更

| 操作 | 对象 | 说明 |
|------|------|------|
| 新增 | likes 表 | 点赞记录 |
| 新增 | bookmarks 表 | 收藏记录 |
| 修改 | posts 表 | like_count 字段 |

## 依赖

- 已有 PostgreSQL：数据存储
- 已有 Post 模型：关联帖子
