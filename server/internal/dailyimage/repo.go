package dailyimage

import (
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// RepoInterface 仓储接口（便于 service 层注入 mock 做单元测试）
type RepoInterface interface {
	RandomActive() (*DailyImage, error)
	List(page, pageSize int, onlyActive bool) ([]DailyImage, int64, error)
	GetByID(id int64) (*DailyImage, error)
	Create(img *DailyImage) error
	Update(img *DailyImage) error
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

// RandomActive 随机获取一张活跃图片
func (r *Repo) RandomActive() (*DailyImage, error) {
	var img DailyImage
	err := r.db.Where("is_active = ?", true).Order("RANDOM()").First(&img).Error
	return &img, err
}

// List 获取图片列表（管理后台）
func (r *Repo) List(page, pageSize int, onlyActive bool) ([]DailyImage, int64, error) {
	var images []DailyImage
	var total int64

	query := r.db.Model(&DailyImage{})
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("sort_order ASC, id ASC").Offset(offset).Limit(pageSize).Find(&images).Error; err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

// ListAll 获取所有活跃图片（按排序）
func (r *Repo) ListAll() ([]DailyImage, error) {
	var images []DailyImage
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC, id ASC").Find(&images).Error
	return images, err
}

// GetByID 根据 ID 获取
func (r *Repo) GetByID(id int64) (*DailyImage, error) {
	var img DailyImage
	err := r.db.First(&img, id).Error
	return &img, err
}

// Create 创建图片
func (r *Repo) Create(img *DailyImage) error {
	return r.db.Create(img).Error
}

// Update 更新图片
func (r *Repo) Update(img *DailyImage) error {
	return r.db.Save(img).Error
}

// Delete 删除图片（软删除）
func (r *Repo) Delete(id int64) error {
	return r.db.Delete(&DailyImage{}, id).Error
}

// UpdateSortOrder 批量更新排序
func (r *Repo) UpdateSortOrder(ids []int64) error {
	for i, id := range ids {
		if err := r.db.Model(&DailyImage{}).Where("id = ?", id).Update("sort_order", i+1).Error; err != nil {
			return err
		}
	}
	return nil
}

// CountActive 统计活跃图片数量
func (r *Repo) CountActive() int64 {
	var count int64
	r.db.Model(&DailyImage{}).Where("is_active = ?", true).Count(&count)
	return count
}
