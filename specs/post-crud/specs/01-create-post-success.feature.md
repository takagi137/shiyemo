# Feature: 创建帖子成功

## 场景：创建情绪帖子

### 前置条件
- 用户已登录（ID=1）

### 测试步骤

1. 发送 POST 请求到 `/api/v1/posts`
2. 请求头携带 `Authorization: Bearer <token>`
3. 请求体：
```json
{
  "type": "emotion",
  "title": "失业第30天，我开始焦虑了",
  "content": "今天是失业的第30天，投了很多简历，但是面试机会很少。",
  "images": [],
  "tags": ["失业", "焦虑"]
}
```

### 预期结果

- HTTP 状态码：200
- 响应体：
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "user_id": 1,
    "type": "emotion",
    "title": "失业第30天，我开始焦虑了",
    "content": "今天是失业的第30天...",
    "images": [],
    "tags": ["失业", "焦虑"],
    "like_count": 0,
    "comment_count": 0,
    "status": 1,
    "created_at": "2026-05-19T16:30:00Z",
    "updated_at": "2026-05-19T16:30:00Z"
  }
}
```

### 验证点

- [x] 返回完整帖子信息
- [x] user_id 为当前用户
- [x] like_count 和 comment_count 为 0
- [x] status 为 1（已发布）

## 场景：创建认知帖子

### 测试步骤

1. 创建 type="insight" 的帖子

### 预期结果

- 创建成功
- type 为 "insight"

## 场景：创建案例帖子

### 测试步骤

1. 创建 type="case" 的帖子

### 预期结果

- 创建成功
- type 为 "case"

## 场景：创建行动帖子

### 测试步骤

1. 创建 type="action" 的帖子

### 预期结果

- 创建成功
- type 为 "action"

## 场景：未登录创建帖子失败

### 测试步骤

1. 不携带 token 发送 POST 请求

### 预期结果

- HTTP 状态码：401
- 响应体：
```json
{
  "code": 40101,
  "message": "未登录",
  "data": null
}
```

## 场景：无效帖子类型失败

### 测试步骤

1. 创建 type="invalid" 的帖子

### 预期结果

- HTTP 状态码：400
- 响应体包含参数错误信息
