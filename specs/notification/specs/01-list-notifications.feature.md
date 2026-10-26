# 场景：获取通知列表

## 场景 1：成功获取通知列表

### 前置条件
- 用户已登录
- 用户有多条通知

### 测试步骤
1. 发送 GET 请求到 `/api/v1/notifications?page=1&page_size=20`
2. 携带有效的 JWT token

### 预期结果
- 状态码：200
- 响应体包含通知列表
- 每条通知包含 id、type、actor、content、is_read、created_at
- 按创建时间倒序排列

---

## 场景 2：分页获取通知

### 前置条件
- 用户有 25 条通知

### 测试步骤
1. 发送 GET 请求到 `/api/v1/notifications?page=1&page_size=10`
2. 发送 GET 请求到 `/api/v1/notifications?page=2&page_size=10`

### 预期结果
- 第一页返回 10 条
- 第二页返回 10 条
- total 为 25

---

## 场景 3：无通知时返回空列表

### 前置条件
- 用户是新用户，无通知

### 测试步骤
1. 发送 GET 请求到 `/api/v1/notifications`

### 预期结果
- 状态码：200
- items 为空数组
- total 为 0
