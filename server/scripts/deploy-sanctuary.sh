#!/bin/bash
# 静思空间模块部署脚本（服务器端版本）
# 用途：在服务器上执行数据库迁移 + 重启 Docker 容器
# 使用：ssh root@<server> 'bash -s' < scripts/deploy-sanctuary.sh
#
# 可选环境变量：REMOTE_PORT=8090（默认；:80/:443 已被 Caddy 占用）
#
# ⚠️ 仓库不写死服务器地址，避免 IP 在 GitHub 公开暴露。

set -e

REMOTE_PORT="${REMOTE_PORT:-8090}"

echo "=== 失业么 - 静思空间部署（服务器端） ==="
echo "时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo "入口端口: ${REMOTE_PORT}"

# 1. 检查 Docker 状态
echo ""
echo "[1/4] 检查 Docker 状态..."
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行，尝试启动..."
    systemctl start docker
    sleep 3
fi
echo "✅ Docker 运行中"

# 2. 检查容器状态
echo ""
echo "[2/4] 检查容器状态..."
cd /opt/shiyemo/deploy
docker compose -f docker-compose.yml -f docker-compose.server.yml ps

# 3. 执行数据库迁移
echo ""
echo "[3/4] 执行数据库迁移..."
# 进入 PostgreSQL 容器执行 SQL
echo "  → 迁移 020_create_journals.sql..."
docker compose -f docker-compose.yml -f docker-compose.server.yml exec -T postgres psql -U shiyemo -d shiyemo -f /dev/stdin < /opt/shiyemo/deploy/migrations/020_create_journals.sql 2>&1 || echo "  ⚠️ 020 可能已存在，继续..."

echo "  → 迁移 021_create_readings.sql..."
docker compose -f docker-compose.yml -f docker-compose.server.yml exec -T postgres psql -U shiyemo -d shiyemo -f /dev/stdin < /opt/shiyemo/deploy/migrations/021_create_readings.sql 2>&1 || echo "  ⚠️ 021 可能已存在，继续..."

# 备用方案：直接通过 SQL 命令
echo "  → 验证表是否创建..."
docker compose -f docker-compose.yml -f docker-compose.server.yml exec -T postgres psql -U shiyemo -d shiyemo -c "\dt journals" 2>&1
docker compose -f docker-compose.yml -f docker-compose.server.yml exec -T postgres psql -U shiyemo -d shiyemo -c "\dt readings" 2>&1

# 4. 重启应用容器（重新构建镜像以包含最新代码）
echo ""
echo "[4/4] 重启应用容器..."
cd /opt/shiyemo/deploy

# 拉取最新代码（如果有 git）
if [ -d "/opt/shiyemo/.git" ]; then
    echo "  → 拉取最新代码..."
    cd /opt/shiyemo && git pull
fi

# 重新构建并启动
cd /opt/shiyemo/deploy
echo "  → 重新构建镜像..."
docker compose -f docker-compose.yml -f docker-compose.server.yml build --no-cache

echo "  → 重启容器..."
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d

# 等待服务启动
echo "  → 等待服务启动（10s）..."
sleep 10

# 5. 验证
echo ""
echo "=== 部署验证 ==="
echo "→ 后端健康检查..."
curl -s "http://localhost:${REMOTE_PORT}/api/v1/health" | head -1 || echo "❌ 后端未响应"

echo "→ 前端检查..."
curl -s "http://localhost:${REMOTE_PORT}" | head -1 | cut -c1-80 || echo "❌ 前端未响应"

echo ""
echo "=== 部署完成 ==="
echo "本机访问:"
echo "  前端: http://localhost:${REMOTE_PORT}"
echo "  后端: http://localhost:${REMOTE_PORT}/api/v1"
echo "  静思空间: http://localhost:${REMOTE_PORT}/sanctuary"
