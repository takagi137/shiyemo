package search

import (
	"github.com/meilisearch/meilisearch-go"
)

// SearchIndex 搜索索引接口（封装 meilisearch.IndexManager）
type SearchIndex interface {
	Search(query string, request *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error)
	AddDocuments(documentsPtr interface{}, opts *meilisearch.DocumentOptions) (*meilisearch.TaskInfo, error)
	DeleteDocument(documentID string) (*meilisearch.TaskInfo, error)
	UpdateFilterableAttributes(attributes *[]interface{}) (*meilisearch.TaskInfo, error)
}

// SearchClient 搜索客户端接口（封装 meilisearch.ServiceManager）
type SearchClient interface {
	Index(uid string) SearchIndex
	CreateIndex(config *meilisearch.IndexConfig) (*meilisearch.TaskInfo, error)
}

// MeiliClient Meilisearch 客户端实现
type MeiliClient struct {
	client meilisearch.ServiceManager
}

// NewMeiliClient 创建 Meilisearch 客户端
func NewMeiliClient(client meilisearch.ServiceManager) *MeiliClient {
	return &MeiliClient{client: client}
}

// Index 获取索引
func (c *MeiliClient) Index(uid string) SearchIndex {
	return &MeiliIndex{index: c.client.Index(uid)}
}

// CreateIndex 创建索引
func (c *MeiliClient) CreateIndex(config *meilisearch.IndexConfig) (*meilisearch.TaskInfo, error) {
	return c.client.CreateIndex(config)
}

// MeiliIndex Meilisearch 索引实现
type MeiliIndex struct {
	index meilisearch.IndexManager
}

// Search 搜索
func (i *MeiliIndex) Search(query string, request *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error) {
	return i.index.Search(query, request)
}

// AddDocuments 添加文档
func (i *MeiliIndex) AddDocuments(documentsPtr interface{}, opts *meilisearch.DocumentOptions) (*meilisearch.TaskInfo, error) {
	return i.index.AddDocuments(documentsPtr, opts)
}

// DeleteDocument 删除文档
func (i *MeiliIndex) DeleteDocument(documentID string) (*meilisearch.TaskInfo, error) {
	return i.index.DeleteDocument(documentID, nil)
}

// UpdateFilterableAttributes 更新可过滤属性
func (i *MeiliIndex) UpdateFilterableAttributes(attributes *[]interface{}) (*meilisearch.TaskInfo, error) {
	return i.index.UpdateFilterableAttributes(attributes)
}
