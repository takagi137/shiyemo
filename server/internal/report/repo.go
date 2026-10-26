package report

import (
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// ReportRepo 举报仓储
type ReportRepo struct {
	db *gorm.DB
}

// NewReportRepo 创建举报仓储
func NewReportRepo() *ReportRepo {
	return &ReportRepo{
		db: database.GetDB(),
	}
}

// NewReportRepoWithDB 使用注入的 DB 创建举报仓储（测试用）
func NewReportRepoWithDB(db *gorm.DB) *ReportRepo {
	return &ReportRepo{
		db: db,
	}
}

// Create 创建举报
func (r *ReportRepo) Create(report *Report) error {
	return r.db.Create(report).Error
}

// ExistsByUserAndTarget 检查用户是否已举报过该目标
func (r *ReportRepo) ExistsByUserAndTarget(userID int64, targetType string, targetID int64) (bool, error) {
	var count int64
	err := r.db.Model(&Report{}).Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Count(&count).Error
	return count > 0, err
}

// List 获取举报列表（管理员）
func (r *ReportRepo) List(page, size int) ([]Report, int64, error) {
	var reports []Report
	var total int64

	query := r.db.Model(&Report{})
	query.Count(&total)

	offset := (page - 1) * size
	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&reports).Error
	return reports, total, err
}

// UpdateStatus 更新举报状态
func (r *ReportRepo) UpdateStatus(id int64, status int) error {
	return r.db.Model(&Report{}).Where("id = ?", id).Update("status", status).Error
}
