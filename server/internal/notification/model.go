package notification

import "time"

// Notification 通知模型
type Notification struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	UserID     int64     `json:"user_id" gorm:"not null;index"`
	ActorID    int64     `json:"actor_id" gorm:"not null"`
	Type       string    `json:"type" gorm:"size:20;not null"`
	TargetID   *int64    `json:"target_id"`
	TargetType *string   `json:"target_type" gorm:"size:20"`
	Content    string    `json:"content" gorm:"type:text"`
	IsRead     bool      `json:"is_read" gorm:"default:false;index:idx_notifications_user_read"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 表名
func (Notification) TableName() string {
	return "notifications"
}

// NotificationType 通知类型常量
const (
	TypeLike    = "like"
	TypeComment = "comment"
	TypeReply   = "reply"
	TypeFollow  = "follow"
)
