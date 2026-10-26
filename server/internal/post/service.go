package post

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/anonymous"
	"github.com/shiyemo/server/internal/common/cache"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/moderation"
	"gorm.io/gorm"
)

// UserGetter 用户查询接口（用于帖子获取作者信息）
type UserGetter interface {
	GetByID(id int64) (*auth.User, error)
}

// ContentModerator 内容审核接口（用于帖子创建审核）
type ContentModerator interface {
	CheckTitle(title string) *moderation.CheckResult
	CheckContent(text string) *moderation.CheckResult
}

// PostIndexer 帖子搜索索引接口（解耦 search 包的循环依赖）
type PostIndexer interface {
	IndexPost(p *Post) error
	DeletePost(postID int64) error
}

// PostService 帖子服务
type PostService struct {
	postRepo     PostRepoInterface
	userRepo     UserGetter
	moderator    ContentModerator
	indexer      PostIndexer
	sf           *cache.Singleflight // 防缓存雪崩
}

// NewPostService 创建帖子服务
func NewPostService() *PostService {
	return &PostService{
		postRepo:  NewPostRepo(),
		userRepo:  auth.NewAuthRepo(),
		moderator: moderation.NewModerationService(),
		sf:        cache.NewSingleflight(),
	}
}

// SetIndexer 设置搜索索引器（延迟注入，避免循环依赖）
func (s *PostService) SetIndexer(indexer PostIndexer) {
	s.indexer = indexer
}

// NewPostServiceWithDeps 使用注入的依赖创建帖子服务（测试用）
func NewPostServiceWithDeps(postRepo PostRepoInterface, userRepo UserGetter, moderator ContentModerator) *PostService {
	return &PostService{
		postRepo:  postRepo,
		userRepo:  userRepo,
		moderator: moderator,
		sf:        cache.NewSingleflight(),
	}
}

// Create 创建帖子
func (s *PostService) Create(userID int64, req CreateRequest) (*Post, error) {
	// 内容审核 - 检查标题
	if result := s.moderator.CheckTitle(req.Title); !result.Passed {
		return nil, bizerr.ErrContentBlocked
	}

	// 内容审核 - 检查正文
	if result := s.moderator.CheckContent(req.Content); !result.Passed {
		return nil, bizerr.ErrContentBlocked
	}

	post := &Post{
		UserID:      userID,
		Type:        req.Type,
		Title:       req.Title,
		Content:     req.Content,
		IsAnonymous: req.IsAnonymous,
		Status:      1,
	}

	// 如果是匿名发布，生成匿名名称
	if req.IsAnonymous {
		post.AnonymousName = anonymous.Generate()
	}

	post.SetImages(req.Images)
	post.SetTags(req.Tags)

	if err := s.postRepo.Create(post); err != nil {
		return nil, bizerr.ErrDB
	}

	// 清除帖子列表缓存
	s.invalidatePostListCache()

	// 同步到搜索索引
	if s.indexer != nil {
		if err := s.indexer.IndexPost(post); err != nil {
			log.Printf("索引帖子 %d 失败: %v", post.ID, err)
		}
	}

	return post, nil
}

// GetDetail 获取帖子详情
func (s *PostService) GetDetail(id int64) (map[string]interface{}, error) {
	ctx := context.Background()
	cacheKey := cache.FormatPostDetailKey(id)

	// 尝试从缓存获取
	var cachedResult map[string]interface{}
	if err := cache.Get(ctx, cacheKey, &cachedResult); err == nil {
		return cachedResult, nil
	}

	// 缓存未命中，使用 singleflight 防止缓存雪崩
	v, err := s.sf.Do(cacheKey, func() (interface{}, error) {
		return s.fetchPostDetail(id)
	})
	if err != nil {
		return nil, err
	}
	result := v.(map[string]interface{})

	// 写入缓存
	if err := cache.Set(ctx, cacheKey, result, cache.CacheExpirePostDetail); err != nil {
		log.Printf("Failed to cache post detail %d: %v", id, err)
	}

	return result, nil
}

