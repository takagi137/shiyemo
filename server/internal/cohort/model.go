package cohort

import "time"

// CohortStats 同温层统计数据
type CohortStats struct {
	City       string `json:"city"`
	Industry   string `json:"industry"`
	UserCount  int64  `json:"user_count"`
	PostCount  int64  `json:"post_count"`
	ActiveDays int    `json:"active_days"` // 近期活跃天数
}

// CohortPeer 同温层用户
type CohortPeer struct {
	ID         int64      `json:"id"`
	Nickname   string     `json:"nickname"`
	Avatar     string     `json:"avatar,omitempty"`
	Bio        string     `json:"bio,omitempty"`
	City       string     `json:"city,omitempty"`
	Industry   string     `json:"industry,omitempty"`
	LeftJobAt  *time.Time `json:"left_job_at,omitempty"`
	FreedomDays *int      `json:"freedom_days,omitempty"`
}

// CohortTimeRange 时间段筛选
type CohortTimeRange string

const (
	TimeRangeWeek  CohortTimeRange = "week"   // 本周
	TimeRangeMonth CohortTimeRange = "month"  // 本月
	TimeRangeYear  CohortTimeRange = "year"   // 本年
	TimeRangeAll   CohortTimeRange = "all"    // 全部
)

// CohortFilter 同温层筛选条件
type CohortFilter struct {
	City      string          `json:"city,omitempty"`
	Industry  string          `json:"industry,omitempty"`
	TimeRange CohortTimeRange `json:"time_range,omitempty"`
	Page      int             `json:"page"`
	PageSize  int             `json:"page_size"`
}
