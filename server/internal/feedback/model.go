package feedback

import "time"

// Feedback 用户反馈
type Feedback struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id" gorm:"index"`
	Type      string    `json:"type" gorm:"size:20;not null"` // bug, suggestion, other
	Title     string    `json:"title" gorm:"size:200;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Contact   string    `json:"contact" gorm:"size:100"` // 联系方式（可选）
	Status    int8      `json:"status" gorm:"default:0"` // 0-待处理 1-已处理 2-已忽略
	Reply     string    `json:"reply" gorm:"type:text"`   // 管理员回复
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FeedbackResponse 反馈响应（带用户信息）
type FeedbackResponse struct {
	Feedback
	UserNickname string `json:"user_nickname"`
	UserEmail    string `json:"user_email"`
}
