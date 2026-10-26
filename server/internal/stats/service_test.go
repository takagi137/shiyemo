package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockStatsRepo 模拟统计仓储
type mockStatsRepo struct {
	stats       []DailyStat
	aggregateErr error
	trendsErr    error
	latestDate   string
}

func (m *mockStatsRepo) AggregateDailyStats(date string) error {
	return m.aggregateErr
}

func (m *mockStatsRepo) GetTrends(days int) ([]DailyStat, error) {
	if m.trendsErr != nil {
		return nil, m.trendsErr
	}
	return m.stats, nil
}

func (m *mockStatsRepo) GetLatestStatDate() (string, error) {
	return m.latestDate, nil
}

func TestStatsService_RunDailyAggregation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockStatsRepo{}
		service := NewStatsServiceWithRepo(repo)

		err := service.RunDailyAggregation()
		require.NoError(t, err)
	})

	t.Run("repo_error", func(t *testing.T) {
		repo := &mockStatsRepo{aggregateErr: assert.AnError}
		service := NewStatsServiceWithRepo(repo)

		err := service.RunDailyAggregation()
		require.Error(t, err)
	})
}

func TestStatsService_GetTrends(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockStatsRepo{
			stats: []DailyStat{
				{StatDate: "2026-05-25", Dimension: "city", Label: "北京", Count: 10},
				{StatDate: "2026-05-25", Dimension: "city", Label: "上海", Count: 8},
				{StatDate: "2026-05-25", Dimension: "industry", Label: "互联网/IT", Count: 15},
				{StatDate: "2026-05-26", Dimension: "city", Label: "北京", Count: 12},
				{StatDate: "2026-05-26", Dimension: "city", Label: "上海", Count: 9},
				{StatDate: "2026-05-26", Dimension: "industry", Label: "互联网/IT", Count: 18},
				{StatDate: "2026-05-26", Dimension: "industry", Label: "金融", Count: 5},
			},
			latestDate: "2026-05-26",
		}
		service := NewStatsServiceWithRepo(repo)

		resp, err := service.GetTrends(7)
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "2026-05-26", resp.UpdatedAt)
		assert.NotEmpty(t, resp.CityTrends)
		assert.NotEmpty(t, resp.IndustryTrends)
	})

	t.Run("repo_error", func(t *testing.T) {
		repo := &mockStatsRepo{trendsErr: assert.AnError}
		service := NewStatsServiceWithRepo(repo)

		_, err := service.GetTrends(7)
		require.Error(t, err)
	})

	t.Run("default_days", func(t *testing.T) {
		repo := &mockStatsRepo{stats: []DailyStat{}, latestDate: "2026-05-26"}
		service := NewStatsServiceWithRepo(repo)

		// days=0 should default to 7
		resp, err := service.GetTrends(0)
		require.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("max_days", func(t *testing.T) {
		repo := &mockStatsRepo{stats: []DailyStat{}, latestDate: "2026-05-26"}
		service := NewStatsServiceWithRepo(repo)

		// days=365 should be capped at 90
		resp, err := service.GetTrends(365)
		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}