// fetchPostDetail 从数据库获取帖子详情（singleflight 回调）
func (s *PostService) fetchPostDetail(id int64) (map[string]interface{}, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerr.ErrPostNotFound
		}
		return nil, bizerr.ErrDB
	}

	author, err := s.userRepo.GetByID(post.UserID)
	if err != nil {
		return nil, bizerr.ErrDB
	}

	result := map[string]interface{}{
		"id":             post.ID,
		"user_id":        post.UserID,
		"type":           post.Type,
		"title":          post.Title,
		"content":        post.Content,
		"images":         post.GetImages(),
		"tags":           post.GetTags(),
		"like_count":     post.LikeCount,
		"comment_count":  post.CommentCount,
		"is_anonymous":   post.IsAnonymous,
		"status":         post.Status,
		"created_at":     post.CreatedAt,
		"updated_at":     post.UpdatedAt,
	}

	// 根据匿名设置返回作者信息
	if post.IsAnonymous {
		result["author"] = map[string]interface{}{
			"id":       0,
			"nickname": post.AnonymousName,
			"avatar":   "",
		}
	} else {
		result["author"] = map[string]interface{}{
			"id":       author.ID,
			"nickname": author.Nickname,
			"avatar":   author.Avatar,
		}
	}

	return result, nil
}

// List 获取帖子列表
func (s *PostService) List(postType string, page, size int) ([]map[string]interface{}, int64, error) {
	ctx := context.Background()
	cacheKey := cache.FormatPostListKey(postType, page, size)

	// 尝试从缓存获取
	type cachedList struct {
		Items []map[string]interface{} `json:"items"`
		Total int64                    `json:"total"`
	}
	var cached cachedList
	if err := cache.Get(ctx, cacheKey, &cached); err == nil {
		log.Printf("Cache hit for post list: %s page %d", postType, page)
		return cached.Items, cached.Total, nil
	}

	// 缓存未命中，从数据库获取
	posts, total, err := s.postRepo.List(postType, page, size)
	if err != nil {
		return nil, 0, bizerr.ErrDB
	}

	var result []map[string]interface{}
	for _, post := range posts {
		item := map[string]interface{}{
			"id":             post.ID,
			"user_id":        post.UserID,
			"type":           post.Type,
			"title":          post.Title,
			"content":        post.Content,
			"images":         post.GetImages(),
			"tags":           post.GetTags(),
			"like_count":     post.LikeCount,
			"comment_count":  post.CommentCount,
			"is_anonymous":   post.IsAnonymous,
			"created_at":     post.CreatedAt,
		}

		// 匿名帖子不返回真实用户ID
		if post.IsAnonymous {
			item["user_id"] = 0
		}

		result = append(result, item)
	}

	// 写入缓存
	cacheData := cachedList{Items: result, Total: total}
	if err := cache.Set(ctx, cacheKey, cacheData, cache.CacheExpirePostList); err != nil {
		log.Printf("Failed to cache post list: %v", err)
	}

	return result, total, nil
}

// Update 更新帖子
func (s *PostService) Update(userID, postID int64, req UpdateRequest) error {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerr.ErrPostNotFound
		}
		return bizerr.ErrDB
	}

	// 检查权限
	if post.UserID != userID {
		return bizerr.ErrPostForbidden
	}

	// 更新字段
	fields := make(map[string]interface{})
	if req.Title != "" {
		fields["title"] = req.Title
	}
	if req.Content != "" {
		fields["content"] = req.Content
	}
	if req.Type != "" {
		fields["type"] = req.Type
	}

	if len(fields) > 0 {
		if err := s.postRepo.UpdateFields(postID, fields); err != nil {
			return err
		}
	}

	// 更新图片
	if req.Images != nil {
		post.SetImages(req.Images)
		if err := s.postRepo.UpdateFields(postID, map[string]interface{}{"images": post.Images}); err != nil {
			return err
		}
	}

	// 更新标签
	if req.Tags != nil {
		post.SetTags(req.Tags)
		if err := s.postRepo.UpdateFields(postID, map[string]interface{}{"tags": post.Tags}); err != nil {
			return err
		}
	}

	// 清除缓存
	s.invalidatePostCache(postID)

	return nil
}

// Delete 删除帖子
func (s *PostService) Delete(userID, postID int64) error {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerr.ErrPostNotFound
		}
		return bizerr.ErrDB
	}

	// 检查权限
	if post.UserID != userID {
		return bizerr.ErrPostForbidden
	}

	if err := s.postRepo.Delete(postID); err != nil {
		return err
	}

	// 清除缓存
	s.invalidatePostCache(postID)
	s.invalidatePostListCache()

	// 从搜索索引删除
	if s.indexer != nil {
		if err := s.indexer.DeletePost(postID); err != nil {
			log.Printf("从索引删除帖子 %d 失败: %v", postID, err)
		}
	}

	return nil
}

