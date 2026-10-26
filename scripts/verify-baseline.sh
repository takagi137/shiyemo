#!/usr/bin/env bash
#
# baseline 守护 —— 确保「两套逻辑」不漂移。
#
# 校验两件事：
#   1. covered.txt == 当前 server/migrations/*.sql 清单
#   2. schema_baseline.sql == 「重放全部迁移」后的 pg_dump 结果
# 任一不符即退出非 0（说明新增/改了迁移却没重新生成 baseline）。
#
# 需要一个干净空 PG（走标准 PG 环境变量），本脚本会往里应用全部迁移。
# 适合放进 CI 的 fresh-deploy job。
# 用法：
#   PGHOST=localhost PGUSER=shiyemo PGPASSWORD=xxx PGDATABASE=shiyemo \
#     bash scripts/verify-baseline.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIG_DIR="$ROOT/server/migrations"
BASE_DIR="$MIG_DIR/baseline"
BASELINE="$BASE_DIR/schema_baseline.sql"
COVERED="$BASE_DIR/covered.txt"

# 1) covered.txt 与当前迁移清单一致？
current="$(ls "$MIG_DIR"/*.sql | xargs -n1 basename | sort)"
if ! diff <(printf '%s\n' "$current") "$COVERED" >/dev/null; then
  echo "❌ covered.txt 与当前迁移清单不一致 —— 新增/删除迁移后请重新生成 baseline"
  diff <(printf '%s\n' "$current") "$COVERED" | head
  exit 1
fi
echo "✓ covered.txt 与迁移清单一致（$(wc -l < "$COVERED") 个）"

# 2) 重放全部迁移到干净库 → pg_dump → 与 baseline 比对
# 归一化：去掉 pg_dump 的随机 \restrict 行 + 环境相关的版本注释行（-- Dumped from/by ...）
normalize() { grep -vE '(restrict|^-- Dumped (from|by) )'; }

bash "$ROOT/scripts/apply-migrations.sh" >/dev/null
pg_dump --schema-only --no-owner --no-privileges --no-tablespaces | normalize > /tmp/replay_schema.sql

if diff <(normalize < "$BASELINE") /tmp/replay_schema.sql >/dev/null; then
  echo "✅ baseline 与「重放全部迁移」产出的 schema 完全一致"
else
  echo "❌ baseline 已过期：与重放迁移的结果不符 —— 请重新生成 baseline（见 baseline/README.md）"
  diff <(grep -v restrict "$BASELINE") /tmp/replay_schema.sql | head -40
  exit 1
fi
