package interaction

import (
	"fmt"

	"github.com/shiyemo/server/internal/auth"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/notification"
	"github.com/shiyemo/server/internal/post"
)

// InteractionService 互动服务
type InteractionService struct {
	repo       InteractionRepository
	postRepo   PostProvider
	userRepo   UserProvider
	notifySvc  Notifier
	countUpd   LikeCountUpdater
}

// NewInteractionService 创建互动服务（使用数据库）
func NewInteractionService() *InteractionService {
	return &InteractionService{
		repo:      NewInteractionRepo(),
		postRepo:  post.NewPostRepo(),
		userRepo:  auth.NewAuthRepo(),
		notifySvc: notification.NewService(),
		countUpd:  &DBLikeCountUpdater{},
	}
}

// NewInteractionServiceWithDeps 创建互动服务（注入依赖，用于测试）
func NewInteractionServiceWithDeps(
	repo InteractionRepository,
	postRepo PostProvider,
	userRepo UserProvider,
	notifySvc Notifier,
	countUpd LikeCountUpdater,
) *InteractionService {
	return &InteractionService{
		repo:      repo,
		postRepo:  postRepo,
		userRepo:  userRepo,
		notifySvc: notifySvc,
		countUpd:  countUpd,
	}
}

// Like 点赞/取消点赞（切换）。返回：切换后的点赞态 + 重算后的权威点赞数。
// 返回权威 count 让前端直接采用，避免"乐观 ±1"在状态不同步时计数漂移。
func (s *InteractionService) Like(userID, postID int64) (bool, int64, error) {
	// 检查帖子是否存在
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		return false, 0, bizerr.ErrPostNotFound
	}

	// 检查是否已点赞
	if s.repo.IsLiked(userID, postID) {
		// 取消点赞
		if err := s.repo.Unlike(userID, postID); err != nil {
			return false, 0, bizerr.ErrDB
		}
		// 更新计数
		count, _ := s.repo.CountLikes(postID)
		s.countUpd.UpdateLikeCount(postID, count)
		return false, count, nil
	}

	// 点赞
	if err := s.repo.Like(userID, postID); err != nil {
		return false, 0, bizerr.ErrDB
	}

	// 更新计数
	count, _ := s.repo.CountLikes(postID)
	s.countUpd.UpdateLikeCount(postID, count)

	// 发送通知给帖子作者
	if p, err := s.postRepo.GetByID(postID); err == nil {
		if user, err := s.userRepo.GetByID(userID); err == nil {
			targetID := postID
			targetType := "post"
			nContent := fmt.Sprintf("%s 点赞了你的帖子\"%s\"", user.Nickname, p.Title)
			s.notifySvc.Create(p.UserID, userID, "like", &targetID, &targetType, nContent)
		}
	}

	return true, count, nil
}

// GetUserPostState 返回当前用户对某帖子的点赞/收藏状态（供前端初始化按钮态）。
func (s *InteractionService) GetUserPostState(userID, postID int64) (liked, bookmarked bool) {
	return s.repo.IsLiked(userID, postID), s.repo.IsBookmarked(userID, postID)
}

// Bookmark 收藏/取消收藏（切换）
func (s *InteractionService) Bookmark(userID, postID int64) (bool, error) {
	// 检查帖子是否存在
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		return false, bizerr.ErrPostNotFound
	}

	// 检查是否已收藏
	if s.repo.IsBookmarked(userID, postID) {
		// 取消收藏
		if err := s.repo.Unbookmark(userID, postID); err != nil {
			return false, bizerr.ErrDB
		}
		return false, nil
	}

	// 收藏
	if err := s.repo.Bookmark(userID, postID); err != nil {
		return false, bizerr.ErrDB
	}

	return true, nil
}

// ListBookmarks 获取用户收藏列表
func (s *InteractionService) ListBookmarks(userID int64, page, size int) ([]map[string]interface{}, int64, error) {
	postIDs, total, err := s.repo.ListBookmarks(userID, page, size)
	if err != nil {
		return nil, 0, bizerr.ErrDB
	}

	var result []map[string]interface{}
	for _, postID := range postIDs {
		p, err := s.postRepo.GetByID(postID)
		if err != nil {
			continue
		}
		result = append(result, map[string]interface{}{
			"id":         p.ID,
			"type":       p.Type,
			"title":      p.Title,
			"content":    p.Content,
			"like_count": p.LikeCount,
			"created_at": p.CreatedAt,
		})
	}

	return result, total, nil
}
