# design：用户资料功能

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 头像存储 | 本地文件系统 | MVP 阶段简单，后续迁移到 OSS |
| 头像访问 | 静态文件服务 | Nginx/Caddy 直接服务 |
| 资料更新 | PATCH 部分更新 | 只更新传入的字段 |
| 公开资料 | 隐藏敏感信息 | 不暴露手机号、邮箱 |

## 接口设计

### 1. 获取当前用户资料

```
GET /api/v1/auth/profile
Authorization: Bearer <token>

响应：
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "phone": "13800138000",
    "nickname": "测试用户",
    "avatar": "/uploads/avatars/1.jpg",
    "bio": "这个人很懒，什么都没写",
    "company": "前公司",
    "occupation": "职业",
    "status": 1,
    "created_at": "2026-05-19T16:13:07Z",
    "updated_at": "2026-05-19T16:13:07Z"
  }
}
```

### 2. 获取指定用户公开资料

```
GET /api/v1/users/:id

响应：
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "nickname": "测试用户",
    "avatar": "/uploads/avatars/1.jpg",
    "bio": "这个人很懒，什么都没写",
    "company": "前公司",
    "occupation": "职业",
    "created_at": "2026-05-19T16:13:07Z"
  }
}
```

### 3. 更新用户资料

```
PUT /api/v1/users/:id
Authorization: Bearer <token>
Content-Type: application/json

请求：
{
  "nickname": "新昵称",
  "bio": "新的个人简介",
  "company": "新公司",
  "occupation": "新职业"
}

响应：
{
  "code": 0,
  "message": "ok",
  "data": null
}
```

### 4. 上传头像

```
POST /api/v1/users/avatar
Authorization: Bearer <token>
Content-Type: multipart/form-data

请求：
- avatar: 文件

响应：
{
  "code": 0,
  "message": "ok",
  "data": {
    "avatar_url": "/uploads/avatars/1.jpg"
  }
}
```

## 数据变更

| 操作 | 对象 | 说明 |
|------|------|------|
| 修改 | users 表 | 添加 bio 字段 |
| 新增 | uploads 目录 | 存储头像文件 |

## 依赖

- 已有 PostgreSQL：用户数据存储
- 本地文件系统：头像存储（MVP 阶段）

## 风险点

1. 头像文件大小限制（建议 2MB）
2. 头像格式限制（jpg/png/gif）
3. 文件名冲突（使用用户ID作为文件名）
