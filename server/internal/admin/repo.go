package admin

import (
	"time"

	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// RepoInterface 管理员仓储接口（用于测试 mock）
type RepoInterface interface {
	ListPosts(page, pageSize int, status string) ([]map[string]interface{}, int64, error)
	DeletePost(postID int64) error
	ListUsers(page, pageSize int, search string, status *int) ([]auth.User, int64, error)
	UpdateUserStatus(userID int64, status int8) error
	GetUserByID(userID int64) (*auth.User, error)
	GetStats() (map[string]interface{}, error)
	SetPostTop(postID int64, isTop bool) error
	SetPostFeatured(postID int64, isFeatured bool) error
	ListPendingPosts(page, pageSize int) ([]map[string]interface{}, int64, error)
	ApprovePost(postID, reviewerID int64) error
	RejectPost(postID, reviewerID int64, reason string) error
	ListComments(page, pageSize int, postID int64) ([]map[string]interface{}, int64, error)
	UpdateCommentStatus(commentID int64, status int) error
}

// Repo 管理员仓储
type Repo struct {
	db *gorm.DB
}

// NewRepo 创建管理员仓储
func NewRepo() *Repo {
	return &Repo{
		db: database.GetDB(),
	}
}

// NewRepoWithDB 使用注入的 DB 创建管理员仓储（测试用）
func NewRepoWithDB(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// ListPosts 获取帖子列表（管理员视角）
func (r *Repo) ListPosts(page, pageSize int, status string) ([]map[string]interface{}, int64, error) {
	var total int64
	var posts []map[string]interface{}

	query := r.db.Table("posts")

	// 根据状态筛选
	switch status {
	case "deleted":
		query = query.Where("deleted_at IS NOT NULL")
	case "active":
		query = query.Where("deleted_at IS NULL AND status = 1")
	default: // all
		query = query.Where("1=1")
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	err := query.Select("id, user_id, type, title, content, status, created_at, deleted_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

// ListUsers 获取用户列表（支持搜索和状态筛选）
func (r *Repo) ListUsers(page, pageSize int, search string, status *int) ([]auth.User, int64, error) {
	var total int64
	var users []auth.User

	query := r.db.Model(&auth.User{})

	// 搜索：按昵称/邮箱/手机号模糊匹配
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("nickname LIKE ? OR email LIKE ? OR phone LIKE ?", like, like, like)
	}

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

// UpdateUserStatus 更新用户状态
func (r *Repo) UpdateUserStatus(userID int64, status int8) error {
	return r.db.Model(&auth.User{}).Where("id = ?", userID).Update("status", status).Error
}

// GetStats 获取平台统计数据
func (r *Repo) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总用户数
	var totalUsers int64
	r.db.Model(&auth.User{}).Count(&totalUsers)
	stats["total_users"] = totalUsers

	// 总帖子数
	var totalPosts int64
	r.db.Table("posts").Where("deleted_at IS NULL").Count(&totalPosts)
	stats["total_posts"] = totalPosts

	// 总评论数
	var totalComments int64
	r.db.Table("comments").Count(&totalComments)
	stats["total_comments"] = totalComments

	// 今日活跃用户数（24小时内有登录或发帖/评论的用户）
	var activeToday int64
	r.db.Model(&auth.User{}).Where("updated_at >= CURRENT_DATE").Count(&activeToday)
	stats["active_users_today"] = activeToday

	return stats, nil
}

// ListPendingPosts 获取待审核帖子列表
func (r *Repo) ListPendingPosts(page, pageSize int) ([]map[string]interface{}, int64, error) {
	var total int64
	var posts []map[string]interface{}

	query := r.db.Table("posts").Where("review_status = 'pending' AND deleted_at IS NULL")

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	err := query.Select("id, user_id, type, title, content, status, created_at").
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

// ApprovePost 审核通过帖子
func (r *Repo) ApprovePost(postID, reviewerID int64) error {
	now := time.Now()
	result := r.db.Table("posts").Where("id = ?", postID).Updates(map[string]interface{}{
		"review_status": "approved",
		"reviewed_at":   now,
		"reviewed_by":   reviewerID,
		"reject_reason": "",
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// RejectPost 审核拒绝帖子
func (r *Repo) RejectPost(postID, reviewerID int64, reason string) error {
	now := time.Now()
	result := r.db.Table("posts").Where("id = ?", postID).Updates(map[string]interface{}{
		"review_status": "rejected",
		"reviewed_at":   now,
		"reviewed_by":   reviewerID,
		"reject_reason": reason,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeletePost 删除帖子（管理员权限）
func (r *Repo) DeletePost(postID int64) error {
	return r.db.Table("posts").Where("id = ?", postID).Update("deleted_at", "NOW()").Error
}

// SetPostTop 设置帖子置顶
func (r *Repo) SetPostTop(postID int64, isTop bool) error {
	return r.db.Table("posts").Where("id = ?", postID).Update("is_top", isTop).Error
}

// SetPostFeatured 设置帖子加精
func (r *Repo) SetPostFeatured(postID int64, isFeatured bool) error {
	return r.db.Table("posts").Where("id = ?", postID).Update("is_featured", isFeatured).Error
}

// GetUserByID 根据 ID 获取用户
func (r *Repo) GetUserByID(userID int64) (*auth.User, error) {
	var user auth.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ListComments 获取评论列表（管理员）
func (r *Repo) ListComments(page, pageSize int, postID int64) ([]map[string]interface{}, int64, error) {
	var total int64
	var comments []map[string]interface{}

	query := r.db.Table("comments").Where("deleted_at IS NULL")
	if postID > 0 {
		query = query.Where("post_id = ?", postID)
	}

	query.Count(&total)

	rows := []struct {
		ID        int64     `gorm:"column:id"`
		PostID    int64     `gorm:"column:post_id"`
		UserID    int64     `gorm:"column:user_id"`
		Content   string    `gorm:"column:content"`
		Status    int       `gorm:"column:status"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}{}

	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	for _, row := range rows {
		comments = append(comments, map[string]interface{}{
			"id":         row.ID,
			"post_id":    row.PostID,
			"user_id":    row.UserID,
			"content":    row.Content,
			"status":     row.Status,
			"created_at": row.CreatedAt,
		})
	}

	return comments, total, nil
}

// UpdateCommentStatus 更新评论状态（屏蔽/恢复）
func (r *Repo) UpdateCommentStatus(commentID int64, status int) error {
	result := r.db.Table("comments").Where("id = ?", commentID).Update("status", status)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
