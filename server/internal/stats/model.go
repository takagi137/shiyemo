package stats

import "time"

// TrendPoint 趋势数据点
type TrendPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// TrendSeries 趋势数据系列
type TrendSeries struct {
	Label  string       `json:"label"`
	Points []TrendPoint `json:"points"`
}

// TrendResponse 趋势响应
type TrendResponse struct {
	CityTrends     []TrendSeries `json:"city_trends"`
	IndustryTrends []TrendSeries `json:"industry_trends"`
	UpdatedAt      string        `json:"updated_at"`
}

// DailyStat 每日统计记录（数据库表）
type DailyStat struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	StatDate  string    `json:"stat_date" gorm:"size:10;not null;index"`
	Dimension string    `json:"dimension" gorm:"size:20;not null;index"` // city / industry
	Label     string    `json:"label" gorm:"size:100;not null;index"`
	Count     int64     `json:"count" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (DailyStat) TableName() string {
	return "daily_stats"
}
