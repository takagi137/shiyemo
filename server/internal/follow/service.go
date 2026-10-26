package follow

import (
	"errors"
	"fmt"

	"github.com/shiyemo/server/internal/auth"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/notification"
	"gorm.io/gorm"
)

// FollowService 关注服务
type FollowService struct {
	repo      FollowRepository
	userRepo  UserProvider
	notifySvc Notifier
	countUpd  FollowCountUpdater
}

// NewFollowService 创建关注服务（使用数据库）
func NewFollowService() *FollowService {
	return &FollowService{
		repo:      NewFollowRepo(),
		userRepo:  auth.NewAuthRepo(),
		notifySvc: notification.NewService(),
		countUpd:  &DBFollowCountUpdater{},
	}
}

// NewFollowServiceWithDeps 创建关注服务（注入依赖，用于测试）
func NewFollowServiceWithDeps(
	repo FollowRepository,
	userRepo UserProvider,
	notifySvc Notifier,
	countUpd FollowCountUpdater,
) *FollowService {
	return &FollowService{
		repo:      repo,
		userRepo:  userRepo,
		notifySvc: notifySvc,
		countUpd:  countUpd,
	}
}

// Follow 关注/取消关注（切换）
func (s *FollowService) Follow(followerID, followingID int64) (bool, error) {
	// 检查是否关注自己
	if followerID == followingID {
		return false, bizerr.New(40301, "不能关注自己")
	}

	// 检查用户是否存在
	_, err := s.userRepo.GetByID(followingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, bizerr.ErrUserNotFound
		}
		return false, bizerr.ErrDB
	}

	// 检查是否已关注
	if s.repo.IsFollowing(followerID, followingID) {
		// 取消关注
		if err := s.repo.Unfollow(followerID, followingID); err != nil {
			return false, bizerr.ErrDB
		}
		// 更新计数
		s.updateCounts(followerID, followingID)
		return false, nil
	}

	// 关注
	if err := s.repo.Follow(followerID, followingID); err != nil {
		return false, bizerr.ErrDB
	}

	// 更新计数
	s.updateCounts(followerID, followingID)

	// 发送通知给被关注者
	if follower, err := s.userRepo.GetByID(followerID); err == nil {
		nContent := fmt.Sprintf("%s 关注了你", follower.Nickname)
		s.notifySvc.Create(followingID, followerID, "follow", nil, nil, nContent)
	}

	return true, nil
}

// ListFollowing 获取关注列表
func (s *FollowService) ListFollowing(userID int64, page, size int) ([]map[string]interface{}, int64, error) {
	userIDs, total, err := s.repo.ListFollowing(userID, page, size)
	if err != nil {
		return nil, 0, bizerr.ErrDB
	}

	var result []map[string]interface{}
	for _, uid := range userIDs {
		user, err := s.userRepo.GetByID(uid)
		if err != nil {
			continue
		}
		result = append(result, map[string]interface{}{
			"id":       user.ID,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"bio":      user.Bio,
		})
	}

	return result, total, nil
}

// ListFollowers 获取粉丝列表
func (s *FollowService) ListFollowers(userID int64, page, size int) ([]map[string]interface{}, int64, error) {
	userIDs, total, err := s.repo.ListFollowers(userID, page, size)
	if err != nil {
		return nil, 0, bizerr.ErrDB
	}

	var result []map[string]interface{}
	for _, uid := range userIDs {
		user, err := s.userRepo.GetByID(uid)
		if err != nil {
			continue
		}
		result = append(result, map[string]interface{}{
			"id":       user.ID,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"bio":      user.Bio,
		})
	}

	return result, total, nil
}

// updateCounts 更新关注/粉丝计数
func (s *FollowService) updateCounts(followerID, followingID int64) {
	followingCount, _ := s.repo.CountFollowing(followerID)
	followerCount, _ := s.repo.CountFollowers(followingID)
	s.countUpd.UpdateCounts(followerID, followingID, followingCount, followerCount)
}
