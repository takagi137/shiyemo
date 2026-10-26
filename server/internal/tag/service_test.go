package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// GetHotTags 测试
// ============================================================

func TestTagService_GetHotTags(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetHotTags([]map[string]interface{}{
			{"id": 1, "name": "Go", "count": 10},
			{"id": 2, "name": "Python", "count": 8},
			{"id": 3, "name": "Java", "count": 5},
		})
		service := NewTagServiceWithDeps(repo)

		tags, err := service.GetHotTags(10)

		require.NoError(t, err)
		assert.Len(t, tags, 3)
	})

	t.Run("limit", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetHotTags([]map[string]interface{}{
			{"id": 1, "name": "Go", "count": 10},
			{"id": 2, "name": "Python", "count": 8},
			{"id": 3, "name": "Java", "count": 5},
		})
		service := NewTagServiceWithDeps(repo)

		tags, err := service.GetHotTags(2)

		require.NoError(t, err)
		assert.Len(t, tags, 2)
	})

	t.Run("default_limit", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetHotTags([]map[string]interface{}{
			{"id": 1, "name": "Go", "count": 10},
		})
		service := NewTagServiceWithDeps(repo)

		// limit = 0 应该使用默认值 20
		tags, err := service.GetHotTags(0)

		require.NoError(t, err)
		assert.Len(t, tags, 1)
	})

	t.Run("db_error", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetGetError(assert.AnError)
		service := NewTagServiceWithDeps(repo)

		_, err := service.GetHotTags(10)

		require.Error(t, err)
	})
}

// ============================================================
// UpdatePostTags 测试
// ============================================================

func TestTagService_UpdatePostTags(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewTagServiceWithDeps(repo)

		names, err := service.UpdatePostTags(1, []string{"Go", "Python"})

		require.NoError(t, err)
		assert.Len(t, names, 2)
		assert.Contains(t, names, "Go")
		assert.Contains(t, names, "Python")
	})

	t.Run("new_tags", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewTagServiceWithDeps(repo)

		// 第一次更新
		names1, err := service.UpdatePostTags(1, []string{"Go"})
		require.NoError(t, err)
		assert.Len(t, names1, 1)

		// 第二次更新（应该复用已有标签）
		names2, err := service.UpdatePostTags(1, []string{"Go", "Python"})
		require.NoError(t, err)
		assert.Len(t, names2, 2)
	})

	t.Run("get_or_create_error", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetGetError(assert.AnError)
		service := NewTagServiceWithDeps(repo)

		_, err := service.UpdatePostTags(1, []string{"Go"})

		require.Error(t, err)
	})

	t.Run("update_error", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetUpdateError(assert.AnError)
		service := NewTagServiceWithDeps(repo)

		_, err := service.UpdatePostTags(1, []string{"Go"})

		require.Error(t, err)
	})
}

// ============================================================
// GetPostTags 测试
// ============================================================

func TestTagService_GetPostTags(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewTagServiceWithDeps(repo)

		// 先添加标签
		service.UpdatePostTags(1, []string{"Go", "Python"})

		// 获取标签
		names, err := service.GetPostTags(1)

		require.NoError(t, err)
		assert.Len(t, names, 2)
		assert.Contains(t, names, "Go")
		assert.Contains(t, names, "Python")
	})

	t.Run("empty", func(t *testing.T) {
		repo := NewMockTagRepo()
		service := NewTagServiceWithDeps(repo)

		names, err := service.GetPostTags(999)

		require.NoError(t, err)
		assert.Len(t, names, 0)
	})

	t.Run("db_error", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetGetError(assert.AnError)
		service := NewTagServiceWithDeps(repo)

		_, err := service.GetPostTags(1)

		require.Error(t, err)
	})
}
