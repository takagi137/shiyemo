package tips

import "time"

// SurvivalTip 生存小贴士
type SurvivalTip struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Content   string    `json:"content" gorm:"not null;type:text"`
	Category  string    `json:"category" gorm:"size:50;default:'general'"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (SurvivalTip) TableName() string {
	return "survival_tips"
}
