// Package migrationx 提供跨语言共享的迁移排序规则。
//
// 单一信息源：脚本侧（scripts/migrate.sh / scripts/apply-migrations.sh）
// 和测试侧（server/tests/integration/setup_test.go）都通过此包实现。
//
// 规则：编号 → create_ 优先 → 文件名。
//   类型秩 0 = create_（建表），1 = 其它（add_/alter_ 等改表）。
//   纯字典序会出错：008_add_notification_type 排在 008_create_notifications
//   之前，但前者要 ALTER 后者建的表 —— 依赖倒置。
package migrationx

import (
	"sort"
	"strings"
)

// RankFor 给出「create_ 优先」的排序权重。
//   0 = create_  → 先建表再改表
//   1 = 其它（add_/alter_ 等）→ 后跑
func RankFor(name string) int {
	if strings.Contains(name, "_create_") {
		return 0
	}
	return 1
}

// NumFor 从文件名提取编号前缀。
//   "008_create_notifications.sql" → "008"
//   "no_prefix.sql"                → "no_prefix"（保留原名）
func NumFor(name string) string {
	if i := strings.IndexByte(name, '_'); i > 0 {
		return name[:i]
	}
	return name
}

// SortInPlace 按「编号 → 类型秩 → 文件名」原地排序迁移文件名切片。
// 与 scripts/_lib_migrations.sh 的 sorted_migrations 行为必须保持一致。
func SortInPlace(files []string) {
	sort.Slice(files, func(i, j int) bool {
		ni, nj := NumFor(files[i]), NumFor(files[j])
		if ni != nj {
			return ni < nj
		}
		ri, rj := RankFor(files[i]), RankFor(files[j])
		if ri != rj {
			return ri < rj
		}
		return files[i] < files[j]
	})
}