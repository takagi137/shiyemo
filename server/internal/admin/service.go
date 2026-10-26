package admin

import (
	"errors"

	"github.com/shiyemo/server/internal/auth"
	"gorm.io/gorm"
)

var ErrPostNotFound = errors.New("帖子不存在")

// Service 管理员服务
type Service struct {
	repo RepoInterface
}

// NewService 创建管理员服务
func NewService() *Service {
	return &Service{
		repo: NewRepo(),
	}
}

// NewServiceWithRepo 使用注入的 repo 创建管理员服务（测试用）
func NewServiceWithRepo(repo RepoInterface) *Service {
	return &Service{repo: repo}
}

// ListPosts 获取帖子列表
func (s *Service) ListPosts(page, pageSize int, status string) ([]map[string]interface{}, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	if status == "" {
		status = "active"
	}

	return s.repo.ListPosts(page, pageSize, status)
}

// DeletePost 删除帖子（管理员权限）
func (s *Service) DeletePost(postID int64) error {
	return s.repo.DeletePost(postID)
}

// ListUsers 获取用户列表（支持搜索和状态筛选）
func (s *Service) ListUsers(page, pageSize int, search string, status *int) ([]auth.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	return s.repo.ListUsers(page, pageSize, search, status)
}

// UpdateUserStatus 更新用户状态
func (s *Service) UpdateUserStatus(userID int64, status int8) error {
	// 检查用户是否存在
	if _, err := s.repo.GetUserByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	return s.repo.UpdateUserStatus(userID, status)
}

// GetStats 获取平台统计数据
func (s *Service) GetStats() (map[string]interface{}, error) {
	return s.repo.GetStats()
}

// SetPostTop 设置帖子置顶
func (s *Service) SetPostTop(postID int64, isTop bool) error {
	return s.repo.SetPostTop(postID, isTop)
}

// SetPostFeatured 设置帖子加精
func (s *Service) SetPostFeatured(postID int64, isFeatured bool) error {
	return s.repo.SetPostFeatured(postID, isFeatured)
}

// ListPendingPosts 获取待审核帖子列表
func (s *Service) ListPendingPosts(page, pageSize int) ([]map[string]interface{}, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	return s.repo.ListPendingPosts(page, pageSize)
}

// ApprovePost 审核通过帖子
func (s *Service) ApprovePost(postID, reviewerID int64) error {
	err := s.repo.ApprovePost(postID, reviewerID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrPostNotFound
	}
	return err
}

// RejectPost 审核拒绝帖子
func (s *Service) RejectPost(postID, reviewerID int64, reason string) error {
	if reason == "" {
		return errors.New("拒绝原因不能为空")
	}
	err := s.repo.RejectPost(postID, reviewerID, reason)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrPostNotFound
	}
	return err
}

// ListComments 获取评论列表（管理员）
func (s *Service) ListComments(page, pageSize int, postID int64) ([]map[string]interface{}, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.ListComments(page, pageSize, postID)
}

// HideComment 屏蔽评论
func (s *Service) HideComment(commentID int64) error {
	return s.repo.UpdateCommentStatus(commentID, 2)
}

// RestoreComment 恢复评论
func (s *Service) RestoreComment(commentID int64) error {
	return s.repo.UpdateCommentStatus(commentID, 1)
}
