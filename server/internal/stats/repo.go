package stats

import (
	"time"

	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// StatsRepoInterface 仓储接口
type StatsRepoInterface interface {
	// AggregateDailyStats 聚合每日统计数据
	AggregateDailyStats(date string) error
	// GetTrends 获取趋势数据
	GetTrends(days int) ([]DailyStat, error)
	// GetLatestStatDate 获取最新统计日期
	GetLatestStatDate() (string, error)
}

// StatsRepo 统计仓储
type StatsRepo struct {
	db *gorm.DB
}

// NewStatsRepo 创建统计仓储
func NewStatsRepo() *StatsRepo {
	return &StatsRepo{db: database.GetDB()}
}

// NewStatsRepoWithDB 使用指定 DB 创建（测试用）
func NewStatsRepoWithDB(db *gorm.DB) *StatsRepo {
	return &StatsRepo{db: db}
}

// AggregateDailyStats 聚合每日统计数据
func (r *StatsRepo) AggregateDailyStats(date string) error {
	// 按城市聚合
	type CityCount struct {
		City  string
		Count int64
	}
	var cityCounts []CityCount
	if err := r.db.Raw(`
		SELECT COALESCE(city, '') as city, COUNT(*) as count
		FROM users
		WHERE city != '' AND city IS NOT NULL
		GROUP BY city
	`).Scan(&cityCounts).Error; err != nil {
		return err
	}

	// 按行业聚合
	type IndustryCount struct {
		Industry string
		Count    int64
	}
	var industryCounts []IndustryCount
	if err := r.db.Raw(`
		SELECT COALESCE(industry, '') as industry, COUNT(*) as count
		FROM users
		WHERE industry != '' AND industry IS NOT NULL
		GROUP BY industry
	`).Scan(&industryCounts).Error; err != nil {
		return err
	}

	// 先删除当天旧数据
	r.db.Where("stat_date = ?", date).Delete(&DailyStat{})

	// 批量插入城市统计
	now := time.Now()
	for _, c := range cityCounts {
		stat := DailyStat{
			StatDate:  date,
			Dimension: "city",
			Label:     c.City,
			Count:     c.Count,
			CreatedAt: now,
		}
		if err := r.db.Create(&stat).Error; err != nil {
			return err
		}
	}

	// 批量插入行业统计
	for _, ind := range industryCounts {
		stat := DailyStat{
			StatDate:  date,
			Dimension: "industry",
			Label:     ind.Industry,
			Count:     ind.Count,
			CreatedAt: now,
		}
		if err := r.db.Create(&stat).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetTrends 获取趋势数据
func (r *StatsRepo) GetTrends(days int) ([]DailyStat, error) {
	var stats []DailyStat
	cutoff := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	if err := r.db.Where("stat_date >= ?", cutoff).
		Order("stat_date ASC, dimension ASC, label ASC").
		Find(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

// GetLatestStatDate 获取最新统计日期
func (r *StatsRepo) GetLatestStatDate() (string, error) {
	var stat DailyStat
	if err := r.db.Order("stat_date DESC").First(&stat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return stat.StatDate, nil
}
