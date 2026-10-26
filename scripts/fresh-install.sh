#!/usr/bin/env bash
#
# 全新安装（不重放历史迁移）—— 「两套逻辑」中的 fresh 路径。
#
#   - 全新空库   → fresh-install.sh（本脚本）：加载 schema baseline 一步建出
#                  当前完整 schema，再把 baseline 覆盖的迁移标记为已应用，
#                  最后只追加 baseline 之后新增的迁移。
#   - 已有库升级 → scripts/migrate.sh：只应用 schema_migrations 里缺的增量迁移。
#
# baseline 由 scripts/verify-baseline.sh 守护，始终 == 重放全部迁移的结果，
# 因此正常情况下「追加」步骤为空。新增迁移后请重新生成 baseline（见 baseline/README.md）。
#
# 连接走标准 PG 环境变量：PGHOST/PGPORT/PGUSER/PGPASSWORD/PGDATABASE。
# 用法：
#   PGHOST=localhost PGUSER=shiyemo PGPASSWORD=xxx PGDATABASE=shiyemo \
#     bash scripts/fresh-install.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIG_DIR="$ROOT/server/migrations"
BASE_DIR="$MIG_DIR/baseline"
BASELINE="$BASE_DIR/schema_baseline.sql"
COVERED="$BASE_DIR/covered.txt"

[ -f "$BASELINE" ] || { echo "❌ baseline 不存在: $BASELINE"; exit 1; }
[ -f "$COVERED" ]  || { echo "❌ covered.txt 不存在: $COVERED"; exit 1; }

echo "[1/3] 加载 schema baseline（一步建出当前 schema）..."
psql -v ON_ERROR_STOP=1 -q -f "$BASELINE"

echo "[2/3] 标记 baseline 覆盖的迁移为已应用..."
psql -v ON_ERROR_STOP=1 -q -c \
  "CREATE TABLE IF NOT EXISTS schema_migrations (filename text PRIMARY KEY, applied_at timestamptz NOT NULL DEFAULT now());"
while IFS= read -r name; do
  [ -z "$name" ] && continue
  psql -v ON_ERROR_STOP=1 -q -c \
    "INSERT INTO schema_migrations (filename) VALUES ('$name') ON CONFLICT DO NOTHING;"
done < "$COVERED"

echo "[3/3] 追加 baseline 之后新增的迁移（不在 covered.txt 里的）..."
applied_tail=0
for f in $(ls "$MIG_DIR"/*.sql | sort); do
  name="$(basename "$f")"
  if ! grep -qxF "$name" "$COVERED"; then
    printf '  apply  %-45s ' "$name"
    psql -v ON_ERROR_STOP=1 -q -f "$f"
    psql -v ON_ERROR_STOP=1 -q -c \
      "INSERT INTO schema_migrations (filename) VALUES ('$name') ON CONFLICT DO NOTHING;"
    echo 'OK'
    applied_tail=$((applied_tail + 1))
  fi
done

echo ""
echo "✓ 全新安装完成：load baseline + 标记 $(wc -l < "$COVERED") 个已覆盖迁移 + 追加 $applied_tail 个新迁移"
