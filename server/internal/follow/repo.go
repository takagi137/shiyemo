package follow

import (
	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// FollowRepository 关注仓储接口
type FollowRepository interface {
	Follow(followerID, followingID int64) error
	Unfollow(followerID, followingID int64) error
	IsFollowing(followerID, followingID int64) bool
	ListFollowing(userID int64, page, size int) ([]int64, int64, error)
	ListFollowers(userID int64, page, size int) ([]int64, int64, error)
	CountFollowing(userID int64) (int64, error)
	CountFollowers(userID int64) (int64, error)
}

// UserProvider 用户查询接口
type UserProvider interface {
	GetByID(id int64) (*auth.User, error)
}

// Notifier 通知发送接口
type Notifier interface {
	Create(userID, actorID int64, nType string, targetID *int64, targetType *string, content string) error
}

// FollowCountUpdater 关注计数更新接口
type FollowCountUpdater interface {
	UpdateCounts(followerID, followingID int64, followingCount, followerCount int64)
}

// FollowRepo 关注仓储（实现 FollowRepository）
type FollowRepo struct {
	db *gorm.DB
}

// NewFollowRepo 创建关注仓储
func NewFollowRepo() *FollowRepo {
	return &FollowRepo{
		db: database.GetDB(),
	}
}

// NewFollowRepoWithDB 使用注入的 DB 创建关注仓储（测试用）
func NewFollowRepoWithDB(db *gorm.DB) *FollowRepo {
	return &FollowRepo{
		db: db,
	}
}

// Follow 关注
func (r *FollowRepo) Follow(followerID, followingID int64) error {
	follow := &Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	return r.db.Create(follow).Error
}

// Unfollow 取消关注
func (r *FollowRepo) Unfollow(followerID, followingID int64) error {
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&Follow{}).Error
}

// IsFollowing 是否已关注
func (r *FollowRepo) IsFollowing(followerID, followingID int64) bool {
	var count int64
	r.db.Model(&Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count)
	return count > 0
}

// ListFollowing 获取关注列表
func (r *FollowRepo) ListFollowing(userID int64, page, size int) ([]int64, int64, error) {
	var follows []Follow
	var total int64

	query := r.db.Model(&Follow{}).Where("follower_id = ?", userID)

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * size
	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&follows).Error
	if err != nil {
		return nil, 0, err
	}

	var userIDs []int64
	for _, f := range follows {
		userIDs = append(userIDs, f.FollowingID)
	}

	return userIDs, total, nil
}

// ListFollowers 获取粉丝列表
func (r *FollowRepo) ListFollowers(userID int64, page, size int) ([]int64, int64, error) {
	var follows []Follow
	var total int64

	query := r.db.Model(&Follow{}).Where("following_id = ?", userID)

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * size
	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&follows).Error
	if err != nil {
		return nil, 0, err
	}

	var userIDs []int64
	for _, f := range follows {
		userIDs = append(userIDs, f.FollowerID)
	}

	return userIDs, total, nil
}

// CountFollowing 统计关注数
func (r *FollowRepo) CountFollowing(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&Follow{}).Where("follower_id = ?", userID).Count(&count).Error
	return count, err
}

// CountFollowers 统计粉丝数
func (r *FollowRepo) CountFollowers(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&Follow{}).Where("following_id = ?", userID).Count(&count).Error
	return count, err
}

// DBFollowCountUpdater 基于数据库的关注计数更新器
type DBFollowCountUpdater struct{}

func (u *DBFollowCountUpdater) UpdateCounts(followerID, followingID int64, followingCount, followerCount int64) {
	database.GetDB().Model(&auth.User{}).Where("id = ?", followerID).Update("following_count", followingCount)
	database.GetDB().Model(&auth.User{}).Where("id = ?", followingID).Update("follower_count", followerCount)
}
