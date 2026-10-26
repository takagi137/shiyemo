@echo off
REM 开发环境启动脚本 (Windows)

echo 🚀 启动失业么开发环境...

REM 检查 Docker 是否运行
docker info >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker 未运行，请先启动 Docker
    exit /b 1
)

REM 启动依赖服务
echo 📦 启动 PostgreSQL, Redis, Meilisearch...
cd deploy
docker-compose up -d
cd ..

REM 等待 PostgreSQL 就绪
echo ⏳ 等待 PostgreSQL 就绪...
timeout /t 5 /nobreak >nul

REM 运行数据库迁移
echo 🗄️ 运行数据库迁移...
set PGPASSWORD=shiyemo123
psql -h localhost -U shiyemo -d shiyemo -f server\migrations\001_create_users.sql
if errorlevel 1 (
    echo ⚠️ 迁移可能已存在，继续...
)

REM 启动后端服务
echo 🔧 启动后端服务...
cd server
set APP_PORT=8080
set APP_MODE=debug
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=shiyemo
set DB_PASSWORD=shiyemo123
set DB_NAME=shiyemo
set REDIS_HOST=localhost
set REDIS_PORT=6379
set JWT_SECRET=shiyemo-jwt-secret-change-in-production
go run cmd\main.go
cd ..

echo ✅ 开发环境已启动！
echo.
echo 📌 服务地址：
echo    - API: http://localhost:8080
echo    - PostgreSQL: localhost:5432
echo    - Redis: localhost:6379
echo    - Meilisearch: localhost:7700
echo.
echo 📌 测试接口：
echo    - 邮箱注册: POST http://localhost:8080/api/v1/auth/register/email
echo    - 手机号注册: POST http://localhost:8080/api/v1/auth/register/phone
echo    - 邮箱登录: POST http://localhost:8080/api/v1/auth/login/email
echo    - 手机号登录: POST http://localhost:8080/api/v1/auth/login/phone
