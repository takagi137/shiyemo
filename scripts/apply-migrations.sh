#!/usr/bin/env bash
#
# 便携迁移应用器（psql 版）—— 用于 CI / 干净环境的一次性全量应用。
#
# 与 scripts/migrate.sh 的分工：
#   - migrate.sh    ：生产/docker 环境，用 docker exec + schema_migrations 表记录，
#                     幂等可重复运行。
#   - apply-migrations.sh：CI/全新库，直接 psql，一把梭全量应用，任一失败即退出非 0。
#     不需要记录表（CI 每次都是干净库）。
#
# 二者共享同一套排序规则：编号 → create_ 优先 → 文件名。
# 排序错误会导致 008_add_notification_type 排在 008_create_notifications 之前。
#
# 连接参数走标准 PG 环境变量：PGHOST/PGPORT/PGUSER/PGPASSWORD/PGDATABASE。
#
# 用法：
#   PGHOST=localhost PGUSER=shiyemo PGPASSWORD=shiyemo123 PGDATABASE=shiyemo \
#     bash scripts/apply-migrations.sh

set -euo pipefail

MIGRATIONS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/server/migrations"

# 共享排序规则（与 scripts/migrate.sh / setup_test.go 保持一致）
source "$(dirname "${BASH_SOURCE[0]}")/_lib_migrations.sh"

applied=0
for f in $(sorted_migrations "$MIGRATIONS_DIR"); do
  name="$(basename "$f")"
  printf '  apply  %-45s ' "$name"
  # ON_ERROR_STOP=1 让任一语句失败即非 0 退出（set -e 会终止脚本）
  psql -v ON_ERROR_STOP=1 -q -f "$f"
  echo 'OK'
  applied=$((applied + 1))
done

echo ""
echo "✓ 全量应用完成：$applied 个迁移，0 失败"
