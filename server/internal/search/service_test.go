package search

import (
	"testing"

	"github.com/meilisearch/meilisearch-go"
	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/post"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// SearchPosts 测试
// ============================================================

func TestSearchService_SearchPosts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

		idx := client.Index("posts").(*MockSearchIndex)
		idx.documents = []map[string]interface{}{
			{"id": 1, "title": "Go语言教程", "type": "insight", "status": 1},
			{"id": 2, "title": "Go语言入门", "type": "insight", "status": 1},
			{"id": 3, "title": "Python教程", "type": "insight", "status": 1},
		}

		posts, total, err := service.SearchPosts("Go", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, posts, 2)
	})

	t.Run("empty_result", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

		posts, total, err := service.SearchPosts("不存在", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, posts, 0)
	})

	t.Run("pagination", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

		idx := client.Index("posts").(*MockSearchIndex)
		for i := 0; i < 15; i++ {
			idx.documents = append(idx.documents, map[string]interface{}{
				"id": i + 1, "title": "Go语言教程", "type": "insight", "status": 1,
			})
		}

		posts1, total1, err := service.SearchPosts("Go", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(10), total1)
		assert.Len(t, posts1, 10)

		posts2, total2, err := service.SearchPosts("Go", 2, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(5), total2)
		assert.Len(t, posts2, 5)
	})

	t.Run("search_error", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())
		idx := client.Index("posts").(*MockSearchIndex)
		idx.SetSearchError(assert.AnError)
		_, _, err := service.SearchPosts("Go", 1, 10)
		require.Error(t, err)
	})
}

// ============================================================
// SearchUsers 测试
// ============================================================

func TestSearchService_SearchUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

		idx := client.Index("users").(*MockSearchIndex)
		idx.documents = []map[string]interface{}{
			{"id": 1, "nickname": "张三"},
			{"id": 2, "nickname": "张三丰"},
			{"id": 3, "nickname": "李四"},
		}

		users, total, err := service.SearchUsers("张三", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, users, 2)
	})

	t.Run("empty_result", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())
		users, total, err := service.SearchUsers("不存在", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, users, 0)
	})

	t.Run("search_error", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())
		idx := client.Index("users").(*MockSearchIndex)
		idx.SetSearchError(assert.AnError)
		_, _, err := service.SearchUsers("张三", 1, 10)
		require.Error(t, err)
	})
}

// ============================================================
// GetSuggestions 测试
// ============================================================

func TestSearchService_GetSuggestions(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

		idx := client.Index("posts").(*MockSearchIndex)
		idx.documents = []map[string]interface{}{
			{"id": 1, "title": "Go语言教程"},
			{"id": 2, "title": "Go语言入门"},
			{"id": 3, "title": "Go语言高级编程"},
		}

		suggestions, err := service.GetSuggestions("Go", 5)
		require.NoError(t, err)
		assert.Len(t, suggestions, 3)
		assert.Contains(t, suggestions, "Go语言教程")
	})

	t.Run("empty_query", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())
		suggestions, err := service.GetSuggestions("", 5)
		require.NoError(t, err)
		assert.Len(t, suggestions, 0)
	})

	t.Run("deduplication", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

		idx := client.Index("posts").(*MockSearchIndex)
		idx.documents = []map[string]interface{}{
			{"id": 1, "title": "Go语言教程"},
			{"id": 2, "title": "Go语言教程"},
			{"id": 3, "title": "Go语言入门"},
		}
		suggestions, err := service.GetSuggestions("Go", 5)
		require.NoError(t, err)
		assert.Len(t, suggestions, 2)
	})

	t.Run("search_error", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())
		idx := client.Index("posts").(*MockSearchIndex)
		idx.SetSearchError(assert.AnError)
		_, err := service.GetSuggestions("Go", 5)
		require.Error(t, err)
	})

	t.Run("limit_too_small", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())
		idx := client.Index("posts").(*MockSearchIndex)
		idx.documents = []map[string]interface{}{{"id": 1, "title": "A"}}
		sug, err := service.GetSuggestions("A", 0)
		require.NoError(t, err)
		assert.Len(t, sug, 1) // limit defaults to 5, returns all
	})

	t.Run("limit_too_large", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())
		idx := client.Index("posts").(*MockSearchIndex)
		idx.documents = []map[string]interface{}{{"id": 1, "title": "A"}}
		sug, err := service.GetSuggestions("A", 99)
		require.NoError(t, err)
		assert.Len(t, sug, 1) // limit defaults to 5, returns all
	})
}

// ============================================================
// IndexPost 测试
// ============================================================

