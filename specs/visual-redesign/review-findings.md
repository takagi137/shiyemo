# 审查遗留清单：visual-redesign 分支

> 来源：2026-07-21 四维度 SDD 审查（code-review-sdd v2.4），总分 73.9 / 100
> 处置决策：P0 当场修复，其余非 bug 项仅记录，不在本次合入阻塞
> 审查范围：`master...worktree-visual-redesign`，18 个提交 / 57 文件 / +2820 −464

---

## 已修复（P0）

| 编号 | 问题 | 处置 |
|------|------|------|
| H1 | `check-design.mjs` 是白名单机制，新增页面默认免检，已实证失效 | 改为全量扫描 + 遗留基线棘轮（`design-baseline.json`），新增文件零容忍；已用探针文件反向验证能报红 |
| H1x | `admin/settings/page.tsx`（本分支新建）含 23 处原子色，含 `text-blue-600` | 全部迁移至语义令牌，实测归零 |
| H3 | 后端改动零执行证据 | 装 Go 1.24.11 后补跑：build / vet / 全量单测均通过，见下方「执行记录」 |

**顺带修掉的第二个真 bug**：`.github/workflows/ci.yml` 两处 `go-version` 钉在 `"1.23"`，
而 `server/go.mod` 要求 `go 1.24.0` + `toolchain go1.24.11`。靠 `GOTOOLCHAIN=auto`
隐式下载 1.24.11 才没红 —— **CI 声明的版本和实际跑的版本不是同一个**，
一旦 runner 禁网或设 `GOTOOLCHAIN=local` 就会硬失败。已对齐到 `"1.24"`。

**顺带修掉的一个真 bug**：原 `check-design.mjs` 的 `[A]` 段用全局 `failures === 0`
判断是否打印「15 个标签令牌齐全」。只要 `[B]` 段有任何失败，即便令牌完全正常也不会打印，
误导排查方向。已改为对比进入 `[A]` 前的快照 `before`。

---

## 待你拍板（P0，未处置）

### H2 — `deploy/.env` 曾被 git 追踪，凭据可能仍在历史 — **已复检：无真凭据**

- `master` 上该文件曾被追踪，已于提交 `2cbc6bd` 删除
- 历史中两份 `.env`（`31396b4` MVP init / `dd4cdc0` 举报功能）键名一致：11 个键（`APP_PORT/APP_MODE/DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME/REDIS_HOST/REDIS_PORT/REDIS_PASSWORD/JWT_SECRET`）
- **2026-07-28 复检结果**（commit 期间）：
  - 全部键值经结构特征判定为 placeholder / 弱占位符 / 空字符串
  - 关键确认：`DB_PASSWORD`（len=10）为弱占位符（`password`/`postgres` 类）
  - `JWT_SECRET`（len=39）为占位符风格字符串，非随机串
  - 无 MinIO / SMTP 凭据暴露（MVP 期尚未引入）
  - 未发现可解密为真实密钥的内容
- **结论**：无需轮换、无需历史重写
- **已加固**（同提交）：
  - `.gitignore` 显式忽略 `deploy/.env` / `deploy/.env.*` / `.env` / `.env.*.local`
  - `deploy/.env.example` 补齐 11 个键（+ MinIO/Meili/SMTP 注释）
  - 本条目状态改为「已复检无真凭据」

---

## 仅记录（非 bug，不阻塞合入）

### 文档与 spec 状态滞后

| 编号 | 位置 | 内容 |
|------|------|------|
| M1 | `proposal.md:3` | 状态仍为「待评审」，但 tasks.md 已写入评审决策并全部实施 |
| M2 | `proposal.md:63` | 成功标准 7 条全为 `- [ ]` 未勾选；「全站不再引用 blue-/gray-」实测 880 处 / 40 文件残留（已由基线锁定，只减不增）；tasks 未做区写「16 页」，实际 40 个文件 |
| M4 | `proposal.md:68` | 成功标准「对比度满足 WCAG AA」无任何验证证据，check-design 不含对比度计算 |
| L1 | `tasks.md:39` | 5.2「等待推送（代理恢复后）」未勾选，但分支已同步 origin |

### 测试覆盖缺口

| 编号 | 位置 | 内容 |
|------|------|------|
| M5 | `admin/settings/page.tsx`(267行) / `Icon.tsx`(140行) | 新增前端代码零测试；Icon 是主链路全量依赖 |
| M6 | `post/service_extra_test.go:187` | `types := []string{"emotion","insight","case","action"}` 未加 `achievement`，新类型 List 过滤无覆盖 |
| M7 | `internal/setting/` | 13 个测试全在 service/repo 层，handler 三个端点与 `TestSend` 真实投递路径零覆盖 |

