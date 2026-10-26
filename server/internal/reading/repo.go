package reading

import (
	"time"

	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// ReadingRepoInterface 阅读仓储接口
type ReadingRepoInterface interface {
	GetByID(id int64) (*Reading, error)
	GetToday() (*Reading, error)
	GetRandom(excludeIDs []int64) (*Reading, error)
	Count() (int64, error)
}

// ReadingRepo 阅读仓储实现
type ReadingRepo struct {
	db *gorm.DB
}

// NewReadingRepo 创建阅读仓储
func NewReadingRepo() *ReadingRepo {
	return &ReadingRepo{
		db: database.GetDB(),
	}
}

// NewReadingRepoWithDB 使用注入的 DB 创建阅读仓储（测试用）
func NewReadingRepoWithDB(db *gorm.DB) *ReadingRepo {
	return &ReadingRepo{
		db: db,
	}
}

// GetByID 根据 ID 获取阅读
func (r *ReadingRepo) GetByID(id int64) (*Reading, error) {
	var reading Reading
	err := r.db.Where("id = ? AND is_active = ?", id, true).First(&reading).Error
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

// GetToday 获取今日阅读（按日期轮换）
func (r *ReadingRepo) GetToday() (*Reading, error) {
	var reading Reading
	// 用日期的 day-of-year 模总数实现轮换
oday := time.Now().YearDay()
	var total int64
	if err := r.db.Model(&Reading{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, err
	}
	if total == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	offset := int64(oday-1) % total
	err := r.db.Where("is_active = ?", true).Order("id ASC").Offset(int(offset)).Limit(1).First(&reading).Error
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

// GetRandom 随机获取一篇（排除已读）
func (r *ReadingRepo) GetRandom(excludeIDs []int64) (*Reading, error) {
	var reading Reading
	query := r.db.Where("is_active = ?", true)
	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}
	err := query.Order("RANDOM()").First(&reading).Error
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

// Count 统计有效阅读数
func (r *ReadingRepo) Count() (int64, error) {
	var total int64
	err := r.db.Model(&Reading{}).Where("is_active = ?", true).Count(&total).Error
	return total, err
}
