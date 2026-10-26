# design：搜索功能

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 搜索引擎 | Meilisearch | 轻量、易部署、支持中文 |
| 索引方式 | 分别索引 | posts 和 users 分开 |
| 同步方式 | 实时同步 | 创建/更新时同步 |

## 接口设计

### 1. 搜索帖子

```
GET /api/v1/search/posts?q=关键词&page=1&size=20

响应：
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "type": "emotion",
        "title": "标题",
        "content": "内容",
        "tags": ["标签"],
        "like_count": 10,
        "created_at": "2026-05-19T16:30:00Z"
      }
    ],
    "page": 1,
    "size": 20,
    "has_more": false
  }
}
```

### 2. 搜索用户

```
GET /api/v1/search/users?q=关键词&page=1&size=20

响应：
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "nickname": "用户昵称",
        "avatar": "/uploads/avatars/1.jpg",
        "bio": "个人简介"
      }
    ],
    "page": 1,
    "size": 20,
    "has_more": false
  }
}
```

## 数据变更

| 操作 | 对象 | 说明 |
|------|------|------|
| 新增 | Meilisearch 索引 | posts 和 users |

## 依赖

- Meilisearch：全文搜索
- PostgreSQL：数据源
