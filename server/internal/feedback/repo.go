package feedback

import (
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// Repo 数据仓库
type Repo struct {
	db *gorm.DB
}

// NewRepo 创建仓库
func NewRepo() *Repo {
	return &Repo{db: database.GetDB()}
}

// Create 创建反馈
func (r *Repo) Create(feedback *Feedback) error {
	return r.db.Create(feedback).Error
}

// GetByID 根据 ID 获取反馈
func (r *Repo) GetByID(id int64) (*Feedback, error) {
	var feedback Feedback
	err := r.db.First(&feedback, id).Error
	return &feedback, err
}

// List 获取反馈列表（管理员）
func (r *Repo) List(page, pageSize int, status *int, feedbackType string) ([]FeedbackResponse, int64, error) {
	var feedbacks []FeedbackResponse
	var total int64

	query := r.db.Model(&Feedback{}).
		Joins("LEFT JOIN users ON users.id = feedbacks.user_id").
		Select("feedbacks.*, users.nickname as user_nickname, users.email as user_email")

	if status != nil {
		query = query.Where("feedbacks.status = ?", *status)
	}
	if feedbackType != "" {
		query = query.Where("feedbacks.type = ?", feedbackType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("feedbacks.created_at DESC").Offset(offset).Limit(pageSize).Scan(&feedbacks).Error; err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

// ListByUser 获取用户的反馈列表
func (r *Repo) ListByUser(userID int64, page, pageSize int) ([]Feedback, int64, error) {
	var feedbacks []Feedback
	var total int64

	query := r.db.Model(&Feedback{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&feedbacks).Error; err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

// UpdateStatus 更新反馈状态
func (r *Repo) UpdateStatus(id int64, status int8, reply string) error {
	return r.db.Model(&Feedback{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": status,
		"reply":  reply,
	}).Error
}
