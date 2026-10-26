package search

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/post"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// IndexPost - 错误路径覆盖（AddDocuments 失败）
// ============================================================

func TestSearchService_IndexPost_AddDocsError(t *testing.T) {
	client := NewMockSearchClient()
	service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

	idx := client.Index("posts").(*MockSearchIndex)
	idx.SetAddDocsError(assert.AnError)

	p := &post.Post{
		ID: 5, Type: "emotion", Title: "测试", Content: "内容", Status: 1,
	}
	p.SetImages(nil)
	p.SetTags(nil)

	err := service.IndexPost(p)
	require.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

// ============================================================
// IndexUser - 错误路径覆盖（AddDocuments 失败）
// ============================================================

func TestSearchService_IndexUser_AddDocsError(t *testing.T) {
	client := NewMockSearchClient()
	service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

	idx := client.Index("users").(*MockSearchIndex)
	idx.SetAddDocsError(assert.AnError)

	u := &auth.User{ID: 5, Nickname: "测试", Avatar: "a.png", Bio: "简介"}
	err := service.IndexUser(u)
	require.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

// ============================================================
// InitIndexes - CreateIndex 失败分支覆盖
// ============================================================

func TestSearchService_InitIndexes_CreateIndexError(t *testing.T) {
	client := NewMockSearchClient()
	client.createErr = assert.AnError // CreateIndex 返回错误
	service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

	// CreateIndex 失败 → 两个 log.Printf 分支被触发，但不阻断流程
	err := service.InitIndexes()
	require.NoError(t, err)
}

// ============================================================
// IndexPost - 覆盖带 tags 和无 indexer 路径
// ============================================================

func TestSearchService_IndexPost_WithTags(t *testing.T) {
	client := NewMockSearchClient()
	service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

	p := &post.Post{
		ID:         10,
		Type:       "insight",
		Title:      "带标签帖子",
		Content:    "内容",
		LikeCount:  5,
		CommentCount: 2,
		Status:     1,
	}
	p.SetImages([]string{"a.png"})
	p.SetTags([]string{"Go", "测试"})

	err := service.IndexPost(p)
	require.NoError(t, err)

	idx := client.Index("posts").(*MockSearchIndex)
	require.Len(t, idx.documents, 1)
	assert.Equal(t, "带标签帖子", idx.documents[0]["title"])
}

// ============================================================
// IndexUser - 覆盖正常路径
// ============================================================

func TestSearchService_IndexUser_NormalPath(t *testing.T) {
	client := NewMockSearchClient()
	service := MockSearchServiceWithDeps(client, NewMockSearchRepo())

	u := &auth.User{
		ID:       20,
		Nickname: "张三",
		Avatar:   "avatar.jpg",
		Bio:      "简介",
		Company:  "公司",
	}
	err := service.IndexUser(u)
	require.NoError(t, err)

	idx := client.Index("users").(*MockSearchIndex)
	require.Len(t, idx.documents, 1)
	assert.Equal(t, "张三", idx.documents[0]["nickname"])
}

// ============================================================
// InitIndexes - 预建索引后触发已有路径
// ============================================================

func TestSearchService_InitIndexes_IndexesAlreadyExist(t *testing.T) {
	client := NewMockSearchClient()
	// 预建索引
	client.indexes["posts"] = NewMockSearchIndex()
	client.indexes["users"] = NewMockSearchIndex()

	service := MockSearchServiceWithDeps(client, NewMockSearchRepo())
	err := service.InitIndexes()
	require.NoError(t, err)
}

// ============================================================
// SyncAllPosts 测试
// ============================================================

func TestSearchService_SyncAllPosts_Success(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	repo.activePosts = []post.Post{
		{ID: 1, Type: "insight", Title: "Go语言教程", Content: "入门", Status: 1, CreatedAt: time.Now()},
		{ID: 2, Type: "emotion", Title: "今天心情好", Content: "不错", Status: 1, CreatedAt: time.Now()},
	}
	repo.activePosts[0].SetTags([]string{"Go"})
	repo.activePosts[1].SetTags(nil)
	service := MockSearchServiceWithDeps(client, repo)

	err := service.SyncAllPosts()
	require.NoError(t, err)

	// 确认索引中添加了 2 条文档
	idx := client.Index("posts").(*MockSearchIndex)
	assert.Len(t, idx.documents, 2)
}

func TestSearchService_SyncAllPosts_Empty(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	service := MockSearchServiceWithDeps(client, repo)

	err := service.SyncAllPosts()
	require.NoError(t, err)

	idx := client.Index("posts").(*MockSearchIndex)
	assert.Len(t, idx.documents, 0)
}

func TestSearchService_SyncAllPosts_RepoError(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	repo.postsErr = errors.New("database error")
	service := MockSearchServiceWithDeps(client, repo)

	err := service.SyncAllPosts()
	require.Error(t, err)
	assert.Equal(t, err, errors.New("database error"))
}

func TestSearchService_SyncAllPosts_IndexError(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	p := &post.Post{ID: 10, Type: "emotion", Title: "测试", Status: 1, CreatedAt: time.Now()}
	p.SetTags(nil)
	repo.activePosts = []post.Post{{ID: 10, Type: "emotion", Title: "测试", Status: 1, CreatedAt: time.Now()}}
	repo.activePosts[0].SetTags(nil)

	// 配置 posts 索引的 AddDocuments 失败
	idx := client.Index("posts").(*MockSearchIndex)
	idx.SetAddDocsError(errors.New("index error"))

	service := MockSearchServiceWithDeps(client, repo)

	// SyncAllPosts 不会因为单条帖子失败而中止，仍返回 nil
	err := service.SyncAllPosts()
	require.NoError(t, err)
	_ = p
}

// ============================================================
// SyncAllUsers 测试
// ============================================================

func TestSearchService_SyncAllUsers_Success(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	repo.activeUsers = []auth.User{
		{ID: 1, Nickname: "张三", Avatar: "a.png", Bio: "简介", Company: "公司A"},
		{ID: 2, Nickname: "李四", Avatar: "b.png", Bio: "简介B", Company: "公司B"},
	}
	service := MockSearchServiceWithDeps(client, repo)

	err := service.SyncAllUsers()
	require.NoError(t, err)

	idx := client.Index("users").(*MockSearchIndex)
	assert.Len(t, idx.documents, 2)
}

func TestSearchService_SyncAllUsers_Empty(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	service := MockSearchServiceWithDeps(client, repo)

	err := service.SyncAllUsers()
	require.NoError(t, err)

	idx := client.Index("users").(*MockSearchIndex)
	assert.Len(t, idx.documents, 0)
}

func TestSearchService_SyncAllUsers_RepoError(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	repo.usersErr = errors.New("database error")
	service := MockSearchServiceWithDeps(client, repo)

	err := service.SyncAllUsers()
	require.Error(t, err)
}

func TestSearchService_SyncAllUsers_IndexError(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	repo.activeUsers = []auth.User{
		{ID: 5, Nickname: "测试", Avatar: "a.png", Bio: "简介"},
	}
	// 配置 users 索引的 AddDocuments 失败
	idx := client.Index("users").(*MockSearchIndex)
	idx.SetAddDocsError(errors.New("index error"))

	service := MockSearchServiceWithDeps(client, repo)
	// 单条失败不中止
	err := service.SyncAllUsers()
	require.NoError(t, err)
}

// ============================================================
// NewSearchServiceWithDeps 测试
// ============================================================

func TestSearchService_NewSearchServiceWithDeps(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	svc := NewSearchServiceWithDeps(client, repo)
	assert.NotNil(t, svc)
	assert.Equal(t, client, svc.client)
	assert.Equal(t, repo, svc.repo)
}

// ============================================================
// repo.go 测试
// ============================================================

func TestSearchRepo_Constructor(t *testing.T) {
	r := NewSearchRepo()
	assert.NotNil(t, r)
}

func TestSearchRepo_WithDB(t *testing.T) {
	// nil DB 不会 panic，用于测试构造验证
	r := NewSearchRepoWithDB(nil)
	assert.NotNil(t, r)
}

// ============================================================
// Handler 测试（handler.go）
// ============================================================

// setupSearchTestRouter 创建带搜索 handler 的测试路由
func setupSearchTestRouter(handler *SearchHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/search/posts", handler.SearchPosts)
	r.GET("/search/users", handler.SearchUsers)
	r.POST("/search/sync", handler.SyncIndexes)
	r.GET("/suggest", handler.GetSuggestions)
	return r
}

func TestSearchHandler_NewSearchHandler(t *testing.T) {
	h := NewSearchHandler()
	assert.NotNil(t, h)
}

func TestSearchHandler_SearchPosts_MissingQuery(t *testing.T) {
	handler := &SearchHandler{service: MockSearchServiceWithDeps(NewMockSearchClient(), NewMockSearchRepo())}
	r := setupSearchTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/search/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchHandler_SearchPosts_Success(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	handler := &SearchHandler{service: MockSearchServiceWithDeps(client, repo)}
	r := setupSearchTestRouter(handler)

	idx := client.Index("posts").(*MockSearchIndex)
	idx.documents = []map[string]interface{}{{"id": 1, "title": "Go教程", "status": 1}}

	req := httptest.NewRequest(http.MethodGet, "/search/posts?q=Go&page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchHandler_SearchPosts_SearchError(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	handler := &SearchHandler{service: MockSearchServiceWithDeps(client, repo)}
	r := setupSearchTestRouter(handler)

	idx := client.Index("posts").(*MockSearchIndex)
	idx.SetSearchError(errors.New("search failed"))

	req := httptest.NewRequest(http.MethodGet, "/search/posts?q=Go&page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSearchHandler_SearchPosts_DefaultPagination(t *testing.T) {
	handler := &SearchHandler{service: MockSearchServiceWithDeps(NewMockSearchClient(), NewMockSearchRepo())}
	r := setupSearchTestRouter(handler)

	// page/size 缺省，走 default 值
	req := httptest.NewRequest(http.MethodGet, "/search/posts?q=test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchHandler_SearchUsers_MissingQuery(t *testing.T) {
	handler := &SearchHandler{service: MockSearchServiceWithDeps(NewMockSearchClient(), NewMockSearchRepo())}
	r := setupSearchTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/search/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchHandler_SearchUsers_Success(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	handler := &SearchHandler{service: MockSearchServiceWithDeps(client, repo)}
	r := setupSearchTestRouter(handler)

	idx := client.Index("users").(*MockSearchIndex)
	idx.documents = []map[string]interface{}{{"id": 1, "nickname": "张三"}}

	req := httptest.NewRequest(http.MethodGet, "/search/users?q=张三", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchHandler_SearchUsers_SearchError(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	handler := &SearchHandler{service: MockSearchServiceWithDeps(client, repo)}
	r := setupSearchTestRouter(handler)

	idx := client.Index("users").(*MockSearchIndex)
	idx.SetSearchError(errors.New("search failed"))

	req := httptest.NewRequest(http.MethodGet, "/search/users?q=张三", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSearchHandler_SyncIndexes_Success(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	handler := &SearchHandler{service: MockSearchServiceWithDeps(client, repo)}
	r := setupSearchTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/search/sync", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchHandler_SyncIndexes_PostsError(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	repo.postsErr = errors.New("db error")
	handler := &SearchHandler{service: MockSearchServiceWithDeps(client, repo)}
	r := setupSearchTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/search/sync", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSearchHandler_SyncIndexes_UsersError(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	repo.usersErr = errors.New("db error")
	handler := &SearchHandler{service: MockSearchServiceWithDeps(client, repo)}
	r := setupSearchTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/search/sync", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSearchHandler_GetSuggestions_MissingQuery(t *testing.T) {
	handler := &SearchHandler{service: MockSearchServiceWithDeps(NewMockSearchClient(), NewMockSearchRepo())}
	r := setupSearchTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/suggest", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchHandler_GetSuggestions_Success(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	handler := &SearchHandler{service: MockSearchServiceWithDeps(client, repo)}
	r := setupSearchTestRouter(handler)

	idx := client.Index("posts").(*MockSearchIndex)
	idx.documents = []map[string]interface{}{
		{"id": 1, "title": "Go语言教程"},
		{"id": 2, "title": "Python教程"},
	}

	req := httptest.NewRequest(http.MethodGet, "/suggest?q=Go&limit=5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Go语言教程")
}

func TestSearchHandler_GetSuggestions_SearchError(t *testing.T) {
	client := NewMockSearchClient()
	repo := NewMockSearchRepo()
	handler := &SearchHandler{service: MockSearchServiceWithDeps(client, repo)}
	r := setupSearchTestRouter(handler)

	idx := client.Index("posts").(*MockSearchIndex)
	idx.SetSearchError(errors.New("search failed"))

	req := httptest.NewRequest(http.MethodGet, "/suggest?q=Go", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

