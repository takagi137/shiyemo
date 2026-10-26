# specs：二次开发 - BDD 场景

## Phase 1: 用户主页数据

### spec-01: 获取用户主页数据 - 成功

```markdown
### Given 用户已登录，存在目标用户
### When 用户请求 GET /api/v1/users/:id/profile
### Then 返回用户基本信息 + 统计数据
    And post_count >= 0
    And follower_count >= 0
    And following_count >= 0
    And like_received_count >= 0
```

### spec-02: 获取用户主页数据 - 不存在

```markdown
### Given 用户已登录
### When 用户请求 GET /api/v1/users/999999/profile
### Then 返回 40203 用户不存在错误
```

---

## Phase 2: 通知分类 + 批量操作

### spec-03: 按类型筛选通知

```markdown
### Given 用户已登录，存在点赞类型通知
### When 用户请求 GET /api/v1/notifications?type=like
### Then 只返回点赞类型通知
```

### spec-04: 批量标记已读

```markdown
### Given 用户已登录，存在未读通知 [1, 2, 3]
### When 用户请求 POST /api/v1/notifications/batch/read { ids: [1, 2, 3] }
### Then 通知 1, 2, 3 标记为已读
    And 返回成功
```

### spec-05: 批量删除通知

```markdown
### Given 用户已登录，存在通知 [1, 2, 3]
### When 用户请求 POST /api/v1/notifications/batch/delete { ids: [1, 2, 3] }
### Then 通知 1, 2, 3 被删除
    And 返回成功
```

---

## Phase 3: 评论回复通知

### spec-06: 评论回复时发送通知

```markdown
### Given 用户 A 发表了评论
### When 用户 B 回复该评论
### Then 用户 A 收到通知
    And 通知类型为 comment
    And 通知内容包含回复信息
```

---

## Phase 4: 搜索增强

### spec-07: 搜索用户

```markdown
### Given 存在用户 "张三"
### When 用户请求 GET /api/v1/search/users?q=张三
### Then 返回匹配的用户列表
```

### spec-08: 搜索建议

```markdown
### Given 存在帖子标题 "Go语言教程"
### When 用户请求 GET /api/v1/search/suggestions?q=Go
### Then 返回搜索建议列表
```

---

## Phase 5: 帖子编辑增强

### spec-09: 上传图片

```markdown
### Given 用户已登录
### When 用户上传图片到 POST /api/v1/upload/image
### Then 返回图片 URL
    And 图片可访问
```

### spec-10: 发布富文本帖子

```markdown
### Given 用户已登录
### When 用户发布包含图片的帖子
### Then 帖子创建成功
    And 内容包含图片 URL
```

---

## Phase 6: 内容审核队列

### spec-11: 查询待审核帖子

```markdown
### Given 管理员已登录，存在待审核帖子
### When 管理员请求 GET /api/v1/admin/posts/pending
### Then 返回待审核帖子列表
```

### spec-12: 审核通过帖子

```markdown
### Given 管理员已登录，存在待审核帖子
### When 管理员请求 POST /api/v1/admin/posts/:id/approve
### Then 帖子状态变为 approved
    And 返回成功
```

### spec-13: 审核拒绝帖子

```markdown
### Given 管理员已登录，存在待审核帖子
### When 管理员请求 POST /api/v1/admin/posts/:id/reject { reason: "违规" }
### Then 帖子状态变为 rejected
    And reject_reason 为 "违规"
    And 返回成功
```
