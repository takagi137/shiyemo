package stats

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shiyemo/server/internal/common/cache"
)

// StatsServiceInterface 服务接口
type StatsServiceInterface interface {
	// RunDailyAggregation 执行每日统计聚合
	RunDailyAggregation() error
	// GetTrends 获取趋势数据（带缓存）
	GetTrends(days int) (*TrendResponse, error)
}

// StatsService 统计服务
type StatsService struct {
	repo StatsRepoInterface
	sf   *cache.Singleflight
}

// NewStatsService 创建统计服务
func NewStatsService() *StatsService {
	return &StatsService{
		repo: NewStatsRepo(),
		sf:   cache.NewSingleflight(),
	}
}

// NewStatsServiceWithRepo 使用注入的 repo 创建服务（测试用）
func NewStatsServiceWithRepo(repo StatsRepoInterface) *StatsService {
	return &StatsService{
		repo: repo,
		sf:   cache.NewSingleflight(),
	}
}

// RunDailyAggregation 执行每日统计聚合
func (s *StatsService) RunDailyAggregation() error {
	today := time.Now().Format("2006-01-02")
	log.Printf("Running daily stats aggregation for %s", today)

	if err := s.repo.AggregateDailyStats(today); err != nil {
		log.Printf("Daily stats aggregation failed: %v", err)
		return err
	}

	// 清除缓存
	ctx := context.Background()
	cache.Delete(ctx, "stats:trends:7")
	cache.Delete(ctx, "stats:trends:30")
	cache.Delete(ctx, "stats:trends:90")

	log.Printf("Daily stats aggregation completed for %s", today)
	return nil
}

// GetTrends 获取趋势数据（带缓存）
func (s *StatsService) GetTrends(days int) (*TrendResponse, error) {
	if days < 1 {
		days = 7
	}
	if days > 90 {
		days = 90
	}

	ctx := context.Background()
	cacheKey := "stats:trends:" + fmt.Sprintf("%d", days)

	// 尝试从缓存获取
	var cached TrendResponse
	if err := cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	// 使用 singleflight
	v, err := s.sf.Do(cacheKey, func() (interface{}, error) {
		return s.buildTrendResponse(days)
	})
	if err != nil {
		return nil, err
	}
	resp := v.(*TrendResponse)

	// 写入缓存（10 分钟）
	if err := cache.Set(ctx, cacheKey, resp, 10*time.Minute); err != nil {
		log.Printf("Failed to cache trends: %v", err)
	}

	return resp, nil
}

// buildTrendResponse 构建趋势响应
func (s *StatsService) buildTrendResponse(days int) (*TrendResponse, error) {
	stats, err := s.repo.GetTrends(days)
	if err != nil {
		return nil, err
	}

	// 按维度分组
	cityMap := make(map[string][]TrendPoint)
	industryMap := make(map[string][]TrendPoint)

	for _, stat := range stats {
		point := TrendPoint{Date: stat.StatDate, Count: stat.Count}
		switch stat.Dimension {
		case "city":
			cityMap[stat.Label] = append(cityMap[stat.Label], point)
		case "industry":
			industryMap[stat.Label] = append(industryMap[stat.Label], point)
		}
	}

	// 取 Top 5 城市和行业
	cityTrends := topNSeries(cityMap, 5)
	industryTrends := topNSeries(industryMap, 5)

	// 获取最新统计日期
	latestDate, _ := s.repo.GetLatestStatDate()
	if latestDate == "" {
		latestDate = time.Now().Format("2006-01-02")
	}

	return &TrendResponse{
		CityTrends:     cityTrends,
		IndustryTrends: industryTrends,
		UpdatedAt:      latestDate,
	}, nil
}

// topNSeries 取最新数据点的 Count 排序，取前 N 个系列
func topNSeries(m map[string][]TrendPoint, n int) []TrendSeries {
	type kv struct {
		label string
		series []TrendPoint
		latest int64
	}
	var items []kv
	for label, points := range m {
		var latest int64
		if len(points) > 0 {
			latest = points[len(points)-1].Count
		}
		items = append(items, kv{label: label, series: points, latest: latest})
	}
	// 简单排序（冒泡，N 小）
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].latest > items[i].latest {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	if len(items) > n {
		items = items[:n]
	}
	result := make([]TrendSeries, 0, len(items))
	for _, item := range items {
		result = append(result, TrendSeries{Label: item.label, Points: item.series})
	}
	return result
}
