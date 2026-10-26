package search

import (
	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/post"
	"gorm.io/gorm"
)

// SearchRepoInterface 搜索仓储接口
type SearchRepoInterface interface {
	GetActivePosts() ([]post.Post, error)
	GetActiveUsers() ([]auth.User, error)
}

// SearchRepo 搜索仓储实现
type SearchRepo struct {
	db *gorm.DB
}

// NewSearchRepo 创建搜索仓储
func NewSearchRepo() *SearchRepo {
	return &SearchRepo{}
}

// NewSearchRepoWithDB 使用注入的 DB 创建搜索仓储（测试用）
func NewSearchRepoWithDB(db *gorm.DB) *SearchRepo {
	return &SearchRepo{db: db}
}

// GetActivePosts 获取所有已发布帖子
func (r *SearchRepo) GetActivePosts() ([]post.Post, error) {
	var posts []post.Post
	err := r.db.Where("status = 1").Find(&posts).Error
	return posts, err
}

// GetActiveUsers 获取所有正常用户
func (r *SearchRepo) GetActiveUsers() ([]auth.User, error) {
	var users []auth.User
	err := r.db.Where("status = 1").Find(&users).Error
	return users, err
}
