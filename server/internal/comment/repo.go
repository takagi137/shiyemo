package comment

import (
	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/database"
	"github.com/shiyemo/server/internal/post"
	"gorm.io/gorm"
)

// CommentRepository 评论仓储接口
type CommentRepository interface {
	Create(comment *Comment) error
	GetByID(id int64) (*Comment, error)
	GetByPostID(postID int64, page, size int) ([]Comment, int64, error)
	GetReplies(parentID int64) ([]Comment, error)
	Delete(id int64) error
	CountByPostID(postID int64) (int64, error)
}

// UserProvider 用户查询接口
type UserProvider interface {
	GetByID(id int64) (*auth.User, error)
}

// PostProvider 帖子查询接口
type PostProvider interface {
	GetByID(id int64) (*post.Post, error)
}

// Notifier 通知发送接口
type Notifier interface {
	Create(userID, actorID int64, nType string, targetID *int64, targetType *string, content string) error
}

// CommentCountUpdater 评论计数更新接口
type CommentCountUpdater interface {
	UpdateCommentCount(postID int64, count int64)
}

// CommentRepo 评论仓储（实现 CommentRepository）
type CommentRepo struct {
	db *gorm.DB
}

// NewCommentRepo 创建评论仓储
func NewCommentRepo() *CommentRepo {
	return &CommentRepo{
		db: database.GetDB(),
	}
}

// NewCommentRepoWithDB 使用注入的 DB 创建评论仓储（测试用）
func NewCommentRepoWithDB(db *gorm.DB) *CommentRepo {
	return &CommentRepo{
		db: db,
	}
}

// Create 创建评论
func (r *CommentRepo) Create(comment *Comment) error {
	return r.db.Create(comment).Error
}

// GetByID 根据 ID 查询评论
func (r *CommentRepo) GetByID(id int64) (*Comment, error) {
	var comment Comment
	err := r.db.Where("id = ? AND status = 1", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetByPostID 获取帖子的评论列表
func (r *CommentRepo) GetByPostID(postID int64, page, size int) ([]Comment, int64, error) {
	var comments []Comment
	var total int64

	query := r.db.Model(&Comment{}).Where("post_id = ? AND status = 1 AND parent_id IS NULL", postID)

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * size
	err := query.Order("created_at ASC").Offset(offset).Limit(size).Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// GetReplies 获取评论的回复列表
func (r *CommentRepo) GetReplies(parentID int64) ([]Comment, error) {
	var replies []Comment
	err := r.db.Where("parent_id = ? AND status = 1", parentID).Order("created_at ASC").Find(&replies).Error
	if err != nil {
		return nil, err
	}
	return replies, nil
}

// Delete 删除评论
func (r *CommentRepo) Delete(id int64) error {
	return r.db.Model(&Comment{}).Where("id = ?", id).Update("status", 3).Error
}

// CountByPostID 统计帖子的评论数
func (r *CommentRepo) CountByPostID(postID int64) (int64, error) {
	var count int64
	err := r.db.Model(&Comment{}).Where("post_id = ? AND status = 1", postID).Count(&count).Error
	return count, err
}

// DBCommentCountUpdater 基于数据库的评论计数更新器
type DBCommentCountUpdater struct{}

func (u *DBCommentCountUpdater) UpdateCommentCount(postID int64, count int64) {
	database.GetDB().Model(&post.Post{}).Where("id = ?", postID).Update("comment_count", count)
}
