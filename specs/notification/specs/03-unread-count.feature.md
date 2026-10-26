# 场景：获取未读数量

## 场景 1：获取未读通知数量

### 前置条件
- 用户有 3 条未读通知、2 条已读通知

### 测试步骤
1. 发送 GET 请求到 `/api/v1/notifications/unread-count`

### 预期结果
- 状态码：200
- count 为 3

---

## 场景 2：无未读通知时返回 0

### 前置条件
- 用户所有通知已读

### 测试步骤
1. 发送 GET 请求到 `/api/v1/notifications/unread-count`

### 预期结果
- 状态码：200
- count 为 0
