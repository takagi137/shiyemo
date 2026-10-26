#!/bin/bash

# 开发环境启动脚本

set -e

echo "🚀 启动失业么开发环境..."

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行，请先启动 Docker"
    exit 1
fi

# 启动依赖服务
echo "📦 启动 PostgreSQL, Redis, Meilisearch..."
cd deploy
docker-compose up -d
cd ..

# 等待 PostgreSQL 就绪
echo "⏳ 等待 PostgreSQL 就绪..."
sleep 5

# 运行数据库迁移
echo "🗄️ 运行数据库迁移..."
PGPASSWORD=shiyemo123 psql -h localhost -U shiyemo -d shiyemo -f server/migrations/001_create_users.sql || echo "⚠️ 迁移可能已存在，继续..."

# 启动后端服务
echo "🔧 启动后端服务..."
cd server
source ../deploy/.env
go run cmd/main.go &
SERVER_PID=$!
cd ..

echo "✅ 开发环境已启动！"
echo ""
echo "📌 服务地址："
echo "   - API: http://localhost:8080"
echo "   - PostgreSQL: localhost:5432"
echo "   - Redis: localhost:6379"
echo "   - Meilisearch: localhost:7700"
echo ""
echo "📌 测试接口："
echo "   - 邮箱注册: POST http://localhost:8080/api/v1/auth/register/email"
echo "   - 手机号注册: POST http://localhost:8080/api/v1/auth/register/phone"
echo "   - 邮箱登录: POST http://localhost:8080/api/v1/auth/login/email"
echo "   - 手机号登录: POST http://localhost:8080/api/v1/auth/login/phone"
echo ""
echo "按 Ctrl+C 停止服务"

# 捕获退出信号
trap "echo '🛑 停止服务...'; kill $SERVER_PID; cd deploy && docker-compose down; exit" SIGINT SIGTERM

# 等待后端进程
wait $SERVER_PID
