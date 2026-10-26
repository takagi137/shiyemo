#!/usr/bin/env bash
# ============================================================
# 迁移文件排序共享库
# 单一信息源，供以下脚本 source：
#   - scripts/migrate.sh          (生产/docker，幂等)
#   - scripts/apply-migrations.sh (CI/干净库，一把梭)
#
# 规则：编号 → create_ 优先 → 文件名。
#   类型秩 0 = create_（建表），1 = 其它（add_/alter_ 等改表）。
#   同编号下先建表再改表，避免 008_add_notification_type 排在 008_create_notifications 之前。
#
# 用法（在被 source 脚本里）：
#   source "$(dirname "${BASH_SOURCE[0]}")/_lib_migrations.sh"
#   sorted_migrations "$MIGRATIONS_DIR"   # 每行一个完整路径
# ============================================================

# 防御性：本库不应被直接执行
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "错误：${BASH_SOURCE[0]} 是共享库，请用 source 引入" >&2
    exit 1
fi

# sorted_migrations <migrations_dir>
#   把目录下所有 *.sql 按「编号 → 类型秩 → 文件名」排序，每行打印完整路径。
sorted_migrations() {
    local dir="$1"
    local base num rank path
    for path in "$dir"/*.sql; do
        base="$(basename "$path")"
        num="${base%%_*}"
        case "$base" in
            *_create_*) rank=0 ;;
            *)          rank=1 ;;
        esac
        printf '%s %s %s %s\n' "$num" "$rank" "$base" "$path"
    done | sort -k1,1 -k2,2n -k3,3 | awk '{print $4}'
}