// FollowingFeed 获取关注用户的帖子列表
func (s *PostService) FollowingFeed(userID int64, page, size int) ([]map[string]interface{}, int64, error) {
	// 获取关注的用户 ID 列表
	followingIDs, err := s.postRepo.GetFollowingIDs(userID)
	if err != nil {
		return nil, 0, bizerr.ErrDB
	}

	if len(followingIDs) == 0 {
		return []map[string]interface{}{}, 0, nil
	}

	// 获取关注用户的帖子
	posts, total, err := s.postRepo.GetByUserIDs(followingIDs, page, size)
	if err != nil {
		return nil, 0, bizerr.ErrDB
	}

	var result []map[string]interface{}
	for _, post := range posts {
		result = append(result, map[string]interface{}{
			"id":            post.ID,
			"user_id":       post.UserID,
			"type":          post.Type,
			"title":         post.Title,
			"content":       post.Content,
			"images":        post.GetImages(),
			"tags":          post.GetTags(),
			"like_count":    post.LikeCount,
			"comment_count": post.CommentCount,
			"created_at":    post.CreatedAt,
		})
	}

	return result, total, nil
}

// CohortFeed 获取同温层帖子列表
func (s *PostService) CohortFeed(postType, city, industry string, page, size int) ([]map[string]interface{}, int64, error) {
	ctx := context.Background()
	cacheKey := cache.FormatPostListKey("cohort:"+city+":"+industry+":"+postType, page, size)

	// 尝试从缓存获取
	type cachedList struct {
		Items []map[string]interface{} `json:"items"`
		Total int64                    `json:"total"`
	}
	var cached cachedList
	if err := cache.Get(ctx, cacheKey, &cached); err == nil {
		log.Printf("Cache hit for cohort feed: %s %s %s page %d", city, industry, postType, page)
		return cached.Items, cached.Total, nil
	}

	// 缓存未命中，从数据库获取
	posts, total, err := s.postRepo.ListByCohort(postType, city, industry, page, size)
	if err != nil {
		return nil, 0, bizerr.ErrDB
	}

	var result []map[string]interface{}
	for _, post := range posts {
		item := map[string]interface{}{
			"id":             post.ID,
			"user_id":        post.UserID,
			"type":           post.Type,
			"title":          post.Title,
			"content":        post.Content,
			"images":         post.GetImages(),
			"tags":           post.GetTags(),
			"like_count":     post.LikeCount,
			"comment_count":  post.CommentCount,
			"is_anonymous":   post.IsAnonymous,
			"created_at":     post.CreatedAt,
		}

		// 匿名帖子不返回真实用户ID
		if post.IsAnonymous {
			item["user_id"] = 0
		}

		result = append(result, item)
	}

	// 写入缓存
	cacheData := cachedList{Items: result, Total: total}
	if err := cache.Set(ctx, cacheKey, cacheData, cache.CacheExpirePostList); err != nil {
		log.Printf("Failed to cache cohort feed: %v", err)
	}

	return result, total, nil
}

// CreateRequest 创建帖子请求
type CreateRequest struct {
	Type        string   `json:"type" binding:"required,oneof=emotion insight case action achievement"`
	Title       string   `json:"title" binding:"required,max=200"`
	Content     string   `json:"content" binding:"required"`
	Images      []string `json:"images"`
	Tags        []string `json:"tags"`
	IsAnonymous bool     `json:"is_anonymous"`
}

// UpdateRequest 更新帖子请求
type UpdateRequest struct {
	Type    string   `json:"type" binding:"omitempty,oneof=emotion insight case action achievement"`
	Title   string   `json:"title" binding:"omitempty,max=200"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
	Tags    []string `json:"tags"`
}

// ============================================================
// 缓存辅助方法
// ============================================================

// invalidatePostCache 清除单个帖子缓存
func (s *PostService) invalidatePostCache(postID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cacheKey := cache.FormatPostDetailKey(postID)
	if err := cache.Delete(ctx, cacheKey); err != nil {
		log.Printf("Failed to invalidate post cache %d: %v", postID, err)
	}
}

// invalidatePostListCache 清除所有帖子列表缓存
func (s *PostService) invalidatePostListCache() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pattern := "shiyemo:post:list:*"
	if err := cache.DeletePattern(ctx, pattern); err != nil {
		log.Printf("Failed to invalidate post list cache: %v", err)
	}
}
