package search

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/meilisearch/meilisearch-go"
	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/database"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/post"
)

// SearchService 搜索服务
type SearchService struct {
	client SearchClient
	repo   SearchRepoInterface
}

// NewSearchService 创建搜索服务
func NewSearchService() *SearchService {
	meiliURL := os.Getenv("MEILI_URL")
	if meiliURL == "" {
		meiliURL = "http://localhost:7700"
	}
	meiliKey := os.Getenv("MEILI_KEY")
	if meiliKey == "" {
		meiliKey = "shiyemo-meili-key"
	}
	client := meilisearch.New(meiliURL, meilisearch.WithAPIKey(meiliKey))
	return &SearchService{
		client: NewMeiliClient(client),
		repo:   NewSearchRepoWithDB(database.GetDB()),
	}
}

// InitIndexes 初始化索引
func (s *SearchService) InitIndexes() error {
	// 创建帖子索引
	_, err := s.client.CreateIndex(&meilisearch.IndexConfig{
		Uid: "posts",
	})
	if err != nil {
		// 索引可能已存在，忽略错误
		log.Printf("创建帖子索引: %v", err)
	}

	// 创建用户索引
	_, err = s.client.CreateIndex(&meilisearch.IndexConfig{
		Uid: "users",
	})
	if err != nil {
		log.Printf("创建用户索引: %v", err)
	}

	// 配置帖子索引
	postsIndex := s.client.Index("posts")
	_, err = postsIndex.UpdateFilterableAttributes(&[]interface{}{"type", "status"})
	if err != nil {
		log.Printf("配置帖子索引: %v", err)
	}

	return nil
}

// NewSearchServiceWithDeps 创建搜索服务（注入依赖，用于测试）
func NewSearchServiceWithDeps(client SearchClient, repo SearchRepoInterface) *SearchService {
	return &SearchService{
		client: client,
		repo:   repo,
	}
}

// IndexPost 索引帖子（实现 post.PostIndexer 接口）
func (s *SearchService) IndexPost(p *post.Post) error {
	postsIndex := s.client.Index("posts")

	doc := map[string]interface{}{
		"id":         p.ID,
		"type":       p.Type,
		"title":      p.Title,
		"content":    p.Content,
		"tags":       p.GetTags(),
		"like_count": p.LikeCount,
		"status":     p.Status,
		"created_at": p.CreatedAt.Unix(),
	}

	_, err := postsIndex.AddDocuments([]map[string]interface{}{doc}, nil)
	return err
}

// DeletePost 从索引中删除帖子（实现 post.PostIndexer 接口）
func (s *SearchService) DeletePost(postID int64) error {
	postsIndex := s.client.Index("posts")
	_, err := postsIndex.DeleteDocument(fmt.Sprintf("%d", postID))
	return err
}

// IndexUser 索引用户
func (s *SearchService) IndexUser(u *auth.User) error {
	usersIndex := s.client.Index("users")

	doc := map[string]interface{}{
		"id":       u.ID,
		"nickname": u.Nickname,
		"avatar":   u.Avatar,
		"bio":      u.Bio,
		"company":  u.Company,
	}

	_, err := usersIndex.AddDocuments([]map[string]interface{}{doc}, nil)
	return err
}

// SearchPosts 搜索帖子
func (s *SearchService) SearchPosts(query string, page, size int) ([]map[string]interface{}, int64, error) {
	postsIndex := s.client.Index("posts")

	searchRequest := &meilisearch.SearchRequest{
		Limit:  int64(size),
		Offset: int64((page - 1) * size),
		Filter: "status = 1",
	}

	result, err := postsIndex.Search(query, searchRequest)
	if err != nil {
		return nil, 0, bizerr.ErrInternal
	}

	var posts []map[string]interface{}
	for _, hit := range result.Hits {
		hitJSON, _ := json.Marshal(hit)
		var post map[string]interface{}
		json.Unmarshal(hitJSON, &post)
		posts = append(posts, post)
	}

	total := int64(result.EstimatedTotalHits)

	return posts, total, nil
}

// SearchUsers 搜索用户
func (s *SearchService) SearchUsers(query string, page, size int) ([]map[string]interface{}, int64, error) {
	usersIndex := s.client.Index("users")

	searchRequest := &meilisearch.SearchRequest{
		Limit:  int64(size),
		Offset: int64((page - 1) * size),
	}

	result, err := usersIndex.Search(query, searchRequest)
	if err != nil {
		return nil, 0, bizerr.ErrInternal
	}

	var users []map[string]interface{}
	for _, hit := range result.Hits {
		hitJSON, _ := json.Marshal(hit)
		var user map[string]interface{}
		json.Unmarshal(hitJSON, &user)
		users = append(users, user)
	}

	total := int64(result.EstimatedTotalHits)

	return users, total, nil
}

// SyncAllPosts 同步所有帖子到索引
func (s *SearchService) SyncAllPosts() error {
	posts, err := s.repo.GetActivePosts()
	if err != nil {
		return err
	}

	for i := range posts {
		if err := s.IndexPost(&posts[i]); err != nil {
			log.Printf("同步帖子 %d 失败: %v", posts[i].ID, err)
		}
	}

	fmt.Printf("同步了 %d 个帖子到索引\n", len(posts))
	return nil
}

// SyncAllUsers 同步所有用户到索引
func (s *SearchService) SyncAllUsers() error {
	users, err := s.repo.GetActiveUsers()
	if err != nil {
		return err
	}

	for i := range users {
		if err := s.IndexUser(&users[i]); err != nil {
			log.Printf("同步用户 %d 失败: %v", users[i].ID, err)
		}
	}

	fmt.Printf("同步了 %d 个用户到索引\n", len(users))
	return nil
}

// GetSuggestions 获取搜索建议
func (s *SearchService) GetSuggestions(query string, limit int) ([]string, error) {
	if limit < 1 || limit > 10 {
		limit = 5
	}

	// 搜索帖子标题
	searchRequest := &meilisearch.SearchRequest{
		Limit:                int64(limit),
		AttributesToRetrieve: []string{"title"},
	}

	result, err := s.client.Index("posts").Search(query, searchRequest)
	if err != nil {
		return nil, err
	}

	// 提取去重的标题
	seen := make(map[string]bool)
	var suggestions []string
	for _, hit := range result.Hits {
		hitJSON, _ := json.Marshal(hit)
		var hitMap map[string]interface{}
		json.Unmarshal(hitJSON, &hitMap)

		if title, ok := hitMap["title"]; ok {
			titleStr, _ := title.(string)
			if titleStr != "" && !seen[titleStr] {
				seen[titleStr] = true
				suggestions = append(suggestions, titleStr)
			}
		}
	}

	return suggestions, nil
}
