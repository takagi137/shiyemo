# Schema Baseline —— 两套部署逻辑

本目录实现「**全新安装不重放历史迁移**」的 baseline 机制，把部署分成两条路：

| 场景 | 脚本 | 做什么 |
|------|------|--------|
| **全新空库** | `scripts/fresh-install.sh` | 加载 `schema_baseline.sql` 一步建出当前完整 schema + 把 `covered.txt` 里的迁移标记为已应用 + 追加 baseline 之后的新迁移 |
| **已有库升级** | `scripts/migrate.sh` | 只应用 `schema_migrations` 里缺的增量迁移 |
| **CI 守护** | `scripts/verify-baseline.sh` | 校验 baseline == 重放全部迁移的结果（防漂移） |

> 为什么：迁移文件是本项目 schema 的唯一真相源（手写 SQL、非 AutoMigrate）。
> 全新库若重放全部历史迁移，会越来越慢、老迁移还会腐烂。baseline 让全新安装
> 一步到位，同时用 `verify-baseline.sh` 保证 baseline 永远等于迁移链的结果。

## 文件

- `schema_baseline.sql` —— 当前完整 schema 快照（`pg_dump --schema-only`，去掉随机 `\restrict` 行）。
- `covered.txt` —— 此 baseline 覆盖的迁移文件清单（= 生成 baseline 时的全部迁移）。

## 何时要重新生成 baseline

**每次新增/修改迁移后**都要重新生成，否则 `verify-baseline.sh`（CI）会红。
生成步骤（在干净 PG15 上，与生产版本对齐）：

```bash
# 1. 起干净库并应用全部迁移
PGHOST=localhost PGUSER=shiyemo PGPASSWORD=xxx PGDATABASE=shiyemo \
  bash scripts/apply-migrations.sh

# 2. 导出 schema 快照，去掉随机 \restrict 行 + 环境相关的版本注释行
pg_dump --schema-only --no-owner --no-privileges --no-tablespaces \
  | grep -vE '(restrict|^-- Dumped (from|by) )' > server/migrations/baseline/schema_baseline.sql

# 3. 刷新覆盖清单
ls server/migrations/*.sql | xargs -n1 basename | sort \
  > server/migrations/baseline/covered.txt

# 4. 自检
bash scripts/verify-baseline.sh   # 应输出「完全一致」
```

## 注意

- `schema_baseline.sql` 放在 `migrations/baseline/` 子目录，**不会**被
  `apply-migrations.sh`（`migrations/*.sql` 非递归 glob）当成迁移误应用。
- baseline 必须在**干净 PG15** 上生成（与生产版本一致），避免 pg_dump 版本差异。
