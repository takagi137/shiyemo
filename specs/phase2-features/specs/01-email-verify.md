# spec：邮箱验证

## 场景 1：注册后发送验证邮件

### Given 前置条件
- 用户通过邮箱注册成功，获得 JWT Token
- 用户邮箱未验证（email_verified = false）

### When 操作
- 用户调用 `POST /api/v1/auth/verify/send`（携带 Token）

### Then 预期结果
- 响应状态码 200，code=0
- 系统向用户邮箱发送验证邮件
- Redis 中存储验证 Token（key: `verify:email:<token>`，TTL: 24h）
- 邮件内包含验证链接

---

## 场景 2：点击验证链接完成验证

### Given 前置条件
- 用户收到验证邮件，邮件中包含有效 Token

### When 操作
- 用户访问 `GET /api/v1/auth/verify/<token>`

### Then 预期结果
- 响应状态码 200，code=0，message="邮箱验证成功"
- 数据库中 users.email_verified 更新为 true
- Redis 中删除该 Token
- 重复访问同一 Token 返回错误（40001 Token 无效）

---

## 场景 3：验证 Token 过期

### Given 前置条件
- 用户收到验证邮件，但超过 24 小时未操作

### When 操作
- 用户访问 `GET /api/v1/auth/verify/<token>`

### Then 预期结果
- 响应 code=40001，message="验证链接已过期，请重新发送"

---

## 场景 4：重复发送验证邮件

### Given 前置条件
- 用户已发送过验证邮件（60 秒内）

### When 操作
- 用户再次调用 `POST /api/v1/auth/verify/send`

### Then 预期结果
- 响应 code=40002，message="请 60 秒后重试"

---

## 场景 5：已验证用户请求发送

### Given 前置条件
- 用户邮箱已验证（email_verified = true）

### When 操作
- 用户调用 `POST /api/v1/auth/verify/send`

### Then 预期结果
- 响应 code=40201，message="邮箱已验证"
