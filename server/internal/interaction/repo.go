package interaction

import (
	"errors"

	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/database"
	"github.com/shiyemo/server/internal/post"
	"gorm.io/gorm"
)

// InteractionRepository 互动仓储接口
type InteractionRepository interface {
	Like(userID, postID int64) error
	Unlike(userID, postID int64) error
	IsLiked(userID, postID int64) bool
	Bookmark(userID, postID int64) error
	Unbookmark(userID, postID int64) error
	IsBookmarked(userID, postID int64) bool
	CountLikes(postID int64) (int64, error)
	ListBookmarks(userID int64, page, size int) ([]int64, int64, error)
}

// PostProvider 帖子查询接口
type PostProvider interface {
	GetByID(id int64) (*post.Post, error)
}

// UserProvider 用户查询接口
type UserProvider interface {
	GetByID(id int64) (*auth.User, error)
}

// Notifier 通知发送接口
type Notifier interface {
	Create(userID, actorID int64, nType string, targetID *int64, targetType *string, content string) error
}

// LikeCountUpdater 点赞计数更新接口
type LikeCountUpdater interface {
	UpdateLikeCount(postID int64, count int64)
}

// InteractionRepo 互动仓储（实现 InteractionRepository）
type InteractionRepo struct {
	db *gorm.DB
}

// NewInteractionRepo 创建互动仓储
func NewInteractionRepo() *InteractionRepo {
	return &InteractionRepo{
		db: database.GetDB(),
	}
}

// NewInteractionRepoWithDB 使用注入的 DB 创建互动仓储（测试用）
func NewInteractionRepoWithDB(db *gorm.DB) *InteractionRepo {
	return &InteractionRepo{
		db: db,
	}
}

// Like 点赞
func (r *InteractionRepo) Like(userID, postID int64) error {
	like := &Like{
		UserID: userID,
		PostID: postID,
	}
	return r.db.Create(like).Error
}

// Unlike 取消点赞
func (r *InteractionRepo) Unlike(userID, postID int64) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&Like{}).Error
}

// IsLiked 是否已点赞
func (r *InteractionRepo) IsLiked(userID, postID int64) bool {
	var count int64
	r.db.Model(&Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}

// Bookmark 收藏
func (r *InteractionRepo) Bookmark(userID, postID int64) error {
	bookmark := &Bookmark{
		UserID: userID,
		PostID: postID,
	}
	return r.db.Create(bookmark).Error
}

// Unbookmark 取消收藏
func (r *InteractionRepo) Unbookmark(userID, postID int64) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&Bookmark{}).Error
}

// IsBookmarked 是否已收藏
func (r *InteractionRepo) IsBookmarked(userID, postID int64) bool {
	var count int64
	r.db.Model(&Bookmark{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}

// CountLikes 统计点赞数
func (r *InteractionRepo) CountLikes(postID int64) (int64, error) {
	var count int64
	err := r.db.Model(&Like{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

// ListBookmarks 获取用户收藏列表
func (r *InteractionRepo) ListBookmarks(userID int64, page, size int) ([]int64, int64, error) {
	var bookmarks []Bookmark
	var total int64

	query := r.db.Model(&Bookmark{}).Where("user_id = ?", userID)

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * size
	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&bookmarks).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	var postIDs []int64
	for _, b := range bookmarks {
		postIDs = append(postIDs, b.PostID)
	}

	return postIDs, total, nil
}

// DBLikeCountUpdater 基于数据库的点赞计数更新器
type DBLikeCountUpdater struct{}

func (u *DBLikeCountUpdater) UpdateLikeCount(postID int64, count int64) {
	database.GetDB().Model(&post.Post{}).Where("id = ?", postID).Update("like_count", count)
}