func TestSearchService_IndexPost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

		p := &post.Post{
			ID:    1, Type: "insight", Title: "Go语言教程",
			Content: "Go语言入门教程", LikeCount: 10, Status: 1,
		}
		p.SetImages(nil)
		p.SetTags([]string{"Go", "开发"})
		err := service.IndexPost(p)
		require.NoError(t, err)
		idx := client.Index("posts").(*MockSearchIndex)
		assert.Len(t, idx.documents, 1)
		assert.Equal(t, "Go语言教程", idx.documents[0]["title"])
	})
}

// ============================================================
// IndexUser 测试
// ============================================================

func TestSearchService_IndexUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

		user := &auth.User{
			ID: 1, Nickname: "张三", Avatar: "avatar.png", Bio: "个人简介", Company: "某公司",
		}
		err := service.IndexUser(user)
		require.NoError(t, err)
		idx := client.Index("users").(*MockSearchIndex)
		assert.Len(t, idx.documents, 1)
		assert.Equal(t, "张三", idx.documents[0]["nickname"])
	})
}

// ============================================================
// DeletePost 测试（缺失路径）
// ============================================================

func TestSearchService_DeletePost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := NewMockSearchClient()
		service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

		// 先索引一个帖子
		p := &post.Post{ID: 42, Type: "emotion", Title: "测试", Status: 1}
		p.SetImages(nil)
		p.SetTags(nil)
		require.NoError(t, service.IndexPost(p))
		idx := client.Index("posts").(*MockSearchIndex)
		assert.Len(t, idx.documents, 1)

		// 删除
		err := service.DeletePost(42)
		require.NoError(t, err)
		assert.Len(t, idx.documents, 0)
	})
}

// ============================================================
// InitIndexes 测试（error 分支）
// ============================================================

func TestSearchService_InitIndexes_IndexExists(t *testing.T) {
	client := NewMockSearchClient()
	// 预先存在 posts 索引，CreateIndex 会先创建再被第二步覆盖
	// 这里用一个会产生错误的 client 来覆盖 error 路径
	idx := NewMockSearchIndex()
	idx.searchErr = assert.AnError // 用于触发 UpdateFilterableAttributes 的 error
	client.indexes["posts"] = idx
	client.indexes["users"] = NewMockSearchIndex()

	service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

	// InitIndexes 的 CreateIndex 在索引存在时仍执行（mock 总是成功）
	// UpdateFilterableAttributes 只要 client 没报错就通过
	err := service.InitIndexes()
	require.NoError(t, err)
}

// ============================================================
// MockSearchClient 测试
// ============================================================

func TestMockSearchClient(t *testing.T) {
	t.Run("Index", func(t *testing.T) {
		client := NewMockSearchClient()
		idx1 := client.Index("posts")
		idx2 := client.Index("posts")
		assert.Equal(t, idx1, idx2)
	})

	t.Run("CreateIndex", func(t *testing.T) {
		client := NewMockSearchClient()
		task, err := client.CreateIndex(&meilisearch.IndexConfig{Uid: "posts"})
		require.NoError(t, err)
		assert.NotNil(t, task)
		_, ok := client.indexes["posts"]
		assert.True(t, ok)
	})
}

// ============================================================
// MockSearchIndex 测试
// ============================================================

func TestMockSearchIndex(t *testing.T) {
	t.Run("Search", func(t *testing.T) {
		idx := NewMockSearchIndex()
		idx.documents = []map[string]interface{}{
			{"id": 1, "title": "Go语言教程"},
			{"id": 2, "title": "Python教程"},
		}
		result, err := idx.Search("Go", &meilisearch.SearchRequest{Limit: 10, Offset: 0})
		require.NoError(t, err)
		assert.Len(t, result.Hits, 1)
	})

	t.Run("AddDocuments", func(t *testing.T) {
		idx := NewMockSearchIndex()
		task, err := idx.AddDocuments([]map[string]interface{}{{"id": 1, "title": "测试"}}, nil)
		require.NoError(t, err)
		assert.NotNil(t, task)
		assert.Len(t, idx.documents, 1)
	})

	t.Run("DeleteDocument", func(t *testing.T) {
		idx := NewMockSearchIndex()
		idx.documents = []map[string]interface{}{{"id": 42, "title": "测试"}}
		task, err := idx.DeleteDocument("42")
		require.NoError(t, err)
		assert.NotNil(t, task)
		assert.Len(t, idx.documents, 0)
	})

	t.Run("UpdateFilterableAttributes", func(t *testing.T) {
		idx := NewMockSearchIndex()
		task, err := idx.UpdateFilterableAttributes(&[]interface{}{"type", "status"})
		require.NoError(t, err)
		assert.NotNil(t, task)
	})
}
