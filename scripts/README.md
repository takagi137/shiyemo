# 开发脚本

## 快速开始

### Windows

```bash
# 启动开发环境
scripts\dev.bat
```

### Linux/Mac

```bash
# 给脚本执行权限
chmod +x scripts/*.sh

# 启动开发环境
./scripts/dev.sh
```

## 测试功能

### 运行认证功能测试

```bash
# Linux/Mac
chmod +x scripts/test-auth.sh
./scripts/test-auth.sh

# Windows (需要 Git Bash 或 WSL)
bash scripts/test-auth.sh
```

## 脚本说明

### dev.sh / dev.bat

启动完整的开发环境，包括：
- PostgreSQL
- Redis
- Meilisearch
- 后端 API 服务

### test-auth.sh

测试用户认证功能的完整流程：
1. 邮箱注册
2. 邮箱登录
3. 手机号注册（使用 mock 验证码 123456）
4. 手机号登录
5. 访问受保护接口
6. 重复注册测试
7. 密码错误测试
8. 无 token 访问测试

## 手动测试

### 邮箱注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register/email \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "a123456",
    "nickname": "测试用户"
  }'
```

### 手机号注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register/phone \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "code": "123456",
    "nickname": "手机用户"
  }'
```

### 邮箱登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login/email \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "a123456"
  }'
```

### 手机号登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login/phone \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "code": "123456"
  }'
```

### 获取用户信息（需要 token）

```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer <your-token>"
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| APP_PORT | 8080 | API 服务端口 |
| APP_MODE | debug | 运行模式 (debug/release) |
| DB_HOST | localhost | PostgreSQL 主机 |
| DB_PORT | 5432 | PostgreSQL 端口 |
| DB_USER | shiyemo | PostgreSQL 用户 |
| DB_PASSWORD | shiyemo123 | PostgreSQL 密码 |
| DB_NAME | shiyemo | 数据库名 |
| REDIS_HOST | localhost | Redis 主机 |
| REDIS_PORT | 6379 | Redis 端口 |
| REDIS_PASSWORD | (空) | Redis 密码 |
| JWT_SECRET | shiyemo-jwt-secret... | JWT 密钥 (生产环境必须修改) |
