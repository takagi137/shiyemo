package ambientsound

import (
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// RepoInterface 仓储接口（便于 service 层注入 mock 做单元测试）
type RepoInterface interface {
	ListActive() ([]AmbientSound, error)
	List(page, pageSize int, onlyActive bool) ([]AmbientSound, int64, error)
	GetByID(id int64) (*AmbientSound, error)
	Create(s *AmbientSound) error
	Update(s *AmbientSound) error
	Delete(id int64) error
	UpdateSortOrder(ids []int64) error
}

// Repo 数据仓库
type Repo struct {
	db *gorm.DB
}

// NewRepo 创建仓库
func NewRepo() *Repo {
	return &Repo{db: database.GetDB()}
}

// ListActive 获取所有启用音效（按排序，公开接口用）
func (r *Repo) ListActive() ([]AmbientSound, error) {
	var sounds []AmbientSound
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC, id ASC").Find(&sounds).Error
	return sounds, err
}

// List 获取音效列表（管理后台，分页）
func (r *Repo) List(page, pageSize int, onlyActive bool) ([]AmbientSound, int64, error) {
	var sounds []AmbientSound
	var total int64

	query := r.db.Model(&AmbientSound{})
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("sort_order ASC, id ASC").Offset(offset).Limit(pageSize).Find(&sounds).Error; err != nil {
		return nil, 0, err
	}

	return sounds, total, nil
}

// GetByID 根据 ID 获取
func (r *Repo) GetByID(id int64) (*AmbientSound, error) {
	var s AmbientSound
	err := r.db.First(&s, id).Error
	return &s, err
}

// Create 创建音效
func (r *Repo) Create(s *AmbientSound) error {
	return r.db.Create(s).Error
}

// Update 更新音效
func (r *Repo) Update(s *AmbientSound) error {
	return r.db.Save(s).Error
}

// Delete 删除音效（软删除）
func (r *Repo) Delete(id int64) error {
	return r.db.Delete(&AmbientSound{}, id).Error
}

// UpdateSortOrder 批量更新排序
func (r *Repo) UpdateSortOrder(ids []int64) error {
	for i, id := range ids {
		if err := r.db.Model(&AmbientSound{}).Where("id = ?", id).Update("sort_order", i+1).Error; err != nil {
			return err
		}
	}
	return nil
}
