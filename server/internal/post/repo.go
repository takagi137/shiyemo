package post

import (
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// PostRepoInterface 帖子仓储接口（用于测试 mock）
type PostRepoInterface interface {
	Create(post *Post) error
	GetByID(id int64) (*Post, error)
	List(postType string, page, size int) ([]Post, int64, error)
	ListByCohort(postType, city, industry string, page, size int) ([]Post, int64, error)
	UpdateFields(id int64, fields map[string]interface{}) error
	Delete(id int64) error
	Count() (int64, error)
	GetByUserIDs(userIDs []int64, page, size int) ([]Post, int64, error)
	GetFollowingIDs(userID int64) ([]int64, error)
}

// PostRepo 帖子仓储
type PostRepo struct {
	db *gorm.DB
}

// NewPostRepo 创建帖子仓储
func NewPostRepo() *PostRepo {
	return &PostRepo{
		db: database.GetDB(),
	}
}

// NewPostRepoWithDB 使用注入的 DB 创建帖子仓储（测试用）
func NewPostRepoWithDB(db *gorm.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

// Create 创建帖子
func (r *PostRepo) Create(post *Post) error {
	return r.db.Create(post).Error
}

// GetByID 根据 ID 查询帖子
func (r *PostRepo) GetByID(id int64) (*Post, error) {
	var post Post
	err := r.db.Where("id = ? AND status = 1", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// List 获取帖子列表
func (r *PostRepo) List(postType string, page, size int) ([]Post, int64, error) {
	var posts []Post
	var total int64

	query := r.db.Model(&Post{}).Where("status = 1")

	if postType != "" {
		query = query.Where("type = ?", postType)
	}

	// 计算总数
	query.Count(&total)

	// 分页查询（置顶优先，然后按时间倒序）
	offset := (page - 1) * size
	err := query.Order("is_top DESC, is_featured DESC, created_at DESC").Offset(offset).Limit(size).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// Update 更新帖子
func (r *PostRepo) Update(post *Post) error {
	return r.db.Save(post).Error
}

// UpdateFields 更新指定字段
func (r *PostRepo) UpdateFields(id int64, fields map[string]interface{}) error {
	return r.db.Model(&Post{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 软删除帖子
func (r *PostRepo) Delete(id int64) error {
	return r.db.Model(&Post{}).Where("id = ?", id).Update("status", 3).Error
}

// Count 统计帖子数量
func (r *PostRepo) Count() (int64, error) {
	var count int64
	err := r.db.Model(&Post{}).Where("status = 1").Count(&count).Error
	return count, err
}

// GetByUserIDs 根据用户 ID 列表获取帖子
func (r *PostRepo) GetByUserIDs(userIDs []int64, page, size int) ([]Post, int64, error) {
	var posts []Post
	var total int64

	if len(userIDs) == 0 {
		return posts, 0, nil
	}

	query := r.db.Model(&Post{}).Where("status = 1 AND user_id IN ?", userIDs)

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * size
	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetFollowingIDs 获取用户关注的用户 ID 列表
func (r *PostRepo) GetFollowingIDs(userID int64) ([]int64, error) {
	var followingIDs []int64
	err := r.db.Table("follows").Where("follower_id = ?", userID).Pluck("following_id", &followingIDs).Error
	return followingIDs, err
}

// ListByCohort 根据同温层条件获取帖子列表
func (r *PostRepo) ListByCohort(postType, city, industry string, page, size int) ([]Post, int64, error) {
	var posts []Post
	var total int64

	// 构建查询，关联用户表过滤城市和行业
	query := r.db.Table("posts p").
		Joins("JOIN users u ON p.user_id = u.id").
		Where("p.status = 1 AND u.status = 1")

	if postType != "" {
		query = query.Where("p.type = ?", postType)
	}
	if city != "" {
		query = query.Where("u.city = ?", city)
	}
	if industry != "" {
		query = query.Where("u.industry = ?", industry)
	}

	// 计算总数
	query.Count(&total)

	// 分页查询（置顶优先，然后按时间倒序）
	offset := (page - 1) * size
	err := query.Select("p.*").
		Order("p.is_top DESC, p.is_featured DESC, p.created_at DESC").
		Offset(offset).Limit(size).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
