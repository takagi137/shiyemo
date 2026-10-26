package ambientsound

import "time"

// AmbientSound 静思空间环境音
// 与 dailyimage 一致：后端维护列表，url 存外部下载地址（不落库二进制）。
type AmbientSound struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"size:50;not null"`   // 稳定标识：rain/piano/forest...
	Name      string    `json:"name" gorm:"size:50;not null"`  // 显示名：雨声
	Emoji     string    `json:"emoji" gorm:"size:16;default:''"`
	URL       string    `json:"url" gorm:"size:500;not null"`  // 音频下载地址
	Credit    string    `json:"credit" gorm:"size:200;default:''"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 表名
func (AmbientSound) TableName() string {
	return "ambient_sounds"
}
