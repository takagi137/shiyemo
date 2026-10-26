# design：关注功能

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 存储方式 | follows 表 | 记录关注关系 |
| 唯一约束 | (follower_id, following_id) | 防止重复关注 |
| 自关注 | 禁止 | 业务规则 |
| 取消方式 | 物理删除 | 简化逻辑 |

## 接口设计

### 1. 关注/取消关注

```
POST /api/v1/users/:id/follow
Authorization: Bearer <token>

响应（关注）：
{
  "code": 0,
  "message": "ok",
  "data": { "following": true }
}

响应（取消关注）：
{
  "code": 0,
  "message": "ok",
  "data": { "following": false }
}
```

### 2. 获取关注列表

```
GET /api/v1/users/:id/following?page=1&size=20

响应：
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "id": 2,
        "nickname": "用户2",
        "avatar": "/uploads/avatars/2.jpg",
        "bio": "个人简介"
      }
    ],
    "page": 1,
    "size": 20,
    "has_more": false
  }
}
```

### 3. 获取粉丝列表

```
GET /api/v1/users/:id/followers?page=1&size=20

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
| 新增 | follows 表 | 关注关系表 |
| 修改 | users 表 | following_count, follower_count |

## 依赖

- 已有 PostgreSQL：数据存储
- 已有 User 模型：关联用户
