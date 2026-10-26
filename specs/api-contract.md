# API 契约

> 本文件定义失业么项目的 REST API 接口规范。
> **后端 agent** 新增/修改 API 时必须同步更新本文件。
> **前端 agent** 对接 API 时从本文件读取接口定义。

## 通用约定

- **Base URL**: `/api/v1`
- **认证**: JWT Bearer Token（`Authorization: Bearer <token>`）
- **Content-Type**: `application/json`
- **时间格式**: RFC3339（`2026-05-27T12:00:00Z`）

## 认证模块 (auth)

### POST /api/v1/auth/register
- 请求: `{ "email": "string", "password": "string", "nickname": "string" }`
- 响应: `{ "id": 1, "email": "...", "nickname": "...", "token": "..." }`
- 状态码: 201 / 400 / 409（邮箱已存在）

### POST /api/v1/auth/login
- 请求: `{ "email": "string", "password": "string" }`
- 响应: `{ "token": "...", "user": { "id": 1, "email": "...", "nickname": "..." } }`
- 状态码: 200 / 401

### GET /api/v1/auth/profile
- 认证: 需要
- 响应: `{ "id": 1, "email": "...", "nickname": "...", ... }`
- 状态码: 200 / 401

## 帖子模块 (post)

### GET /api/v1/posts
- 参数: `page` (int), `type` (string, 可选)
- 响应: `{ "posts": [...], "total": 100, "page": 1, "has_more": true }`

### GET /api/v1/posts/:id
- 响应: `{ "id": 1, "title": "...", "content": "...", "type": "emotion", ... }`

### POST /api/v1/posts
- 认证: 需要
- 请求: `{ "title": "string", "content": "string", "type": "emotion|insight|case|action" }`
- 响应: 201 + 帖子对象

### PUT /api/v1/posts/:id
- 认证: 需要（作者）
- 请求: `{ "title": "...", "content": "..." }`

### DELETE /api/v1/posts/:id
- 认证: 需要（作者）

## 评论模块 (comment)

### GET /api/v1/posts/:id/comments
- 响应: `{ "comments": [...], "total": 10 }`

### POST /api/v1/posts/:id/comments
- 认证: 需要
- 请求: `{ "content": "string", "parent_id": null }`

## 互动模块 (interaction)

### POST /api/v1/posts/:id/like
- 认证: 需要

### DELETE /api/v1/posts/:id/like
- 认证: 需要

### POST /api/v1/posts/:id/bookmark
- 认证: 需要

### DELETE /api/v1/posts/:id/bookmark
- 认证: 需要

## 关注模块 (follow)

### POST /api/v1/users/:id/follow
- 认证: 需要

### DELETE /api/v1/users/:id/follow
- 认证: 需要

### GET /api/v1/users/:id/following
- 响应: `{ "users": [...], "total": 50 }`

### GET /api/v1/users/:id/followers
- 响应: `{ "users": [...], "total": 30 }`

## 通知模块 (notification)

### GET /api/v1/notifications
- 认证: 需要
- 参数: `page`, `type` (可选)
- 响应: `{ "notifications": [...], "unread_count": 5 }`

### PUT /api/v1/notifications/:id/read
- 认证: 需要

### PUT /api/v1/notifications/read-all
- 认证: 需要

## 搜索模块 (search)

### GET /api/v1/search
- 参数: `q` (string), `type` (string, 可选), `page`
- 响应: `{ "results": [...], "total": 20 }`

## 管理后台模块 (admin)

### GET /api/v1/admin/stats
- 认证: 需要（管理员）
- 响应: `{ "total_users": 100, "total_posts": 500, "total_comments": 2000, "active_users_today": 50 }`

### GET /api/v1/admin/posts
- 认证: 需要（管理员）
- 参数: `page`, `page_size`, `status` (active/deleted)
- 响应: `{ "items": [...], "total": 100, "page": 1, "page_size": 20 }`

### DELETE /api/v1/admin/posts/:id
- 认证: 需要（管理员）
- 说明: 删除帖子

### PUT /api/v1/admin/posts/:id/top
- 认证: 需要（管理员）
- 请求: `{ "is_top": true/false }`
- 说明: 设置/取消置顶

### PUT /api/v1/admin/posts/:id/featured
- 认证: 需要（管理员）
- 请求: `{ "is_featured": true/false }`
- 说明: 设置/取消加精

### GET /api/v1/admin/posts/pending
- 认证: 需要（管理员）
- 参数: `page`, `page_size`
- 响应: `{ "items": [...], "total": 10, "page": 1, "page_size": 20 }`
- 说明: 获取待审核帖子列表

### POST /api/v1/admin/posts/:id/approve
- 认证: 需要（管理员）
- 说明: 审核通过帖子

### POST /api/v1/admin/posts/:id/reject
- 认证: 需要（管理员）
- 请求: `{ "reason": "拒绝原因" }`
- 说明: 审核拒绝帖子

### GET /api/v1/admin/users
- 认证: 需要（管理员）
- 参数: `page`, `page_size`, `search`, `status`, `role`
- 响应: `{ "items": [...], "total": 100, "page": 1, "page_size": 20 }`
- 说明: search 支持按昵称/邮箱/手机号搜索

