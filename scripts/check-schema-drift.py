#!/usr/bin/env python3
"""Schema 漂移检测：对比 Go 模型字段与数据库实际列。

背景
----
本项目的迁移是手写 SQL、与 Go 模型分离维护，历史上多次脱节：
- 模型有 email_verified / last_work_date 字段，但无建表/加列迁移
- 4 张表（daily_checkins/email_verifications/sensitive_words/survival_tips）
  有模型和路由，却从无建表迁移
- 迁移 012 引用列名 read，实际列为 is_read
这些在「全新部署」时才暴露，而 CI 从不跑真实迁移，导致潜伏两个月。

本脚本把「模型与库是否一致」变成可自动化的红/绿判断，供 CI 卡点使用。

用法
----
    python scripts/check-schema-drift.py \
        --host localhost --port 5432 --user shiyemo --password shiyemo123 --db shiyemo

    # 或喂一个已导出的列清单文件（每行 "table:col1 col2 ..."）
    python scripts/check-schema-drift.py --cols-file /tmp/cols.txt

退出码非 0 表示存在漂移（表缺失或列缺失），供 CI 判定失败。

限制
----
列名由 Go 结构体 tag 推导（column: 优先，否则字段名转 snake_case）。
这是静态启发式，个别特殊映射可能需要在 KNOWN_SKIP 中登记。
"""
import argparse
import glob
import os
import re
import subprocess
import sys

# 相对本脚本定位 server/internal
REPO_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
INTERNAL = os.path.join(REPO_ROOT, "server", "internal")

# 已知无需校验的表（如纯连接表由 gorm 隐式管理、或无独立模型）
KNOWN_SKIP = {"schema_migrations"}


def to_snake(name: str) -> str:
    """Go 字段名 → 数据库列名（与 GORM 默认 NamingStrategy 一致）。"""
    if name in ("ID", "URL", "SQL"):
        return name.lower()
    name = re.sub(r"(ID|URL|SQL)$", lambda m: "_" + m.group(1).lower(), name)
    name = re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", name)
    return name.lower().replace("__", "_")


def parse_models() -> dict:
    """扫描 server/internal 下所有 *model.go，返回 {table: (file, {cols})}。"""
    models = {}
    for f in glob.glob(os.path.join(INTERNAL, "**", "model.go"), recursive=True):
        src = open(f, encoding="utf-8").read()
        for m in re.finditer(r"type (\w+) struct \{(.*?)\n\}", src, re.S):
            name, body = m.group(1), m.group(2)
            tn = re.search(
                r"func \(" + name + r"\) TableName\(\) string \{\s*return \"(\w+)\"",
                src,
            )
            if not tn:
                continue
            cols = set()
            for line in body.split("\n"):
                line = line.strip()
                if not line or line.startswith("//"):
                    continue
                fm = re.match(r"(\w+)\s+[\*\[\]\w\.]+\s*`(.*)`", line)
                if not fm:
                    continue
                fname, tags = fm.group(1), fm.group(2)
                if 'gorm:"-"' in tags:
                    continue
                cm = re.search(r"column:(\w+)", tags)
                cols.add(cm.group(1) if cm else to_snake(fname))
            models[tn.group(1)] = (os.path.relpath(f, REPO_ROOT), cols)
    return models


def dump_columns_via_psql(args) -> dict:
    """连库导出 {table: {cols}}。"""
    sql = (
        "SELECT table_name || ':' || string_agg(column_name, ' ' ORDER BY column_name) "
        "FROM information_schema.columns WHERE table_schema='public' "
        "GROUP BY table_name ORDER BY table_name;"
    )
    env = dict(os.environ, PGPASSWORD=args.password)
    out = subprocess.check_output(
        ["psql", "-h", args.host, "-p", str(args.port), "-U", args.user,
         "-d", args.db, "-tA", "-c", sql],
        env=env, text=True,
    )
    return parse_cols_text(out)


def parse_cols_text(text: str) -> dict:
    db = {}
    for line in text.splitlines():
        line = line.strip()
        if ":" in line:
            t, c = line.split(":", 1)
            db[t.strip()] = set(c.split())
    return db


def main() -> int:
    p = argparse.ArgumentParser()
    p.add_argument("--host", default="localhost")
    p.add_argument("--port", default="5432")
    p.add_argument("--user", default="shiyemo")
    p.add_argument("--password", default="shiyemo123")
    p.add_argument("--db", default="shiyemo")
    p.add_argument("--cols-file", help="已导出的列清单文件（每行 table:cols）")
    args = p.parse_args()

    db = parse_cols_text(open(args.cols_file, encoding="utf-8").read()) \
        if args.cols_file else dump_columns_via_psql(args)
    models = parse_models()

    print("=== Go 模型字段 vs 数据库实际列 ===\n")
    issues = 0
    for table, (f, cols) in sorted(models.items()):
        if table in KNOWN_SKIP:
            continue
        if table not in db:
            print(f"[表缺失] {table}   <- {f}")
            issues += 1
            continue
        missing = cols - db[table]
        if missing:
            print(f"[缺列]   {table}: {', '.join(sorted(missing))}   <- {f}")
            issues += 1
        else:
            print(f"[OK]     {table}")

    print(f"\n共 {issues} 处不一致")
    if issues:
        print("\n✗ Schema 漂移：Go 模型与数据库不一致，全新部署会失败。")
        return 1
    print("\n✓ Schema 与模型一致")
    return 0


if __name__ == "__main__":
    sys.exit(main())
