# 获取热门标签

## 场景：成功获取热门标签
- API: `GET /api/v1/tags`
- 无需登录

### 请求
```
GET /api/v1/tags
```

### 响应 (200)
```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {"id": 1, "name": "Python", "count": 10},
    {"id": 2, "name": "产品经理", "count": 8}
  ]
}
```
