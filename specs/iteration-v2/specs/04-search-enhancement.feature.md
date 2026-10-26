# feature: 搜索增强

## 场景 1：搜索用户

### Given 存在用户 "张三"
### When 用户请求 GET /api/v1/search/users?q=张三
### Then 返回匹配的用户列表

```json
// 请求
GET /api/v1/search/users?q=张三&page=1&page_size=10

// 预期响应
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": 1,
        "nickname": "张三"
      }
    ],
    "total": 1
  }
}
```

## 场景 2：搜索建议

### Given 存在帖子标题 "Go语言教程"
### When 用户请求 GET /api/v1/search/suggestions?q=Go
### Then 返回搜索建议列表

```json
// 请求
GET /api/v1/search/suggestions?q=Go

// 预期响应
{
  "code": 0,
  "message": "ok",
  "data": {
    "suggestions": ["Go语言教程", "Go语言入门"]
  }
}
```
