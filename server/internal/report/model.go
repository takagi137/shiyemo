package report

import (
	"time"
)

// Report 举报模型
type Report struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     int64     `json:"user_id" gorm:"index;not null"`
	TargetType string    `json:"target_type" gorm:"size:20;not null"` // post, comment
	TargetID   int64     `json:"target_id" gorm:"index;not null"`
	Reason     string    `json:"reason" gorm:"size:500;not null"`
	Status     int       `json:"status" gorm:"default:1"` // 1: 待处理, 2: 已处理, 3: 已驳回
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// 举报状态常量
const (
	ReportStatusPending   = 1 // 待处理
	ReportStatusProcessed = 2 // 已处理
	ReportStatusRejected  = 3 // 已驳回
)

// TableName 表名
func (Report) TableName() string {
	return "reports"
}
