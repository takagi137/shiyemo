# 反馈管理模块 OpenSpec 设计

> 创建时间: 2026-05-28
> 状态: 进行中

---

## 1. 功能概述

用户可以提交反馈（Bug、建议、其他），管理员可以查看、处理反馈并回复。

---

## 2. 用户故事

### 2.1 用户端

- 作为用户，我可以提交反馈（选择类型、填写标题/内容/联系方式）
- 作为用户，我可以查看我的反馈列表和处理状态
- 作为用户，我可以查看管理员的回复

### 2.2 管理端

- 作为管理员，我可以查看所有反馈列表
- 作为管理员，我可以按状态/类型筛选反馈
- 作为管理员，我可以处理反馈（标记已处理/忽略）
- 作为管理员，我可以回复反馈

---

## 3. API 设计

### 3.1 用户端接口

#### POST /api/v1/feedback

创建反馈

**请求：**
```json
{
  "type": "bug|suggestion|other",
  "title": "必填，200字符以内",
  "content": "必填，反馈详情",
  "contact": "可选，联系方式"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "ok",
  "data": null
}
```

#### GET /api/v1/feedback/mine

获取我的反馈列表

**请求参数：**
- page: 页码（默认1）
- page_size: 每页数量（默认20）

**响应：**
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [...],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

### 3.2 管理端接口

#### GET /api/v1/admin/feedback

获取反馈列表（管理员）

**请求参数：**
- page: 页码
- page_size: 每页数量
- status: 状态筛选（0-待处理 1-已处理 2-已忽略）
- type: 类型筛选（bug/suggestion/other）

**响应：**
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": 1,
        "user_id": 100,
        "user_nickname": "用户昵称",
        "user_email": "user@example.com",
        "type": "bug",
        "title": "标题",
        "content": "内容",
        "contact": "联系方式",
        "status": 0,
        "reply": "",
        "created_at": "2026-05-28T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

#### PUT /api/v1/admin/feedback/:id/status

更新反馈状态（管理员）

**请求：**
```json
{
  "status": 1,
  "reply": "已修复，感谢反馈"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "ok",
  "data": null
}
```

---

## 4. 数据模型

### Feedback 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| user_id | BIGINT | 用户ID，索引 |
| type | VARCHAR(20) | 类型：bug/suggestion/other |
| title | VARCHAR(200) | 标题 |
| content | TEXT | 内容 |
| contact | VARCHAR(100) | 联系方式（可选） |
| status | SMALLINT | 0-待处理 1-已处理 2-已忽略 |
| reply | TEXT | 管理员回复 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

---

## 5. 文件清单

### 后端（server/internal/feedback/）

| 文件 | 内容 | 状态 |
|------|------|------|
| model.go | Feedback 模型定义 | ⏳ |
| repo.go | 数据库操作 | ⏳ |
| service.go | 业务逻辑 | ⏳ |
| handler.go | HTTP 接口 | ⏳ |

### 前端（web/src/app/）

| 文件 | 内容 | 状态 |
|------|------|------|
| admin/feedback/page.tsx | 管理端反馈列表 | ⏳ |
| feedback/page.tsx | 用户端反馈提交页 | ⏳ |

---

## 6. 验证标准

- [ ] 后端：`go build ./...` 通过
- [ ] 后端：`go vet ./...` 无警告
- [ ] 前端：`npm run build` 通过
- [ ] 管理端反馈列表可正常显示
- [ ] 用户端可提交反馈
- [ ] 管理员可处理反馈并回复
