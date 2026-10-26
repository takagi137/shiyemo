package report

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Create 测试
// ============================================================

func TestReportService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockReportRepo()
		service := NewReportServiceWithDeps(repo)

		req := CreateRequest{
			TargetType: "post",
			TargetID:   1,
			Reason:     "违规内容",
		}

		err := service.Create(1, req)

		require.NoError(t, err)
	})

	t.Run("already_exists", func(t *testing.T) {
		repo := NewMockReportRepo()
		repo.SetExists(true)
		service := NewReportServiceWithDeps(repo)

		req := CreateRequest{
			TargetType: "post",
			TargetID:   1,
			Reason:     "违规内容",
		}

		err := service.Create(1, req)

		require.Error(t, err)
	})

	t.Run("check_exists_error", func(t *testing.T) {
		repo := NewMockReportRepo()
		repo.SetError(assert.AnError)
		service := NewReportServiceWithDeps(repo)

		req := CreateRequest{
			TargetType: "post",
			TargetID:   1,
			Reason:     "违规内容",
		}

		err := service.Create(1, req)

		require.Error(t, err)
	})

	t.Run("create_error", func(t *testing.T) {
		repo := NewMockReportRepo()
		repo.SetError(assert.AnError)
		service := NewReportServiceWithDeps(repo)

		req := CreateRequest{
			TargetType: "post",
			TargetID:   1,
			Reason:     "违规内容",
		}

		err := service.Create(1, req)

		require.Error(t, err)
	})
}

// ============================================================
// List 测试
// ============================================================

func TestReportService_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockReportRepo()
		service := NewReportServiceWithDeps(repo)

		// 先创建几个举报
		service.Create(1, CreateRequest{TargetType: "post", TargetID: 1, Reason: "原因1"})
		service.Create(2, CreateRequest{TargetType: "post", TargetID: 2, Reason: "原因2"})
		service.Create(3, CreateRequest{TargetType: "comment", TargetID: 3, Reason: "原因3"})

		reports, total, err := service.List(1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, reports, 3)
	})

	t.Run("pagination", func(t *testing.T) {
		repo := NewMockReportRepo()
		service := NewReportServiceWithDeps(repo)

		// 创建5个举报
		for i := 1; i <= 5; i++ {
			service.Create(1, CreateRequest{TargetType: "post", TargetID: int64(i), Reason: "原因"})
		}

		// 第一页，每页2条
		reports, total, err := service.List(1, 2)

		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, reports, 2)

		// 第二页
		reports, total, err = service.List(2, 2)

		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, reports, 2)

		// 第三页（只剩1条）
		reports, total, err = service.List(3, 2)

		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, reports, 1)
	})

	t.Run("empty", func(t *testing.T) {
		repo := NewMockReportRepo()
		service := NewReportServiceWithDeps(repo)

		reports, total, err := service.List(1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, reports, 0)
	})

	t.Run("db_error", func(t *testing.T) {
		repo := NewMockReportRepo()
		repo.SetError(assert.AnError)
		service := NewReportServiceWithDeps(repo)

		_, _, err := service.List(1, 10)

		require.Error(t, err)
	})
}

// ============================================================
// Handle 测试
// ============================================================

func TestReportService_Handle(t *testing.T) {
	t.Run("success_approve", func(t *testing.T) {
		repo := NewMockReportRepo()
		service := NewReportServiceWithDeps(repo)

		// 先创建举报
		service.Create(1, CreateRequest{TargetType: "post", TargetID: 1, Reason: "违规内容"})

		// 处理举报（状态2=已处理）
		err := service.Handle(1, 2)

		require.NoError(t, err)
	})

	t.Run("success_reject", func(t *testing.T) {
		repo := NewMockReportRepo()
		service := NewReportServiceWithDeps(repo)

		// 先创建举报
		service.Create(1, CreateRequest{TargetType: "post", TargetID: 1, Reason: "违规内容"})

		// 处理举报（状态3=已驳回）
		err := service.Handle(1, 3)

		require.NoError(t, err)
	})

	t.Run("not_found", func(t *testing.T) {
		repo := NewMockReportRepo()
		service := NewReportServiceWithDeps(repo)

		err := service.Handle(999, 2)

		require.Error(t, err)
	})

	t.Run("db_error", func(t *testing.T) {
		repo := NewMockReportRepo()
		repo.SetError(assert.AnError)
		service := NewReportServiceWithDeps(repo)

		err := service.Handle(1, 2)

		require.Error(t, err)
	})
}
