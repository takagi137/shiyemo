# design：用户注册/登录功能

## 技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 认证方式 | JWT (HS256) | 无状态，单人项目简单可靠 |
| 密码哈希 | bcrypt | Go 社区标准，安全性高 |
| Token 过期 | 7 天 | 平衡安全性和用户体验 |
| 验证码存储 | Redis | 自动过期，读写快 |
| 验证码有效期 | 5 分钟 | 行业标准 |
| 验证码发送频率 | 60 秒/次 | 防刷 |
| MVP 验证码 | Mock 123456 | 省成本，先跑通流程 |

## 模块影响

| 模块 | 改动类型 | 改动说明 |
|------|---------|---------|
| auth/handler.go | 新增 | 注册、登录、发送验证码接口 |
| auth/service.go | 新增 | 认证业务逻辑 |
| auth/repo.go | 新增 | 用户查询、创建 |
| auth/model.go | 新增 | 用户模型定义 |
| common/middleware/jwt.go | 完善 | JWT 生成、解析、验证 |
| common/database/database.go | 完善 | 数据库连接初始化 |
| cmd/main.go | 完善 | 路由注册、服务启动 |

## 接口设计

### 1. 邮箱注册

```
POST /api/v1/auth/register/email
请求：
{
  "email": "user@example.com",
  "password": "a123456",
  "nickname": "测试用户"
}

响应（成功）：
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "<JWT token>",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "nickname": "测试用户"
    }
  }
}

响应（失败 - 邮箱已存在）：
{
  "code": 40201,
  "message": "邮箱已注册",
  "data": null
}
```

### 2. 手机号注册

```
POST /api/v1/auth/register/phone
请求：
{
  "phone": "13800138000",
  "code": "123456",
  "nickname": "测试用户"
}

响应（成功）：
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "<JWT token>",
    "user": {
      "id": 1,
      "phone": "13800138000",
      "nickname": "测试用户"
    }
  }
}
```

### 3. 邮箱登录

```
POST /api/v1/auth/login/email
请求：
{
  "email": "user@example.com",
  "password": "a123456"
}

响应（成功）：
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "<JWT token>",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "nickname": "测试用户"
    }
  }
}
```

### 4. 手机号登录

```
POST /api/v1/auth/login/phone
请求：
{
  "phone": "13800138000",
  "code": "123456"
}
```

### 5. 发送验证码

```
POST /api/v1/auth/sms/send
请求：
{
  "phone": "13800138000"
}

响应（成功）：
{
  "code": 0,
  "message": "验证码已发送",
  "data": null
}

响应（失败 - 发送频繁）：
{
  "code": 40106,
  "message": "验证码发送过于频繁，请60秒后重试",
  "data": null
}
```

## 数据变更

| 操作 | 对象 | 说明 |
|------|------|------|
| 新增 | users 表 | 创建用户表 |
| 新增 | 索引 | phone、email 唯一索引 |

## 依赖

- PostgreSQL：用户数据存储
- Redis：验证码存储、频率限制
- JWT 库：`github.com/golang-jwt/jwt/v5`
- 密码哈希：`golang.org/x/crypto/bcrypt`

## 风险点

1. Mock 验证码模式下，需要在生产环境切换为真实短信服务
2. JWT Secret 必须在生产环境更换为强密码
3. 邮箱格式校验需要覆盖各种边界情况
