package search

import (
	"encoding/json"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/post"
)

// ============================================================
// Mock 实现
// ============================================================

// MockSearchClient mock 搜索客户端
type MockSearchClient struct {
	indexes       map[string]*MockSearchIndex
	createErr     error // CreateIndex 返回的错误
}

// NewMockSearchClient 创建 mock 搜索客户端
func NewMockSearchClient() *MockSearchClient {
	return &MockSearchClient{
		indexes: make(map[string]*MockSearchIndex),
	}
}

// Index 获取索引
func (c *MockSearchClient) Index(uid string) SearchIndex {
	if idx, ok := c.indexes[uid]; ok {
		return idx
	}
	idx := NewMockSearchIndex()
	c.indexes[uid] = idx
	return idx
}

// CreateIndex 创建索引（支持模拟错误）
func (c *MockSearchClient) CreateIndex(config *meilisearch.IndexConfig) (*meilisearch.TaskInfo, error) {
	if c.createErr != nil {
		return nil, c.createErr
	}
	if _, ok := c.indexes[config.Uid]; !ok {
		c.indexes[config.Uid] = NewMockSearchIndex()
	}
	return &meilisearch.TaskInfo{TaskUID: 1}, nil
}

// MockSearchIndex mock 搜索索引
type MockSearchIndex struct {
	documents    []map[string]interface{}
	searchErr    error
	addDocsErr   error // AddDocuments 返回的错误
}

// NewMockSearchIndex 创建 mock 搜索索引
func NewMockSearchIndex() *MockSearchIndex {
	return &MockSearchIndex{}
}

// SetSearchError 设置搜索错误
func (i *MockSearchIndex) SetSearchError(err error) {
	i.searchErr = err
}

// SetAddDocsError 设置 AddDocuments 错误
func (i *MockSearchIndex) SetAddDocsError(err error) {
	i.addDocsErr = err
}

// Search 搜索
func (i *MockSearchIndex) Search(query string, request *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error) {
	if i.searchErr != nil {
		return nil, i.searchErr
	}

	// 模拟搜索结果
	var hits []map[string]interface{}
	for _, doc := range i.documents {
		// 简单的关键词匹配
		if query == "" {
			hits = append(hits, doc)
			continue
		}
		if title, ok := doc["title"].(string); ok {
			if contains(title, query) {
				hits = append(hits, doc)
			}
		}
		if nickname, ok := doc["nickname"].(string); ok {
			if contains(nickname, query) {
				hits = append(hits, doc)
			}
		}
	}

	// 模拟分页
	offset := int(request.Offset)
	limit := int(request.Limit)
	if offset >= len(hits) {
		hits = []map[string]interface{}{}
	} else {
		end := offset + limit
		if end > len(hits) {
			end = len(hits)
		}
		hits = hits[offset:end]
	}

	// 转换为 meilisearch.Hits
	meiliHits := make(meilisearch.Hits, len(hits))
	for idx, hit := range hits {
		meiliHit := make(meilisearch.Hit, len(hit))
		for k, v := range hit {
			data, _ := json.Marshal(v)
			meiliHit[k] = data
		}
		meiliHits[idx] = meiliHit
	}

	return &meilisearch.SearchResponse{
		Hits:               meiliHits,
		EstimatedTotalHits: int64(len(hits)),
	}, nil
}

// AddDocuments 添加文档（支持错误注入）
func (i *MockSearchIndex) AddDocuments(documentsPtr interface{}, opts *meilisearch.DocumentOptions) (*meilisearch.TaskInfo, error) {
	if i.addDocsErr != nil {
		return nil, i.addDocsErr
	}
	docs, ok := documentsPtr.([]map[string]interface{})
	if !ok {
		return nil, nil
	}
	i.documents = append(i.documents, docs...)
	return &meilisearch.TaskInfo{TaskUID: 1}, nil
}

// DeleteDocument 删除文档
func (i *MockSearchIndex) DeleteDocument(documentID string) (*meilisearch.TaskInfo, error) {
	for idx, doc := range i.documents {
		if id, ok := doc["id"]; ok {
			if fmt.Sprintf("%v", id) == documentID {
				i.documents = append(i.documents[:idx], i.documents[idx+1:]...)
				break
			}
		}
	}
	return &meilisearch.TaskInfo{TaskUID: 1}, nil
}

// UpdateFilterableAttributes 更新可过滤属性
func (i *MockSearchIndex) UpdateFilterableAttributes(attributes *[]interface{}) (*meilisearch.TaskInfo, error) {
	return &meilisearch.TaskInfo{TaskUID: 1}, nil
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && len(substr) > 0 && containsIgnoreCase(s, substr))
}

func containsIgnoreCase(s, substr string) bool {
	s = toLowerCase(s)
	substr = toLowerCase(substr)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLowerCase(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

// MockSearchRepo mock 搜索仓储（实现 SearchRepoInterface）
type MockSearchRepo struct {
	activePosts []post.Post
	activeUsers []auth.User
	postsErr    error
	usersErr    error
}

// NewMockSearchRepo 创建 mock 搜索仓储
func NewMockSearchRepo() *MockSearchRepo {
	return &MockSearchRepo{}
}

func (r *MockSearchRepo) GetActivePosts() ([]post.Post, error) {
	if r.postsErr != nil {
		return nil, r.postsErr
	}
	return r.activePosts, nil
}

func (r *MockSearchRepo) GetActiveUsers() ([]auth.User, error) {
	if r.usersErr != nil {
		return nil, r.usersErr
	}
	return r.activeUsers, nil
}

// MockDB 保留向后兼容（旧测试中未用到 DB 逻辑）
type MockDB struct{}
func NewMockDB() *MockDB { return &MockDB{} }

// MockSearchServiceWithDeps 创建带 mock 依赖的搜索服务
func MockSearchServiceWithDeps(client SearchClient, repo SearchRepoInterface) *SearchService {
	return &SearchService{
		client: client,
		repo:   repo,
	}
}

// MockSearchResult 模拟搜索结果
func MockSearchResult(hits []map[string]interface{}, total int64) *meilisearch.SearchResponse {
	meiliHits := make(meilisearch.Hits, len(hits))
	for i, hit := range hits {
		meiliHit := make(meilisearch.Hit, len(hit))
		for k, v := range hit {
			data, _ := json.Marshal(v)
			meiliHit[k] = data
		}
		meiliHits[i] = meiliHit
	}
	return &meilisearch.SearchResponse{
		Hits:               meiliHits,
		EstimatedTotalHits: total,
	}
}
