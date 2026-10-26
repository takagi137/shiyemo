package migrationx

import (
	"reflect"
	"testing"
)

func TestRankFor(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{"001_create_users.sql", 0},
		{"008_create_notifications.sql", 0},
		{"008_add_notification_type.sql", 1},
		{"016_create_daily_stats.sql", 0},
		{"020_create_journals.sql", 0},
		{"024_alter_something.sql", 1},
		{"no_prefix.sql", 1}, // 无 _create_ 后缀 → 落到普通
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := RankFor(c.name)
			if got != c.want {
				t.Errorf("RankFor(%q) = %d, want %d", c.name, got, c.want)
			}
		})
	}
}

func TestNumFor(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{"001_create_users.sql", "001"},
		{"008_create_notifications.sql", "008"},
		{"008_add_notification_type.sql", "008"},
		// 原 Go 行为：在第一个 '_' 截断。包含 '_' 的非编号文件会被当作「编号 = 前缀」
		{"no_prefix.sql", "no"},
		{"", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := NumFor(c.name)
			if got != c.want {
				t.Errorf("NumFor(%q) = %q, want %q", c.name, got, c.want)
			}
		})
	}
}

// TestSortInPlace 验证排序规则，避免明天有人改 Go 实现时改坏。
// 关键回归：008_add_notification_type 必须排在 008_create_notifications 之后。
//
// 注：原 Go 实现的字典序比较中，'n'（0x6E）> '0'（0x30），
// 因此 no_prefix.sql 实际排在所有数字编号之后。这是历史行为，
// 本测试用来「锁住」这个行为，防止回归。若要改变需同步改 shell 脚本。
func TestSortInPlace(t *testing.T) {
	input := []string{
		"008_add_notification_type.sql",
		"008_create_notifications.sql",
		"007_create_users.sql",
		"005_create_posts.sql",
		"005_add_post_index.sql", // 同编号下 add > create
		"001_create_users.sql",
		"no_prefix.sql", // 'n' (0x6E) > '0' (0x30)，实际排在数字编号之后
	}
	want := []string{
		"001_create_users.sql",
		"005_create_posts.sql",
		"005_add_post_index.sql",
		"007_create_users.sql",
		"008_create_notifications.sql", // 必须排在 add 之前
		"008_add_notification_type.sql",
		"no_prefix.sql",
	}
	got := append([]string(nil), input...)
	SortInPlace(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("排序错误\n got: %v\nwant: %v", got, want)
	}
}

// TestSortInPlace_Idempotent 重复排序结果稳定。
func TestSortInPlace_Idempotent(t *testing.T) {
	input := []string{"008_add_x.sql", "008_create_x.sql", "001_create_y.sql"}
	first := append([]string(nil), input...)
	SortInPlace(first)
	second := append([]string(nil), first...)
	SortInPlace(second)
	if !reflect.DeepEqual(first, second) {
		t.Errorf("排序不稳定\nfirst: %v\nsecond: %v", first, second)
	}
}

// TestSortInPlace_Empty 边界：空切片不崩。
func TestSortInPlace_Empty(t *testing.T) {
	SortInPlace(nil)
	SortInPlace([]string{})
}