### PUT /api/v1/admin/users/:id/status
- 认证: 需要（管理员）
- 请求: `{ "status": 0/1 }` (0=禁用, 1=正常)

### GET /api/v1/admin/comments
- 认证: 需要（管理员）
- 参数: `page`, `page_size`, `post_id` (可选)
- 响应: `{ "items": [...], "total": 100, "page": 1, "page_size": 20 }`

### PUT /api/v1/admin/comments/:id/hide
- 认证: 需要（管理员）
- 说明: 屏蔽评论

### PUT /api/v1/admin/comments/:id/restore
- 认证: 需要（管理员）
- 说明: 恢复评论

### GET /api/v1/admin/reports
- 认证: 需要（管理员）
- 参数: `page`, `size`, `status` (可选, 1=待处理/2=已处理/3=已忽略)
- 响应: `{ "items": [...], "total": 10, "page": 1, "page_size": 20 }`

### PUT /api/v1/admin/reports/:id
- 认证: 需要（管理员）
- 请求: `{ "status": 2/3 }` (2=已处理, 3=已忽略)
- 说明: 处理举报

### GET /api/v1/admin/moderation/words
- 认证: 需要（管理员）
- 参数: `page`, `page_size`
- 响应: `{ "items": [...], "total": 50, "page": 1, "page_size": 20 }`

### POST /api/v1/admin/moderation/words
- 认证: 需要（管理员）
- 请求: `{ "word": "敏感词" }`

### DELETE /api/v1/admin/moderation/words/:id
- 认证: 需要（管理员）

## 用户模块 (user)

### GET /api/v1/users/:id
- 响应: 用户对象

### PUT /api/v1/users/:id
- 认证: 需要（本人）
- 请求: `{ "nickname": "...", "bio": "...", "company": "...", "occupation": "..." }`

### GET /api/v1/users/:id/posts
- 参数: `page`
- 响应: `{ "posts": [...], "total": 10 }`

## 标签模块 (tag)

### GET /api/v1/users/:id/skills
- 响应: `{ "skills": ["Go", "React", ...] }`

### POST /api/v1/users/:id/skills
- 认证: 需要（本人）
- 请求: `{ "skill": "Go" }`

### DELETE /api/v1/users/:id/skills/:skill
- 认证: 需要（本人）

## 上传模块 (upload)

### POST /api/v1/upload/image
- 认证: 需要
- Content-Type: multipart/form-data
- 响应: `{ "url": "/uploads/xxx.jpg" }`

## 审核模块 (moderation)

### GET /api/v1/admin/moderation/words
- 认证: 需要（管理员）
- 响应: `{ "words": [...], "total": 50 }`

### POST /api/v1/admin/moderation/words
- 认证: 需要（管理员）
- 请求: `{ "word": "敏感词", "level": "high|medium|low" }`

### DELETE /api/v1/admin/moderation/words/:id
- 认证: 需要（管理员）

## 举报模块 (report)

### POST /api/v1/reports
- 认证: 需要
- 请求: `{ "target_type": "post|comment|user", "target_id": 1, "reason": "spam" }`

### GET /api/v1/admin/reports
- 认证: 需要（管理员）

## WebSocket (websocket)

### GET /api/v1/ws
- 升级: WebSocket
- 认证: JWT Token（query param）

## 导出模块 (export)

### GET /api/v1/admin/export/users
- 认证: 需要（管理员）
- Content-Type: text/csv

### GET /api/v1/admin/export/posts
- 认证: 需要（管理员）
- Content-Type: text/csv

## 签到模块 (checkin)

### POST /api/v1/checkin
- 认证: 需要

### GET /api/v1/checkin/status
- 认证: 需要
- 响应: `{ "checked_in": true, "streak": 7 }`

## 同行模块 (cohort)

### GET /api/v1/cohorts
- 响应: `{ "cohorts": [...] }`

### POST /api/v1/cohorts/:id/join
- 认证: 需要

## 验证模块 (verify)

### POST /api/v1/verify/email/send
- 认证: 需要
- 请求: `{ "email": "..." }`

### POST /api/v1/verify/email/confirm
- 认证: 需要
- 请求: `{ "code": "123456" }`

## 健康检查 (health)

### GET /health
- 响应: `{ "status": "ok" }`

### GET /health/detailed
- 响应: `{ "status": "ok", "db": "ok", "redis": "ok", "uptime": "2h30m" }`

## 备份模块 (backup)

### POST /api/v1/backup/create
- 认证: 需要（管理员）

### GET /api/v1/backup/list
- 认证: 需要（管理员）

### POST /api/v1/backup/restore
- 认证: 需要（管理员）
- 请求: `{ "filename": "backup_20260527.sql" }`

## 数据统计模块 (stats)

### GET /api/v1/stats/overview
- 认证: 需要（管理员）
- 响应: `{ "total_users": 100, "total_posts": 500, "today_posts": 10 }`

## 小贴士模块 (tips)

### GET /api/v1/tips/random
- 响应: `{ "content": "...", "category": "job_search" }`
