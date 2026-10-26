# 失业么 API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: JWT Bearer Token
- **响应格式**: JSON

## 响应结构

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400xx | 参数错误 |
| 401xx | 认证错误 |
| 403xx | 权限错误 |
| 404xx | 资源不存在 |
| 500xx | 服务器错误 |

---

## 认证模块

### 邮箱注册

```
POST /auth/register/email
```

**请求体**:
```json
{
  "email": "user@example.com",
  "password": "a123456",
  "nickname": "用户昵称"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "jwt_token",
    "user": {
      "id": 1,
      "nickname": "用户昵称",
      "email": "user@example.com"
    }
  }
}
```

### 邮箱登录

```
POST /auth/login/email
```

**请求体**:
```json
{
  "email": "user@example.com",
  "password": "a123456"
}
```

---

## 帖子模块

### 获取帖子列表

```
GET /posts?type={type}&page={page}&page_size={size}
```

**参数**:
- `type`: 帖子类型 (emotion/insight/case/action)
- `page`: 页码
- `page_size`: 每页数量

### 获取帖子详情

```
GET /posts/{id}
```

### 创建帖子

```
POST /posts
```

**请求体**:
```json
{
  "type": "emotion",
  "title": "帖子标题",
  "content": "帖子内容",
  "tags": ["标签1", "标签2"]
}
```

### 更新帖子

```
PUT /posts/{id}
```

### 删除帖子

```
DELETE /posts/{id}
```

### 关注动态 Feed

```
GET /posts/following?page={page}&page_size={size}
```

---

## 评论模块

### 创建评论

```
POST /posts/{id}/comments
```

**请求体**:
```json
{
  "content": "评论内容"
}
```

### 删除评论

```
DELETE /comments/{id}
```

---

## 互动模块

### 点赞/取消点赞

```
POST /posts/{id}/like
```

### 收藏/取消收藏

```
POST /posts/{id}/bookmark
```

### 获取用户收藏列表

```
GET /users/{id}/bookmarks?page={page}&page_size={size}
```

---

## 关注模块

### 关注/取消关注用户

```
POST /users/{id}/follow
```

---

## 举报模块

### 创建举报

```
POST /reports
```

**请求体**:
```json
{
  "target_type": "post",
  "target_id": 123,
  "reason": "违规内容描述"
}
```

### 获取举报列表（管理员）

```
GET /admin/reports?page={page}&page_size={size}&status={status}
```

### 处理举报（管理员）

```
PUT /admin/reports/{id}
```

**请求体**:
```json
{
  "status": 2,
  "handle_reason": "处理原因"
}
```

---

## 管理员模块

### 获取帖子列表（管理员）

```
GET /admin/posts?page={page}&page_size={size}&status={status}
```

### 删除帖子（管理员）

```
DELETE /admin/posts/{id}
```

### 获取用户列表（管理员）

```
GET /admin/users?page={page}&page_size={size}
```

### 更新用户状态（管理员）

```
PUT /admin/users/{id}/status
```

**请求体**:
```json
{
  "status": 1
}
```

### 获取统计数据（管理员）

```
GET /admin/stats
```

### 设置帖子置顶（管理员）

```
PUT /admin/posts/{id}/top
```

**请求体**:
```json
{
  "is_top": true
}
```

### 设置帖子加精（管理员）

```
PUT /admin/posts/{id}/featured
```

**请求体**:
```json
{
  "is_featured": true
}
```

---

## 通知模块

### 获取通知列表

```
GET /notifications?page={page}&page_size={size}
```

### 获取未读通知数量

```
GET /notifications/unread-count
```

### 标记通知已读

```
PUT /notifications/{id}/read
```

### 标记所有通知已读

```
PUT /notifications/read-all
```

---

## 搜索模块

### 搜索帖子

```
GET /search?q={query}&type={type}&page={page}&page_size={size}
```

---

## 文件上传模块

### 上传单张图片

```
POST /upload/image
```

**请求体**: `multipart/form-data` 格式，字段名为 `file`

### 上传多张图片

```
POST /upload/images
```

**请求体**: `multipart/form-data` 格式，字段名为 `files`，最多 9 张

---

## 内容审核模块

### 获取敏感词列表（管理员）

```
GET /admin/moderation/words?page={page}&page_size={size}
```

### 添加敏感词（管理员）

```
POST /admin/moderation/words
```

**请求体**:
```json
{
  "word": "敏感词"
}
```

### 删除敏感词（管理员）

```
DELETE /admin/moderation/words/{id}
```

---

## 数据导出模块

### 导出用户数据（管理员）

```
GET /admin/export/users?format={format}
```

### 导出帖子数据（管理员）

```
GET /admin/export/posts?format={format}&type={type}
```

### 获取导出统计（管理员）

```
GET /admin/export/stats
```

---

## WebSocket 实时通知

### 连接

```
ws://localhost:8080/ws?token={jwt_token}
```

### 消息格式

```json
{
  "type": "notification",
  "data": {
    "id": 1,
    "type": "like",
    "content": "有人点赞了你的帖子",
    "created_at": "2026-05-20T10:00:00Z"
  }
}
```

---

## 用户模块

### 获取当前用户信息

```
GET /user/me
```

### 更新用户信息

```
PUT /user/me
```

**请求体**:
```json
{
  "nickname": "新昵称",
  "bio": "个人简介"
}
```

### 获取用户信息

```
GET /users/{id}
```
