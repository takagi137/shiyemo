# 设计：管理后台

## 权限设计

MVP 阶段使用简单的 is_admin 字段判断管理员权限。

### users 表扩展

| 字段 | 类型 | 说明 |
|------|------|------|
| is_admin | BOOLEAN | 是否管理员，默认 false |

## API 设计

### 获取帖子列表（管理员）

```
GET /api/v1/admin/posts?page=1&page_size=20&status=all
Authorization: Bearer <admin_token>
```

参数：
- status: all/active/deleted（默认 active）

### 删除帖子（管理员）

```
DELETE /api/v1/admin/posts/:id
Authorization: Bearer <admin_token>
```

### 获取用户列表

```
GET /api/v1/admin/users?page=1&page_size=20
Authorization: Bearer <admin_token>
```

### 禁用/启用用户

```
PUT /api/v1/admin/users/:id/status
Authorization: Bearer <admin_token>
Body: {"status": 0}  // 0-禁用 1-启用
```

### 平台统计数据

```
GET /api/v1/admin/stats
Authorization: Bearer <admin_token>
```

响应：
```json
{
  "code": 0,
  "data": {
    "total_users": 100,
    "total_posts": 500,
    "total_comments": 2000,
    "active_users_today": 50
  }
}
```