### 架构与可维护性

| 编号 | 位置 | 内容 |
|------|------|------|
| M8 | `migrate.sh:45` / `apply-migrations.sh:25` / `setup_test.go:sortMigrations` | 「编号 → create_ 优先 → 文件名」排序规则**三处重复实现**，可各自漂移。本分支立意正是消除漂移 |
| M9 | CI vs 集成测试 | CI 走 `psql -f`，集成测试走自研 `splitSQL`（裸分号切分），两套解析器可各自通过/失败 |
| L4 | `constants.ts:80` | 心情「焦虑」复用 `achievement` 令牌，日后调成果展示配色会连带改变心情色 |

### 性能（均非当前瓶颈）

| 编号 | 位置 | 内容 |
|------|------|------|
| L2 | `setting/service.go:60-71` | `GetSMTPWithPassword` 调 `GetAll()` 两次，每次发信 2 次全表查询 |
| L3 | `middleware/email_verified.go:19` | 每个受保护请求一次 `system_settings` 全表 SELECT，无缓存 |

### 潜在缺陷（当前不触发，属加固项）

| 编号 | 位置 | 内容 |
|------|------|------|
| L5 | `check-schema-drift.py:42` | `KNOWN_SKIP={"schema_migrations"}` 作用于 models 字典，该表无 Go 模型，永不命中，是死代码 |
| L6 | `check-schema-drift.py:45` | `to_snake` 与 GORM NamingStrategy 不等价：`SMTPHost` → `smtphost` vs GORM 的 `smtp_host`。当前无此类字段，属潜在误报 |
| L7 | `setup_test.go:splitSQL` | 按裸分号切分，字符串字面量内含分号会截断语句。**已实测当前 23 个迁移 0 处**，属脆弱性非现存缺陷 |

### 观察备注

- `IsEmailVerifyRequired` 读库失败时 fail-open 放行（`service.go:136`），注释说明是为避免设置表异常导致全站无法发帖。取舍合理，但 DB 抖动会静默关闭邮箱验证闸门，需运行时观测才能判断影响。
- `migrate.sh:60,69` 把文件名插值进 SQL 字面量。文件名由仓库控制，实际不可注入，仅记录。

---

## 审查执行记录

实际跑过的命令与统计（非「应该能过」）：

```
node node_modules/vitest/vitest.mjs run     → Test Files 11 passed (11) / Tests 105 passed (105) / 34.87s
node node_modules/next/dist/bin/next build  → exit 0
node scripts/check-design.mjs --css         → [B] 19/59 零残留 + 基线内  [A] 令牌 6 项 + 15 标签令牌  全绿
探针反向验证（新增含 blue-600 的 tsx）      → exit 1，正确报红 2 项
```

后端（Go 1.24.11，与 go.mod 的 toolchain 声明一致）：

```
go build -o <tmp> ./cmd/...            → exit 0
go vet ./...                            → exit 0
go test ./internal/setting/... -v       → 13/13 PASS（含「密码永不回显」「SMTP 关闭时强制取消邮箱验证」）
go test ./internal/... -count=1         → 32 个模块全 ok，6 个 [no test files]，exit 0
```

### ⚠️ 环境陷阱：Windows 上 go test 必须设 GOTMPDIR

`WORKFLOW.md` 记录的解决方案是「逐模块运行 `go test`」，那是在治症状。

同一条 `go test ./internal/... -count=1`，只改 `GOTMPDIR` 一个变量：

| GOTMPDIR | 结果 |
|----------|------|
| 默认（`%TEMP%`） | **17 个 `[build failed]`**，报 `$WORK\bNNN\vet.cfg: Access is denied` / `importcfg: 系统找不到指定的路径` |
| `<repo>/tmp/gobuild`（D 盘） | **全绿，exit 0** |

根因是 `%TEMP%` 下 Go 的 `go-build*` 临时目录被外部进程中途删除或拒绝访问
（同一现象在本次会话中还连续吞掉了 4 次命令输出文件）。具体是哪个进程未取证。

**结论：本机跑后端测试请先设**

```bash
export GOTMPDIR="<local-project-path>"
```

比逐模块运行更快也更彻底。`WORKFLOW.md` 的环境备注建议在阶段 3.2 一并订正。
