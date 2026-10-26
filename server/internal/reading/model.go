package reading

import "time"

// Reading 经典阅读
type Reading struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Source        string    `json:"source" gorm:"size:200;not null"`        // 来源，如"苏轼《定风波》"
	Dynasty       string    `json:"dynasty" gorm:"size:50"`                 // 朝代
	Author        string    `json:"author" gorm:"size:100"`                 // 作者
	Title         string    `json:"title" gorm:"size:200"`                  // 标题
	Content       string    `json:"content" gorm:"not null;type:text"`      // 原文
	Interpretation string   `json:"interpretation" gorm:"type:text"`        // 背景解读
	Tags          string    `json:"tags" gorm:"size:500"`                   // 标签，逗号分隔
	ImageURL      string    `json:"image_url" gorm:"size:500"`              // 配图
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
}

// TableName 表名
func (Reading) TableName() string {
	return "readings"
}
