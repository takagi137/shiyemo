package comment

import (
	"errors"
	"fmt"

	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/anonymous"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/notification"
	"github.com/shiyemo/server/internal/post"
	"gorm.io/gorm"
)

// CommentService 评论服务
type CommentService struct {
	repo       CommentRepository
	userRepo   UserProvider
	postRepo   PostProvider
	notifySvc  Notifier
	countUpd   CommentCountUpdater
}

// NewCommentService 创建评论服务（使用数据库）
func NewCommentService() *CommentService {
	return &CommentService{
		repo:      NewCommentRepo(),
		userRepo:  auth.NewAuthRepo(),
		postRepo:  post.NewPostRepo(),
		notifySvc: notification.NewService(),
		countUpd:  &DBCommentCountUpdater{},
	}
}

// NewCommentServiceWithDeps 创建评论服务（注入依赖，用于测试）
func NewCommentServiceWithDeps(
	repo CommentRepository,
	userRepo UserProvider,
	postRepo PostProvider,
	notifySvc Notifier,
	countUpd CommentCountUpdater,
) *CommentService {
	return &CommentService{
		repo:      repo,
		userRepo:  userRepo,
		postRepo:  postRepo,
		notifySvc: notifySvc,
		countUpd:  countUpd,
	}
}

// Create 创建评论
func (s *CommentService) Create(userID, postID int64, req CreateRequest) (map[string]interface{}, error) {
	// 检查帖子是否存在
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerr.ErrPostNotFound
		}
		return nil, bizerr.ErrDB
	}

	// 如果是回复，检查父评论是否存在
	if req.ParentID != nil {
		parentComment, err := s.repo.GetByID(*req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, bizerr.New(40301, "父评论不存在")
			}
			return nil, bizerr.ErrDB
		}
		// 检查父评论是否属于同一帖子
		if parentComment.PostID != postID {
			return nil, bizerr.New(40301, "父评论不属于该帖子")
		}
		// 检查是否是回复的回复（最多两级）
		if parentComment.ParentID != nil {
			return nil, bizerr.New(40301, "最多支持两级评论")
		}
	}

	comment := &Comment{
		PostID:      postID,
		UserID:      userID,
		ParentID:    req.ParentID,
		Content:     req.Content,
		IsAnonymous: req.IsAnonymous,
		Status:      1,
	}

	// 如果是匿名发布，生成匿名名称
	if req.IsAnonymous {
		comment.AnonymousName = anonymous.Generate()
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, bizerr.ErrDB
	}

	// 更新帖子的评论计数
	count, _ := s.repo.CountByPostID(postID)
	s.countUpd.UpdateCommentCount(postID, count)

	// 获取作者信息
	author, _ := s.userRepo.GetByID(userID)

	// 发送通知给帖子作者
	if p, err := s.postRepo.GetByID(postID); err == nil {
		targetID := postID
		targetType := "post"
		nContent := fmt.Sprintf("%s 评论了你的帖子\"%s\"", author.Nickname, p.Title)
		s.notifySvc.Create(p.UserID, userID, "comment", &targetID, &targetType, nContent)
	}

	// 如果是回复，发送通知给被回复的评论作者
	if req.ParentID != nil {
		parentComment, err := s.repo.GetByID(*req.ParentID)
		if err == nil && parentComment.UserID != userID {
			targetID := postID
			targetType := "comment"
			nContent := fmt.Sprintf("%s 回复了你的评论", author.Nickname)
			s.notifySvc.Create(parentComment.UserID, userID, "reply", &targetID, &targetType, nContent)
		}
	}

	return map[string]interface{}{
		"id":         comment.ID,
		"post_id":    comment.PostID,
		"user_id":    comment.UserID,
		"parent_id":  comment.ParentID,
		"content":    comment.Content,
		"like_count": comment.LikeCount,
		"created_at": comment.CreatedAt,
		"author": map[string]interface{}{
			"id":       author.ID,
			"nickname": author.Nickname,
			"avatar":   author.Avatar,
		},
	}, nil
}

