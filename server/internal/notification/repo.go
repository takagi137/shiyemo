package notification

import (
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// RepoInterface 通知仓储接口（用于测试 mock）
type RepoInterface interface {
	Create(notification *Notification) error
	ListByUserID(userID int64, page, pageSize int) ([]Notification, int64, error)
	ListByType(userID int64, nType string, unreadOnly bool, page, pageSize int) ([]Notification, int64, error)
	BatchMarkAsRead(ids []int64, userID int64) error
	BatchDelete(ids []int64, userID int64) error
	GetByID(id int64) (*Notification, error)
	MarkAsRead(id int64, userID int64) error
	MarkAllAsRead(userID int64) error
	CountUnread(userID int64) (int64, error)
}

// Repo 通知仓储
type Repo struct {
	db *gorm.DB
}

// NewRepo 创建通知仓储
func NewRepo() *Repo {
	return &Repo{
		db: database.GetDB(),
	}
}

// NewRepoWithDB 使用注入的 DB 创建通知仓储（测试用）
func NewRepoWithDB(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// Create 创建通知
func (r *Repo) Create(notification *Notification) error {
	return r.db.Create(notification).Error
}

// ListByUserID 获取用户通知列表
func (r *Repo) ListByUserID(userID int64, page, pageSize int) ([]Notification, int64, error) {
	var notifications []Notification
	var total int64

	// 查询总数
	err := r.db.Model(&Notification{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	err = r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notifications).Error

	return notifications, total, err
}

// ListByType 按类型获取用户通知列表
func (r *Repo) ListByType(userID int64, nType string, unreadOnly bool, page, pageSize int) ([]Notification, int64, error) {
	var notifications []Notification
	var total int64

	query := r.db.Model(&Notification{}).Where("user_id = ?", userID)
	if nType != "" {
		query = query.Where("type = ?", nType)
	}
	if unreadOnly {
		query = query.Where("is_read = false")
	}

	// 查询总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notifications).Error

	return notifications, total, err
}

// BatchMarkAsRead 批量标记已读
func (r *Repo) BatchMarkAsRead(ids []int64, userID int64) error {
	return r.db.Model(&Notification{}).
		Where("id IN ? AND user_id = ? AND is_read = false", ids, userID).
		Update("is_read", true).Error
}

// BatchDelete 批量删除
func (r *Repo) BatchDelete(ids []int64, userID int64) error {
	return r.db.Where("id IN ? AND user_id = ?", ids, userID).
		Delete(&Notification{}).Error
}

// GetByID 根据 ID 获取通知
func (r *Repo) GetByID(id int64) (*Notification, error) {
	var notification Notification
	err := r.db.Where("id = ?", id).First(&notification).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// MarkAsRead 标记通知为已读
func (r *Repo) MarkAsRead(id int64, userID int64) error {
	return r.db.Model(&Notification{}).
		Where("id = ? AND user_id = ? AND is_read = false", id, userID).
		Update("is_read", true).Error
}

// MarkAllAsRead 标记用户所有通知为已读
func (r *Repo) MarkAllAsRead(userID int64) error {
	return r.db.Model(&Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

// CountUnread 获取未读通知数量
func (r *Repo) CountUnread(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}
