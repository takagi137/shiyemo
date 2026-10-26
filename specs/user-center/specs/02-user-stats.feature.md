# 场景：获取用户统计数据

## 场景 1：成功获取统计数据

### 前置条件
- 用户已登录
- 用户有帖子、有粉丝、有获赞

### 测试步骤
1. 发送 GET 请求到 `/api/v1/users/:id/stats`
2. 携带有效的 JWT token

### 预期结果
- 状态码：200
- 返回 post_count（帖子数量）
- 返回 like_received（获赞数量）
- 返回 follower_count（粉丝数量）
- 返回 following_count（关注数量）

---

## 场景 2：新用户统计数据为 0

### 前置条件
- 用户是新用户

### 测试步骤
1. 发送 GET 请求到 `/api/v1/users/:id/stats`

### 预期结果
- 状态码：200
- 所有统计字段为 0
