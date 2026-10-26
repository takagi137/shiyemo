# spec：帖子编辑 - 成功场景

### Given 用户已登录，是帖子作者
- 用户 ID = 1
- 帖子 ID = 1，作者 ID = 1

### When 用户提交编辑请求
- PUT /api/v1/posts/1
- 请求体：{ "title": "新标题", "content": "新内容", "type": "insight", "tags": ["AI"] }

### Then 帖子内容更新成功
- 响应 code = 0
- 帖子标题、内容、类型、标签已更新

---

# spec：帖子编辑 - 非作者

### Given 用户已登录，不是帖子作者
- 用户 ID = 2
- 帖子 ID = 1，作者 ID = 1

### When 用户尝试编辑帖子
- PUT /api/v1/posts/1

### Then 返回无权限错误
- 响应 code = 40303
