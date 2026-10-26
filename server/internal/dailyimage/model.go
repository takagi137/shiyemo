package dailyimage

import "time"

// DailyImage 每日风景图
type DailyImage struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	URL         string    `json:"url" gorm:"size:500;not null"`
	Description string    `json:"description" gorm:"size:200;default:''"`
	Credit      string    `json:"credit" gorm:"size:200;default:''"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
