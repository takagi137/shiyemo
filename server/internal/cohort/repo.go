package cohort

import (
	"time"

	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/database"
	"gorm.io/gorm"
)

// CohortRepoInterface 同温层仓储接口
type CohortRepoInterface interface {
	// GetCohortStats 获取同温层统计数据
	GetCohortStats(city, industry string, timeRange CohortTimeRange) (*CohortStats, error)
	// ListPeers 获取同温层用户列表
	ListPeers(filter CohortFilter) ([]CohortPeer, int64, error)
	// GetMyCohortStats 获取当前用户的同温层统计
	GetMyCohortStats(userID int64) (*CohortStats, error)
	// GetCities 获取所有城市列表
	GetCities() ([]string, error)
	// GetIndustries 获取所有行业列表
	GetIndustries() ([]string, error)
}

// CohortRepo 同温层仓储实现
type CohortRepo struct {
	db *gorm.DB
}

// NewCohortRepo 创建同温层仓储
func NewCohortRepo() *CohortRepo {
	return &CohortRepo{
		db: database.GetDB(),
	}
}

// NewCohortRepoWithDB 使用指定 DB 创建仓储（测试用）
func NewCohortRepoWithDB(db *gorm.DB) *CohortRepo {
	return &CohortRepo{db: db}
}

// GetCohortStats 获取同温层统计数据
func (r *CohortRepo) GetCohortStats(city, industry string, timeRange CohortTimeRange) (*CohortStats, error) {
	var stats CohortStats
	stats.City = city
	stats.Industry = industry

	// 构建查询条件
	query := r.db.Model(&auth.User{}).Where("status = 1")

	if city != "" {
		query = query.Where("city = ?", city)
	}
	if industry != "" {
		query = query.Where("industry = ?", industry)
	}

	// 时间范围筛选
	if timeRange != "" && timeRange != TimeRangeAll {
		var startTime time.Time
		now := time.Now()
		switch timeRange {
		case TimeRangeWeek:
			startTime = now.AddDate(0, 0, -7)
		case TimeRangeMonth:
			startTime = now.AddDate(0, -1, 0)
		case TimeRangeYear:
			startTime = now.AddDate(-1, 0, 0)
		}
		query = query.Where("left_job_at >= ?", startTime)
	}

	// 统计用户数
	if err := query.Count(&stats.UserCount).Error; err != nil {
		return nil, err
	}

	// 统计帖子数（需要关联帖子表）
	var postCount int64
	postQuery := r.db.Table("posts p").
		Joins("JOIN users u ON p.user_id = u.id").
		Where("u.status = 1 AND p.status = 1")

	if city != "" {
		postQuery = postQuery.Where("u.city = ?", city)
	}
	if industry != "" {
		postQuery = postQuery.Where("u.industry = ?", industry)
	}

	if err := postQuery.Count(&postCount).Error; err != nil {
		return nil, err
	}
	stats.PostCount = postCount

	return &stats, nil
}

// ListPeers 获取同温层用户列表
func (r *CohortRepo) ListPeers(filter CohortFilter) ([]CohortPeer, int64, error) {
	var users []auth.User
	var total int64

	// 构建查询条件
	query := r.db.Model(&auth.User{}).Where("status = 1")

	if filter.City != "" {
		query = query.Where("city = ?", filter.City)
	}
	if filter.Industry != "" {
		query = query.Where("industry = ?", filter.Industry)
	}

	// 时间范围筛选
	if filter.TimeRange != "" && filter.TimeRange != TimeRangeAll {
		var startTime time.Time
		now := time.Now()
		switch filter.TimeRange {
		case TimeRangeWeek:
			startTime = now.AddDate(0, 0, -7)
		case TimeRangeMonth:
			startTime = now.AddDate(0, -1, 0)
		case TimeRangeYear:
			startTime = now.AddDate(-1, 0, 0)
		}
		query = query.Where("left_job_at >= ?", startTime)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 50 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(offset).Limit(filter.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 CohortPeer
	peers := make([]CohortPeer, 0, len(users))
	for _, u := range users {
		peer := CohortPeer{
			ID:        u.ID,
			Nickname:  u.Nickname,
			Avatar:    u.Avatar,
			Bio:       u.Bio,
			City:      u.City,
			Industry:  u.Industry,
			LeftJobAt: u.LeftJobAt,
		}
		peer.FreedomDays = u.FreedomDays()
		peers = append(peers, peer)
	}

	return peers, total, nil
}

// GetMyCohortStats 获取当前用户的同温层统计
func (r *CohortRepo) GetMyCohortStats(userID int64) (*CohortStats, error) {
	// 先获取用户信息
	var user auth.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// 如果用户没有设置城市或行业，返回空统计
	if user.City == "" && user.Industry == "" {
		return &CohortStats{
			City:      user.City,
			Industry:  user.Industry,
			UserCount: 0,
			PostCount: 0,
		}, nil
	}

	// 获取同温层统计
	return r.GetCohortStats(user.City, user.Industry, TimeRangeAll)
}

// GetCities 获取所有城市列表
func (r *CohortRepo) GetCities() ([]string, error) {
	var cities []string
	if err := r.db.Model(&auth.User{}).
		Where("status = 1 AND city != ''").
		Distinct("city").
		Pluck("city", &cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}

// GetIndustries 获取所有行业列表
func (r *CohortRepo) GetIndustries() ([]string, error) {
	var industries []string
	if err := r.db.Model(&auth.User{}).
		Where("status = 1 AND industry != ''").
		Distinct("industry").
		Pluck("industry", &industries).Error; err != nil {
		return nil, err
	}
	return industries, nil
}