// List 获取帖子评论列表
func (s *CommentService) List(postID int64, page, size int) ([]map[string]interface{}, int64, error) {
	comments, total, err := s.repo.GetByPostID(postID, page, size)
	if err != nil {
		return nil, 0, bizerr.ErrDB
	}

	// 收集所有需要查询的用户 ID
	userIDs := make(map[int64]bool)
	for _, comment := range comments {
		userIDs[comment.UserID] = true
	}

	// 批量获取所有回复
	allReplies := make(map[int64][]Comment)
	for _, comment := range comments {
		replies, _ := s.repo.GetReplies(comment.ID)
		allReplies[comment.ID] = replies
		for _, reply := range replies {
			userIDs[reply.UserID] = true
		}
	}

	// 批量查询用户信息
	users := make(map[int64]*auth.User)
	for uid := range userIDs {
		user, err := s.userRepo.GetByID(uid)
		if err == nil {
			users[uid] = user
		}
	}

	var result []map[string]interface{}
	for _, comment := range comments {
		author := users[comment.UserID]
		if author == nil {
			author = &auth.User{ID: comment.UserID, Nickname: "未知用户"}
		}

		var replyList []map[string]interface{}
		for _, reply := range allReplies[comment.ID] {
			replyAuthor := users[reply.UserID]
			if replyAuthor == nil {
				replyAuthor = &auth.User{ID: reply.UserID, Nickname: "未知用户"}
			}

			replyItem := map[string]interface{}{
				"id":             reply.ID,
				"post_id":        reply.PostID,
				"user_id":        reply.UserID,
				"parent_id":      reply.ParentID,
				"content":        reply.Content,
				"like_count":     reply.LikeCount,
				"is_anonymous":   reply.IsAnonymous,
				"created_at":     reply.CreatedAt,
			}

			// 根据匿名设置返回作者信息
			if reply.IsAnonymous {
				replyItem["author"] = map[string]interface{}{
					"id":       0,
					"nickname": reply.AnonymousName,
					"avatar":   "",
				}
				replyItem["user_id"] = 0
			} else {
				replyItem["author"] = map[string]interface{}{
					"id":       replyAuthor.ID,
					"nickname": replyAuthor.Nickname,
					"avatar":   replyAuthor.Avatar,
				}
			}

			replyList = append(replyList, replyItem)
		}

		commentItem := map[string]interface{}{
			"id":             comment.ID,
			"post_id":        comment.PostID,
			"user_id":        comment.UserID,
			"parent_id":      comment.ParentID,
			"content":        comment.Content,
			"like_count":     comment.LikeCount,
			"is_anonymous":   comment.IsAnonymous,
			"created_at":     comment.CreatedAt,
			"replies":        replyList,
		}

		// 根据匿名设置返回作者信息
		if comment.IsAnonymous {
			commentItem["author"] = map[string]interface{}{
				"id":       0,
				"nickname": comment.AnonymousName,
				"avatar":   "",
			}
			commentItem["user_id"] = 0
		} else {
			commentItem["author"] = map[string]interface{}{
				"id":       author.ID,
				"nickname": author.Nickname,
				"avatar":   author.Avatar,
			}
		}

		result = append(result, commentItem)
	}

	return result, total, nil
}

// Delete 删除评论
func (s *CommentService) Delete(userID, commentID int64) error {
	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerr.ErrPostNotFound
		}
		return bizerr.ErrDB
	}

	// 检查权限
	if comment.UserID != userID {
		return bizerr.ErrPostForbidden
	}

	// 删除评论
	if err := s.repo.Delete(commentID); err != nil {
		return bizerr.ErrDB
	}

	// 更新帖子的评论计数
	count, _ := s.repo.CountByPostID(comment.PostID)
	s.countUpd.UpdateCommentCount(comment.PostID, count)

	return nil
}

// CreateRequest 创建评论请求
type CreateRequest struct {
	Content     string `json:"content" binding:"required,max=1000"`
	ParentID    *int64 `json:"parent_id"`
	IsAnonymous bool   `json:"is_anonymous"`
}
