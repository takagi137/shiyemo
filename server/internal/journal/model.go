package journal

import "time"

// Journal 私密日记
type Journal struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"not null;index"`
	Content   string    `json:"content" gorm:"not null;type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 表名
func (Journal) TableName() string {
	return "journals"
}
