package tips

import (
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// TipsRepoInterface 贴士仓储接口
type TipsRepoInterface interface {
	GetRandom() (*SurvivalTip, error)
	GetList(page, pageSize int) ([]SurvivalTip, int64, error)
	Create(tip *SurvivalTip) error
	Update(tip *SurvivalTip) error
	Delete(id int64) error
}

// TipsRepo 贴士仓储实现
type TipsRepo struct {
	db *gorm.DB
}

// NewTipsRepo 创建贴士仓储
func NewTipsRepo() *TipsRepo {
	return &TipsRepo{
		db: database.GetDB(),
	}
}

// NewTipsRepoWithDB 使用注入的 DB 创建贴士仓储（测试用）
func NewTipsRepoWithDB(db *gorm.DB) *TipsRepo {
	return &TipsRepo{
		db: db,
	}
}

// GetRandom 随机获取一条贴士
func (r *TipsRepo) GetRandom() (*SurvivalTip, error) {
	var tip SurvivalTip
	err := r.db.Where("is_active = ?", true).Order("RANDOM()").First(&tip).Error
	if err != nil {
		return nil, err
	}
	return &tip, nil
}

// GetList 获取贴士列表（分页）
func (r *TipsRepo) GetList(page, pageSize int) ([]SurvivalTip, int64, error) {
	var tips []SurvivalTip
	var total int64

	// 统计总数
	if err := r.db.Model(&SurvivalTip{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := r.db.Where("is_active = ?", true).Order("id DESC").Offset(offset).Limit(pageSize).Find(&tips).Error; err != nil {
		return nil, 0, err
	}

	return tips, total, nil
}

// Create 创建贴士
func (r *TipsRepo) Create(tip *SurvivalTip) error {
	return r.db.Create(tip).Error
}

// Update 更新贴士
func (r *TipsRepo) Update(tip *SurvivalTip) error {
	return r.db.Save(tip).Error
}

// Delete 删除贴士（软删除）
func (r *TipsRepo) Delete(id int64) error {
	return r.db.Model(&SurvivalTip{}).Where("id = ?", id).Update("is_active", false).Error
}
