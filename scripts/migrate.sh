#!/usr/bin/env bash
#
# 幂等数据库迁移执行器
#
# 背景：本项目的 migrations/*.sql 一直靠人工 psql 逐个执行，
# 曾因漏执行 016_create_daily_stats.sql 导致新模块路由全部 404。
# 本脚本用 schema_migrations 表记录已应用的文件，可安全重复运行。
#
# 用法：./scripts/migrate.sh [postgres容器名]
#   默认容器名 shiyemo-db，数据库/用户均为 shiyemo。
#
# 注意：迁移文件编号存在重复（005/006/007/008 各两个）且缺号（011/018）。
# 纯字典序会出错：008_add_notification_type.sql 排在 008_create_notifications.sql
# 之前，但前者要 ALTER 后者创建的表 —— 依赖倒置。
# 因此排序规则为「编号 → create_ 优先 → 文件名」。

set -euo pipefail

CONTAINER="${1:-shiyemo-db}"
DB_USER="${DB_USER:-shiyemo}"
DB_NAME="${DB_NAME:-shiyemo}"
MIGRATIONS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/server/migrations"

# 共享排序规则（与 scripts/apply-migrations.sh / setup_test.go 保持一致）
source "$(dirname "${BASH_SOURCE[0]}")/_lib_migrations.sh"

if ! docker ps --format '{{.Names}}' | grep -qx "$CONTAINER"; then
  echo "错误：容器 $CONTAINER 未在运行" >&2
  exit 1
fi

psql_run() {
  docker exec -i "$CONTAINER" psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME" "$@"
}

# 建记录表
psql_run -q -c "CREATE TABLE IF NOT EXISTS schema_migrations (
  filename TEXT PRIMARY KEY,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);"

applied=0
skipped=0

for f in $(sorted_migrations "$MIGRATIONS_DIR"); do
  name="$(basename "$f")"

  if psql_run -tAc "SELECT 1 FROM schema_migrations WHERE filename='$name'" | grep -q 1; then
    printf '  skip   %s\n' "$name"
    skipped=$((skipped + 1))
    continue
  fi

  printf '  apply  %s ... ' "$name"
  # 单个事务内执行迁移 + 记录，失败则整体回滚
  if { echo 'BEGIN;'; cat "$f"; \
       echo "INSERT INTO schema_migrations(filename) VALUES ('$name');"; \
       echo 'COMMIT;'; } | psql_run -q 2>/tmp/migrate_err_$$; then
    echo 'OK'
    applied=$((applied + 1))
  else
    echo 'FAILED'
    echo '--- 错误详情 ---' >&2
    cat /tmp/migrate_err_$$ >&2
    rm -f /tmp/migrate_err_$$
    exit 1
  fi
  rm -f /tmp/migrate_err_$$
done

echo
echo "迁移完成：新应用 $applied 个，跳过 $skipped 个（已应用过）"
