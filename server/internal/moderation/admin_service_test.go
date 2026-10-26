package moderation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// AddWord 测试
// ============================================================

func TestService_AddWord(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockModerationRepo()
		service := NewServiceWithDeps(repo)

		word, err := service.AddWord("赌博")

		require.NoError(t, err)
		assert.NotNil(t, word)
		assert.Equal(t, "赌博", word.Word)
	})

	t.Run("already_exists", func(t *testing.T) {
		repo := NewMockModerationRepo()
		service := NewServiceWithDeps(repo)

		// 第一次添加
		word1, _ := service.AddWord("赌博")
		// 第二次添加同名词（复用）
		word2, err := service.AddWord("赌博")

		require.NoError(t, err)
		assert.Equal(t, word1.ID, word2.ID)
	})

	t.Run("get_error", func(t *testing.T) {
		repo := NewMockModerationRepo()
		repo.SetGetError(assert.AnError)
		service := NewServiceWithDeps(repo)

		_, err := service.AddWord("赌博")

		require.Error(t, err)
	})

	t.Run("create_error", func(t *testing.T) {
		repo := NewMockModerationRepo()
		repo.SetCreateError(assert.AnError)
		service := NewServiceWithDeps(repo)

		_, err := service.AddWord("赌博")

		require.Error(t, err)
	})
}

// ============================================================
// DeleteWord 测试
// ============================================================

func TestService_DeleteWord(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockModerationRepo()
		service := NewServiceWithDeps(repo)

		// 先添加
		word, _ := service.AddWord("赌博")

		// 删除
		err := service.DeleteWord(word.ID)

		require.NoError(t, err)
	})

	t.Run("not_found", func(t *testing.T) {
		repo := NewMockModerationRepo()
		service := NewServiceWithDeps(repo)

		err := service.DeleteWord(999)

		require.Error(t, err)
	})

	t.Run("delete_error", func(t *testing.T) {
		repo := NewMockModerationRepo()
		repo.SetDeleteError(assert.AnError)
		service := NewServiceWithDeps(repo)

		err := service.DeleteWord(1)

		require.Error(t, err)
	})
}

// ============================================================
// ListWords 测试
// ============================================================

func TestService_ListWords(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewMockModerationRepo()
		service := NewServiceWithDeps(repo)

		// 添加几个敏感词
		service.AddWord("赌博")
		service.AddWord("色情")
		service.AddWord("暴力")

		words, total, err := service.ListWords(1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, words, 3)
	})

	t.Run("pagination", func(t *testing.T) {
		repo := NewMockModerationRepo()
		service := NewServiceWithDeps(repo)

		// 添加5个敏感词
		for i := 1; i <= 5; i++ {
			service.AddWord("word" + string(rune('0'+i)))
		}

		// 第一页，每页2条
		words, total, err := service.ListWords(1, 2)

		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, words, 2)
	})

	t.Run("default_values", func(t *testing.T) {
		repo := NewMockModerationRepo()
		service := NewServiceWithDeps(repo)

		// page=0 应该使用默认值1
		words, total, err := service.ListWords(0, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, words, 0)
	})

	t.Run("large_page_size", func(t *testing.T) {
		repo := NewMockModerationRepo()
		service := NewServiceWithDeps(repo)

		// pageSize>100 应该使用默认值50
		words, total, err := service.ListWords(1, 200)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, words, 0)
	})

	t.Run("empty", func(t *testing.T) {
		repo := NewMockModerationRepo()
		service := NewServiceWithDeps(repo)

		words, total, err := service.ListWords(1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, words, 0)
	})

	t.Run("db_error", func(t *testing.T) {
		repo := NewMockModerationRepo()
		repo.SetGetError(assert.AnError)
		service := NewServiceWithDeps(repo)

		_, _, err := service.ListWords(1, 10)

		require.Error(t, err)
	})
}
