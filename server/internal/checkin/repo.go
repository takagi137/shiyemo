package checkin

import (
	"time"

	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// CheckinRepoInterface 打卡仓储接口
type CheckinRepoInterface interface {
	Create(checkin *DailyCheckin) error
	GetTodayByUserID(userID int64) (*DailyCheckin, error)
	GetHistoryByUserID(userID int64, page, pageSize int) ([]DailyCheckin, int64, error)
	GetMoodCurve(userID int64, days int) ([]DailyCheckin, error)
	GetPublicTimeline(page, pageSize int) ([]DailyCheckin, int64, error)
}

// CheckinRepo 打卡仓储实现
type CheckinRepo struct {
	db *gorm.DB
}

// NewCheckinRepo 创建打卡仓储
func NewCheckinRepo() *CheckinRepo {
	return &CheckinRepo{
		db: database.GetDB(),
	}
}

// NewCheckinRepoWithDB 使用注入的 DB 创建打卡仓储（测试用）
func NewCheckinRepoWithDB(db *gorm.DB) *CheckinRepo {
	return &CheckinRepo{
		db: db,
	}
}

// Create 创建打卡记录
func (r *CheckinRepo) Create(checkin *DailyCheckin) error {
	return r.db.Create(checkin).Error
}

// GetTodayByUserID 获取用户今天的打卡记录
func (r *CheckinRepo) GetTodayByUserID(userID int64) (*DailyCheckin, error) {
	var checkin DailyCheckin
	today := time.Now().Format("2006-01-02")
	err := r.db.Where("user_id = ? AND DATE(created_at) = ?", userID, today).First(&checkin).Error
	if err != nil {
		return nil, err
	}
	return &checkin, nil
}

// GetHistoryByUserID 获取用户打卡历史（分页）
func (r *CheckinRepo) GetHistoryByUserID(userID int64, page, pageSize int) ([]DailyCheckin, int64, error) {
	var checkins []DailyCheckin
	var total int64

	// 统计总数
	if err := r.db.Model(&DailyCheckin{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&checkins).Error; err != nil {
		return nil, 0, err
	}

	return checkins, total, nil
}

// GetMoodCurve 获取用户近 N 天的心情数据
func (r *CheckinRepo) GetMoodCurve(userID int64, days int) ([]DailyCheckin, error) {
	var checkins []DailyCheckin
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	err := r.db.Where("user_id = ? AND DATE(created_at) >= ?", userID, startDate).
		Order("created_at ASC").
		Find(&checkins).Error
	if err != nil {
		return nil, err
	}

	return checkins, nil
}

// GetPublicTimeline 获取公开打卡时间线（分页）
func (r *CheckinRepo) GetPublicTimeline(page, pageSize int) ([]DailyCheckin, int64, error) {
	var checkins []DailyCheckin
	var total int64

	// 统计总数
	if err := r.db.Model(&DailyCheckin{}).Where("is_public = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := r.db.Where("is_public = ?", true).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&checkins).Error; err != nil {
		return nil, 0, err
	}

	return checkins, total, nil
}
