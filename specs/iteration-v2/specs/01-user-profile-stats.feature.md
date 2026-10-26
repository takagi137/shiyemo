# feature: 获取用户主页数据

## 场景 1：获取用户主页数据 - 成功

### Given 用户已登录，存在目标用户
### When 用户请求 GET /api/v1/users/:id/profile
### Then 返回用户基本信息 + 统计数据

```json
// 请求
GET /api/v1/users/1/profile
Authorization: Bearer <token>

// 预期响应
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "nickname": "测试用户",
    "post_count": 5,
    "follower_count": 10,
    "following_count": 3,
    "like_received_count": 20
  }
}
```

## 场景 2：获取用户主页数据 - 用户不存在

### Given 用户已登录
### When 用户请求 GET /api/v1/users/999999/profile
### Then 返回 40203 用户不存在错误

```json
// 请求
GET /api/v1/users/999999/profile
Authorization: Bearer <token>

// 预期响应
{
  "code": 40203,
  "message": "用户不存在"
}
```
