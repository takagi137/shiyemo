# 场景：获取用户帖子列表

## 场景 1：成功获取我的帖子列表

### 前置条件
- 用户已登录
- 用户有多个帖子

### 测试步骤
1. 发送 GET 请求到 `/api/v1/users/:id/posts?page=1&page_size=20`
2. 携带有效的 JWT token

### 预期结果
- 状态码：200
- 返回帖子列表
- 按创建时间倒序排列
- 包含该用户的所有帖子

---

## 场景 2：分页获取帖子

### 前置条件
- 用户有 25 个帖子

### 测试步骤
1. 发送 GET 请求到 `/api/v1/users/:id/posts?page=1&page_size=10`
2. 发送 GET 请求到 `/api/v1/users/:id/posts?page=2&page_size=10`

### 预期结果
- 第一页返回 10 条
- 第二页返回 10 条
- total 为 25
