# design：二次开发计划 - 技术设计

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 富文本存储 | PostgreSQL JSONB | 复用现有技术栈，支持结构化查询 |
| 图片存储 | MinIO 对象存储 | 已配置，支持 CDN 加速 |
| 通知类型 | 枚举字段 | 简单高效，支持扩展 |
| 审核状态 | 状态机 | pending → approved/rejected，清晰明确 |

## 模块影响分析

### Phase 1: 用户主页数据

| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| user/service | 修改 | 添加 GetUserStats 方法 |
| user/handler | 修改 | 扩展 Profile 接口返回统计数据 |

### Phase 2: 通知分类 + 批量操作

| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| notification/model | 修改 | 添加 type 字段枚举 |
| notification/repo | 修改 | 添加批量操作方法 |
| notification/service | 修改 | 添加分类查询、批量操作 |
| notification/handler | 修改 | 添加批量 API |

### Phase 3: 评论回复通知

| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| comment/service | 修改 | 发送回复通知 |

### Phase 4: 搜索增强

| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| search/service | 修改 | 添加用户搜索、标签搜索 |
| search/handler | 修改 | 添加搜索建议接口 |

### Phase 5: 帖子编辑增强

| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| post/model | 修改 | content 字段改为 JSONB |
| post/service | 修改 | 支持富文本存储 |
| upload | 新增 | 图片上传接口 |

### Phase 6: 内容审核队列

| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| admin/service | 修改 | 添加审核队列查询、审核操作 |
| admin/handler | 修改 | 添加审核 API |

## 接口设计

### Phase 1: 用户主页数据

```
GET /api/v1/users/:id/profile
响应：{
  "user": { ... },
  "stats": {
    "post_count": 10,
    "follower_count": 50,
    "following_count": 30,
    "like_received_count": 100
  }
}
```

### Phase 2: 通知分类 + 批量操作

```
GET /api/v1/notifications?type=like&unread=true&page=1&page_size=20
POST /api/v1/notifications/batch/read
请求：{ "ids": [1, 2, 3] }
POST /api/v1/notifications/batch/delete
请求：{ "ids": [1, 2, 3] }
```

### Phase 3: 评论回复通知

无新接口，通知自动触发。

### Phase 4: 搜索增强

```
GET /api/v1/search/users?q=keyword&page=1&page_size=20
GET /api/v1/search/suggestions?q=keyword
```

### Phase 5: 帖子编辑增强

```
POST /api/v1/upload/image
请求：multipart/form-data
响应：{ "url": "https://minio.example.com/xxx.jpg" }
```

### Phase 6: 内容审核队列

```
GET /api/v1/admin/posts/pending?page=1&page_size=20
POST /api/v1/admin/posts/:id/approve
POST /api/v1/admin/posts/:id/reject
请求：{ "reason": "违反社区规定" }
```

## 数据变更

### Phase 1: 用户主页数据

无表结构变更，仅查询逻辑。

### Phase 2: 通知分类

```sql
ALTER TABLE notifications ADD COLUMN type VARCHAR(20) DEFAULT 'system';
-- 枚举值：like, comment, follow, system
```

### Phase 3: 评论回复通知

无表结构变更。

### Phase 4: 搜索增强

无表结构变更，依赖 Meilisearch。

### Phase 5: 帖子编辑增强

```sql
ALTER TABLE posts ALTER COLUMN content TYPE JSONB USING jsonb_build_object('text', content);
```

### Phase 6: 内容审核队列

```sql
ALTER TABLE posts ADD COLUMN status VARCHAR(20) DEFAULT 'approved';
ALTER TABLE posts ADD COLUMN reviewed_at TIMESTAMPTZ;
ALTER TABLE posts ADD COLUMN reviewed_by BIGINT;
ALTER TABLE posts ADD COLUMN reject_reason TEXT;
-- 枚举值：pending, approved, rejected
```

## 验证策略

每个 Phase 完成后执行：

| 层级 | 验证内容 | 工具 |
|------|---------|------|
| L1 | 单元测试通过 | `go test -cover` |
| L3 | BDD 场景通过 | `scripts/verify-bdd.sh` |
| L3 | 生产环境部署 | Docker Compose |

## 风险评估

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 富文本内容 XSS | 高 | 服务端 HTML 过滤 |
| 图片上传大小 | 中 | 限制 5MB，压缩处理 |
| 审核队列积压 | 中 | 添加审核人员管理 |
