package journal

import (
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// JournalRepoInterface 日记仓储接口
type JournalRepoInterface interface {
	Create(journal *Journal) error
	GetByUserID(userID int64, page, pageSize int) ([]Journal, int64, error)
	GetByID(id, userID int64) (*Journal, error)
	Delete(id, userID int64) error
}

// JournalRepo 日记仓储实现
type JournalRepo struct {
	db *gorm.DB
}

// NewJournalRepo 创建日记仓储
func NewJournalRepo() *JournalRepo {
	return &JournalRepo{
		db: database.GetDB(),
	}
}

// NewJournalRepoWithDB 使用注入的 DB 创建日记仓储（测试用）
func NewJournalRepoWithDB(db *gorm.DB) *JournalRepo {
	return &JournalRepo{
		db: db,
	}
}

// Create 创建日记
func (r *JournalRepo) Create(journal *Journal) error {
	return r.db.Create(journal).Error
}

// GetByUserID 获取用户日记列表（分页，按时间倒序）
func (r *JournalRepo) GetByUserID(userID int64, page, pageSize int) ([]Journal, int64, error) {
	var journals []Journal
	var total int64

	if err := r.db.Model(&Journal{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&journals).Error; err != nil {
		return nil, 0, err
	}

	return journals, total, nil
}

// GetByID 根据 ID 获取日记（校验用户归属）
func (r *JournalRepo) GetByID(id, userID int64) (*Journal, error) {
	var journal Journal
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&journal).Error
	if err != nil {
		return nil, err
	}
	return &journal, nil
}

// Delete 删除日记（只能删除自己的）
func (r *JournalRepo) Delete(id, userID int64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&Journal{}).Error
}
