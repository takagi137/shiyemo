package checkin

import (
	"time"
)

// DailyCheckin 每日打卡记录
type DailyCheckin struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     int64     `json:"user_id" gorm:"not null;index"`
	Mood       string    `json:"mood" gorm:"size:20;not null"` // happy, calm, sad, angry, anxious
	Content    string    `json:"content" gorm:"type:text"`
	Images     string    `json:"images" gorm:"type:text"` // JSON 数组，存储图片 URL
	IsPublic   bool      `json:"is_public" gorm:"default:true"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName 表名
func (DailyCheckin) TableName() string {
	return "daily_checkins"
}

// MoodType 心情类型常量
const (
	MoodHappy   = "happy"
	MoodCalm    = "calm"
	MoodSad     = "sad"
	MoodAngry   = "angry"
	MoodAnxious = "anxious"
)

// MoodEmoji 心情对应的 emoji
var MoodEmoji = map[string]string{
	MoodHappy:   "😊",
	MoodCalm:    "😌",
	MoodSad:     "😢",
	MoodAngry:   "😠",
	MoodAnxious: "😰",
}

// MoodLabel 心情对应的中文标签
var MoodLabel = map[string]string{
	MoodHappy:   "开心",
	MoodCalm:    "平静",
	MoodSad:     "难过",
	MoodAngry:   "生气",
	MoodAnxious: "焦虑",
}
