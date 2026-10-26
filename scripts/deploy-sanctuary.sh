#!/bin/bash
# 静思空间模块部署脚本
# 用法: SERVER_HOST=<your.server.ip> bash scripts/deploy-sanctuary.sh
# 前提: 能 SSH 到服务器（默认 root，密码走 ssh-agent 或 ~/.ssh/config）
#
# 必需环境变量：
#   SERVER_HOST  目标服务器 IP 或域名
# 可选环境变量：
#   SERVER_USER  SSH 用户（默认 root）
#   REMOTE_PORT  入口端口（默认 8090，因该服务器 :80/:443 被 Caddy 占用）
#
# ⚠️ 仓库不写死服务器地址，避免 IP 在 GitHub 公开暴露。

set -e

if [ -z "$SERVER_HOST" ]; then
    echo "❌ 缺少 SERVER_HOST 环境变量" >&2
    echo "   用法: SERVER_HOST=<ip> bash scripts/deploy-sanctuary.sh" >&2
    exit 1
fi

SERVER_USER="${SERVER_USER:-root}"
SERVER="${SERVER_USER}@${SERVER_HOST}"
REMOTE_DIR="/opt/shiyemo"
REMOTE_PORT="${REMOTE_PORT:-8090}"

echo "=== 静思空间部署脚本 ==="
echo "目标服务器: $SERVER (端口 ${REMOTE_PORT})"
echo ""

# 1. 检查 SSH 连接
echo "[1/5] 检查 SSH 连接..."
if ! ssh -o ConnectTimeout=5 -o BatchMode=yes "$SERVER" "echo ok" 2>/dev/null; then
    echo "❌ 无法连接到 $SERVER，请检查 SSH 密钥和网络"
    exit 1
fi
echo "✅ SSH 连接正常"

# 2. 上传代码
echo "[2/5] 同步代码到服务器..."
cd "$(dirname "$0")/.."
tar czf /tmp/shiyemo-sanctuary.tar.gz \
    --exclude='web/node_modules' \
    --exclude='web/.next' \
    --exclude='server/uploads' \
    --exclude='server/*.exe' \
    --exclude='server/*_test.exe' \
    --exclude='server/cover*.out' \
    -C .. shiyemo

scp /tmp/shiyemo-sanctuary.tar.gz "$SERVER:$REMOTE_DIR/"
ssh "$SERVER" "cd $REMOTE_DIR && tar xzf shiyemo-sanctuary.tar.gz && rm shiyemo-sanctuary.tar.gz"
echo "✅ 代码同步完成"

# 3. 执行数据库迁移
echo "[3/5] 执行数据库迁移..."
ssh "$SERVER" "cd $REMOTE_DIR/shiyemo/server && \
    docker exec -i shiyemo-db psql -U shiyemo -d shiyemo < migrations/020_create_journals.sql && \
    echo '✅ 020 journals 表创建完成' && \
    docker exec -i shiyemo-db psql -U shiyemo -d shiyemo < migrations/021_create_readings.sql && \
    echo '✅ 021 readings 表创建完成'"
echo "✅ 数据库迁移完成"

# 4. 重建后端镜像
echo "[4/5] 重建 Docker 镜像..."
ssh "$SERVER" "cd $REMOTE_DIR/shiyemo/deploy && \
    docker compose build api && \
    docker compose build web && \
    echo '✅ 镜像重建完成'"
echo "✅ 镜像重建完成"

# 5. 重启服务
echo "[5/5] 重启服务..."
ssh "$SERVER" "cd $REMOTE_DIR/shiyemo/deploy && \
    docker compose -f docker-compose.yml -f docker-compose.server.yml up -d && \
    sleep 3 && \
    docker compose ps"
echo "✅ 服务重启完成"

echo ""
echo "=== 部署完成 ==="
echo "前端: http://${SERVER_HOST}:${REMOTE_PORT}"
echo "后端: http://${SERVER_HOST}:${REMOTE_PORT}/api/v1"
echo "静思空间: http://${SERVER_HOST}:${REMOTE_PORT}/sanctuary"
echo ""
echo "验证步骤:"
echo "  curl http://${SERVER_HOST}:${REMOTE_PORT}/api/v1/health"
echo "  curl http://${SERVER_HOST}:${REMOTE_PORT}/api/v1/readings/today"
