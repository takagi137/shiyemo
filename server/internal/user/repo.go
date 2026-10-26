package user

import (
	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// UserRepoInterface 用户仓储接口（用于测试 mock）
type UserRepoInterface interface {
	GetByID(id int64) (*auth.User, error)
	Update(user *auth.User) error
	UpdateFields(id int64, fields map[string]interface{}) error
	GetUserPosts(userID int64, page, pageSize int) ([]map[string]interface{}, int64, error)
	CountUserPosts(userID int64) (int64, error)
	CountUserLikesReceived(userID int64) (int64, error)
}

// UserRepo 用户仓储
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建用户仓储
func NewUserRepo() *UserRepo {
	return &UserRepo{
		db: database.GetDB(),
	}
}

// NewUserRepoWithDB 使用注入的 DB 创建用户仓储（测试用）
func NewUserRepoWithDB(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// GetByID 根据 ID 查询用户
func (r *UserRepo) GetByID(id int64) (*auth.User, error) {
	var user auth.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *UserRepo) Update(user *auth.User) error {
	return r.db.Save(user).Error
}

// UpdateFields 更新指定字段
func (r *UserRepo) UpdateFields(id int64, fields map[string]interface{}) error {
	return r.db.Model(&auth.User{}).Where("id = ?", id).Updates(fields).Error
}

// GetUserPosts 获取用户帖子列表
func (r *UserRepo) GetUserPosts(userID int64, page, pageSize int) ([]map[string]interface{}, int64, error) {
	var total int64
	var posts []map[string]interface{}

	// 查询总数
	err := r.db.Table("posts").Where("user_id = ? AND status = 1 AND deleted_at IS NULL", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	err = r.db.Table("posts").
		Select("id, user_id, type, title, content, images, tags, like_count, comment_count, created_at").
		Where("user_id = ? AND status = 1 AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

// CountUserPosts 统计用户帖子数量
func (r *UserRepo) CountUserPosts(userID int64) (int64, error) {
	var count int64
	err := r.db.Table("posts").Where("user_id = ? AND status = 1 AND deleted_at IS NULL", userID).Count(&count).Error
	return count, err
}

// CountUserLikesReceived 统计用户获赞数量
func (r *UserRepo) CountUserLikesReceived(userID int64) (int64, error) {
	var count int64
	sql := `SELECT COALESCE(SUM(like_count), 0) FROM posts WHERE user_id = ? AND status = 1 AND deleted_at IS NULL`
	err := r.db.Raw(sql, userID).Scan(&count).Error
	return count, err
}
