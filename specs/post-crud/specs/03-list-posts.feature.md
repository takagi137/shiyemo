# Feature: 获取帖子列表

## 场景：获取所有帖子

### 测试步骤

1. 发送 GET 请求到 `/api/v1/posts`

### 预期结果

- HTTP 状态码：200
- 响应体包含分页数据

### 验证点

- [x] 返回 list 数组
- [x] 包含 page, size, has_more 字段
- [x] 按时间倒序排列

## 场景：按类型筛选

### 测试步骤

1. 发送 GET 请求到 `/api/v1/posts?type=emotion`

### 预期结果

- 只返回 emotion 类型的帖子

## 场景：分页查询

### 测试步骤

1. 发送 GET 请求到 `/api/v1/posts?page=1&size=2`

### 预期结果

- 返回最多 2 条记录
- has_more 表示是否有下一页
