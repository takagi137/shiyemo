package cohort

import (
	"context"
	"log"
	"time"

	"github.com/shiyemo/server/internal/common/cache"
)

// CohortService 同温层服务
type CohortService struct {
	repo CohortRepoInterface
	sf   *cache.Singleflight
}

// NewCohortService 创建同温层服务
func NewCohortService() *CohortService {
	return &CohortService{
		repo: NewCohortRepo(),
		sf:   cache.NewSingleflight(),
	}
}

// NewCohortServiceWithRepo 使用注入的 repo 创建服务（测试用）
func NewCohortServiceWithRepo(repo CohortRepoInterface) *CohortService {
	return &CohortService{
		repo: repo,
		sf:   cache.NewSingleflight(),
	}
}

// GetMyStats 获取当前用户的同温层统计
func (s *CohortService) GetMyStats(userID int64) (*CohortStats, error) {
	ctx := context.Background()
	cacheKey := cache.FormatUserDetailKey(userID) + ":cohort"

	// 尝试从缓存获取
	var cachedStats CohortStats
	if err := cache.Get(ctx, cacheKey, &cachedStats); err == nil {
		return &cachedStats, nil
	}

	// 缓存未命中，使用 singleflight
	v, err := s.sf.Do(cacheKey, func() (interface{}, error) {
		return s.repo.GetMyCohortStats(userID)
	})
	if err != nil {
		return nil, err
	}
	stats := v.(*CohortStats)

	// 写入缓存（5 分钟）
	if err := cache.Set(ctx, cacheKey, stats, 5*time.Minute); err != nil {
		log.Printf("Failed to cache cohort stats for user %d: %v", userID, err)
	}

	return stats, nil
}

// GetStats 获取指定同温层统计
func (s *CohortService) GetStats(city, industry string, timeRange CohortTimeRange) (*CohortStats, error) {
	ctx := context.Background()
	cacheKey := "cohort:stats:" + city + ":" + industry + ":" + string(timeRange)

	// 尝试从缓存获取
	var cachedStats CohortStats
	if err := cache.Get(ctx, cacheKey, &cachedStats); err == nil {
		return &cachedStats, nil
	}

	// 缓存未命中，使用 singleflight
	v, err := s.sf.Do(cacheKey, func() (interface{}, error) {
		return s.repo.GetCohortStats(city, industry, timeRange)
	})
	if err != nil {
		return nil, err
	}
	stats := v.(*CohortStats)

	// 写入缓存（5 分钟）
	if err := cache.Set(ctx, cacheKey, stats, 5*time.Minute); err != nil {
		log.Printf("Failed to cache cohort stats: %v", err)
	}

	return stats, nil
}

// ListPeers 获取同温层用户列表
func (s *CohortService) ListPeers(filter CohortFilter) ([]CohortPeer, int64, error) {
	// 参数校验
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 50 {
		filter.PageSize = 20
	}

	return s.repo.ListPeers(filter)
}

// GetCities 获取所有城市列表
func (s *CohortService) GetCities() ([]string, error) {
	ctx := context.Background()
	cacheKey := "cohort:cities"

	// 尝试从缓存获取
	var cachedCities []string
	if err := cache.Get(ctx, cacheKey, &cachedCities); err == nil {
		return cachedCities, nil
	}

	cities, err := s.repo.GetCities()
	if err != nil {
		return nil, err
	}

	// 写入缓存（10 分钟）
	if err := cache.Set(ctx, cacheKey, cities, 10*time.Minute); err != nil {
		log.Printf("Failed to cache cities: %v", err)
	}

	return cities, nil
}

// GetIndustries 获取所有行业列表
func (s *CohortService) GetIndustries() ([]string, error) {
	ctx := context.Background()
	cacheKey := "cohort:industries"

	// 尝试从缓存获取
	var cachedIndustries []string
	if err := cache.Get(ctx, cacheKey, &cachedIndustries); err == nil {
		return cachedIndustries, nil
	}

	industries, err := s.repo.GetIndustries()
	if err != nil {
		return nil, err
	}

	// 写入缓存（10 分钟）
	if err := cache.Set(ctx, cacheKey, industries, 10*time.Minute); err != nil {
		log.Printf("Failed to cache industries: %v", err)
	}

	return industries, nil
}
