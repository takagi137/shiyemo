package cohort

import (
	"errors"
	"testing"
	"time"

	"github.com/shiyemo/server/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================
// Mock 实现
// ============================================================

type mockCohortRepo struct {
	users      map[int64]*auth.User
	stats      map[string]*CohortStats
	cities     []string
	industries []string
	// 控制返回错误
	statsErr   error
	citiesErr  error
	industErr  error
	listErr    error
	myStatsErr error
}

func newMockCohortRepo() *mockCohortRepo {
	return &mockCohortRepo{
		users:      make(map[int64]*auth.User),
		stats:      make(map[string]*CohortStats),
		cities:     []string{"北京", "上海", "深圳"},
		industries: []string{"互联网", "金融", "教育"},
	}
}

func (m *mockCohortRepo) GetCohortStats(city, industry string, timeRange CohortTimeRange) (*CohortStats, error) {
	if m.statsErr != nil {
		return nil, m.statsErr
	}
	key := city + ":" + industry + ":" + string(timeRange)
	if stats, ok := m.stats[key]; ok {
		return stats, nil
	}
	return &CohortStats{
		City:      city,
		Industry:  industry,
		UserCount: 100,
		PostCount: 500,
	}, nil
}

func (m *mockCohortRepo) ListPeers(filter CohortFilter) ([]CohortPeer, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	var peers []CohortPeer
	for _, u := range m.users {
		if filter.City != "" && u.City != filter.City {
			continue
		}
		if filter.Industry != "" && u.Industry != filter.Industry {
			continue
		}
		peers = append(peers, CohortPeer{
			ID:         u.ID,
			Nickname:   u.Nickname,
			Avatar:     u.Avatar,
			Bio:        u.Bio,
			City:       u.City,
			Industry:   u.Industry,
			LeftJobAt:  u.LeftJobAt,
			FreedomDays: u.FreedomDays(),
		})
	}
	return peers, int64(len(peers)), nil
}

func (m *mockCohortRepo) GetMyCohortStats(userID int64) (*CohortStats, error) {
	if m.myStatsErr != nil {
		return nil, m.myStatsErr
	}
	if u, ok := m.users[userID]; ok {
		if u.City == "" && u.Industry == "" {
			return &CohortStats{
				City:      u.City,
				Industry:  u.Industry,
				UserCount: 0,
				PostCount: 0,
			}, nil
		}
		return m.GetCohortStats(u.City, u.Industry, TimeRangeAll)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockCohortRepo) GetCities() ([]string, error) {
	if m.citiesErr != nil {
		return nil, m.citiesErr
	}
	return m.cities, nil
}

func (m *mockCohortRepo) GetIndustries() ([]string, error) {
	if m.industErr != nil {
		return nil, m.industErr
	}
	return m.industries, nil
}

// ============================================================
// GetMyStats 测试
// ============================================================

func TestCohortService_GetMyStats(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockCohortRepo()
		leftJobAt := time.Now().AddDate(0, -1, 0)
		repo.users[100] = &auth.User{
			ID:        100,
			City:      "北京",
			Industry:  "互联网",
			LeftJobAt: &leftJobAt,
		}

		svc := NewCohortServiceWithRepo(repo)
		stats, err := svc.GetMyStats(100)

		require.NoError(t, err)
		assert.Equal(t, "北京", stats.City)
		assert.Equal(t, "互联网", stats.Industry)
		assert.Equal(t, int64(100), stats.UserCount)
	})

	t.Run("user_not_found", func(t *testing.T) {
		repo := newMockCohortRepo()
		svc := NewCohortServiceWithRepo(repo)

		_, err := svc.GetMyStats(999)
		require.Error(t, err)
	})

	t.Run("user_without_city_industry", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.users[200] = &auth.User{
			ID:       200,
			Nickname: "测试用户",
		}

		svc := NewCohortServiceWithRepo(repo)
		stats, err := svc.GetMyStats(200)

		require.NoError(t, err)
		assert.Equal(t, "", stats.City)
		assert.Equal(t, "", stats.Industry)
		assert.Equal(t, int64(0), stats.UserCount)
	})

	t.Run("repo_error", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.myStatsErr = errors.New("database connection lost")

		svc := NewCohortServiceWithRepo(repo)
		_, err := svc.GetMyStats(100)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "database connection lost")
	})

	t.Run("success_with_career_stats", func(t *testing.T) {
		repo := newMockCohortRepo()
		leftJobAt := time.Now().AddDate(0, -3, 0)
		repo.users[300] = &auth.User{
			ID:        300,
			City:      "深圳",
			Industry:  "金融",
			LeftJobAt: &leftJobAt,
		}
		key := "深圳:金融:" + string(TimeRangeAll)
		repo.stats[key] = &CohortStats{
			City:      "深圳",
			Industry:  "金融",
			UserCount: 25,
			PostCount: 120,
		}

		svc := NewCohortServiceWithRepo(repo)
		stats, err := svc.GetMyStats(300)

		require.NoError(t, err)
		assert.Equal(t, "深圳", stats.City)
		assert.Equal(t, "金融", stats.Industry)
		assert.Equal(t, int64(25), stats.UserCount)
		assert.Equal(t, int64(120), stats.PostCount)
	})
}

// ============================================================
// GetStats 测试
// ============================================================

func TestCohortService_GetStats(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockCohortRepo()
		svc := NewCohortServiceWithRepo(repo)

		stats, err := svc.GetStats("北京", "互联网", TimeRangeMonth)

		require.NoError(t, err)
		assert.Equal(t, "北京", stats.City)
		assert.Equal(t, "互联网", stats.Industry)
	})

	t.Run("empty_params", func(t *testing.T) {
		repo := newMockCohortRepo()
		svc := NewCohortServiceWithRepo(repo)

		stats, err := svc.GetStats("", "", TimeRangeAll)

		require.NoError(t, err)
		assert.NotNil(t, stats)
	})

	t.Run("repo_error", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.statsErr = errors.New("redis timeout")

		svc := NewCohortServiceWithRepo(repo)
		_, err := svc.GetStats("北京", "互联网", TimeRangeAll)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "redis timeout")
	})

	t.Run("with_cached_stats", func(t *testing.T) {
		repo := newMockCohortRepo()
		key := "上海:教育:" + string(TimeRangeYear)
		repo.stats[key] = &CohortStats{
			City:      "上海",
			Industry:  "教育",
			UserCount: 300,
			PostCount: 1500,
			ActiveDays: 45,
		}

		svc := NewCohortServiceWithRepo(repo)
		stats, err := svc.GetStats("上海", "教育", TimeRangeYear)

		require.NoError(t, err)
		assert.Equal(t, int64(300), stats.UserCount)
		assert.Equal(t, int64(1500), stats.PostCount)
		assert.Equal(t, 45, stats.ActiveDays)
	})

	t.Run("week_time_range", func(t *testing.T) {
		repo := newMockCohortRepo()

		svc := NewCohortServiceWithRepo(repo)
		stats, err := svc.GetStats("北京", "互联网", TimeRangeWeek)

		require.NoError(t, err)
		assert.NotNil(t, stats)
	})

	t.Run("year_time_range", func(t *testing.T) {
		repo := newMockCohortRepo()

		svc := NewCohortServiceWithRepo(repo)
		stats, err := svc.GetStats("北京", "互联网", TimeRangeYear)

		require.NoError(t, err)
		assert.NotNil(t, stats)
	})
}

// ============================================================
// ListPeers 测试
// ============================================================

func TestCohortService_ListPeers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockCohortRepo()
		leftJobAt := time.Now().AddDate(0, -2, 0)
		repo.users[1] = &auth.User{ID: 1, Nickname: "用户1", City: "北京", Industry: "互联网", LeftJobAt: &leftJobAt}
		repo.users[2] = &auth.User{ID: 2, Nickname: "用户2", City: "北京", Industry: "金融"}
		repo.users[3] = &auth.User{ID: 3, Nickname: "用户3", City: "上海", Industry: "互联网"}

		svc := NewCohortServiceWithRepo(repo)

		filter := CohortFilter{
			City:     "北京",
			Page:     1,
			PageSize: 20,
		}
		peers, total, err := svc.ListPeers(filter)

		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, peers, 2)
	})

	t.Run("default_pagination", func(t *testing.T) {
		repo := newMockCohortRepo()
		svc := NewCohortServiceWithRepo(repo)

		filter := CohortFilter{
			Page:     0,
			PageSize: 0,
		}
		_, _, err := svc.ListPeers(filter)

		require.NoError(t, err)
	})

	t.Run("repo_error", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.listErr = errors.New("query timeout")

		svc := NewCohortServiceWithRepo(repo)
		_, _, err := svc.ListPeers(CohortFilter{City: "北京"})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "query timeout")
	})

	t.Run("filter_by_industry", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.users[1] = &auth.User{ID: 1, Nickname: "用户1", City: "北京", Industry: "互联网"}
		repo.users[2] = &auth.User{ID: 2, Nickname: "用户2", City: "北京", Industry: "金融"}
		repo.users[3] = &auth.User{ID: 3, Nickname: "用户3", City: "北京", Industry: "互联网"}

		svc := NewCohortServiceWithRepo(repo)

		filter := CohortFilter{
			City:      "北京",
			Industry:  "互联网",
			Page:      1,
			PageSize:  10,
		}
		peers, total, err := svc.ListPeers(filter)

		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, peers, 2)
	})

	t.Run("filter_city_and_industry", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.users[1] = &auth.User{ID: 1, Nickname: "用户1", City: "北京", Industry: "互联网"}
		repo.users[2] = &auth.User{ID: 2, Nickname: "用户2", City: "上海", Industry: "互联网"}
		repo.users[3] = &auth.User{ID: 3, Nickname: "用户3", City: "北京", Industry: "金融"}

		svc := NewCohortServiceWithRepo(repo)

		filter := CohortFilter{
			City:      "北京",
			Industry:  "互联网",
			Page:      1,
			PageSize:  10,
		}
		peers, total, err := svc.ListPeers(filter)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, peers, 1)
		assert.Equal(t, "用户1", peers[0].Nickname)
	})

	t.Run("no_match", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.users[1] = &auth.User{ID: 1, Nickname: "用户1", City: "北京", Industry: "互联网"}

		svc := NewCohortServiceWithRepo(repo)

		filter := CohortFilter{
			City:     "广州",
			Page:     1,
			PageSize: 10,
		}
		peers, total, err := svc.ListPeers(filter)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, peers, 0)
	})

	t.Run("page_size_too_large", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.users[1] = &auth.User{ID: 1, Nickname: "用户1", City: "北京"}

		svc := NewCohortServiceWithRepo(repo)

		filter := CohortFilter{
			Page:     1,
			PageSize: 100, // > 50, 应被限制为 20
		}
		peers, total, err := svc.ListPeers(filter)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, peers, 1)
	})

	t.Run("negative_page", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.users[1] = &auth.User{ID: 1, Nickname: "用户1", City: "北京"}

		svc := NewCohortServiceWithRepo(repo)

		filter := CohortFilter{
			Page:     -1,
			PageSize: -1,
		}
		peers, total, err := svc.ListPeers(filter)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, peers, 1)
	})

	t.Run("with_user_data_fields", func(t *testing.T) {
		repo := newMockCohortRepo()
		lastWorkDate := time.Now().AddDate(0, -6, 0)
		leftJobAt := time.Now().AddDate(0, -6, 0)
		repo.users[50] = &auth.User{
			ID:           50,
			Nickname:     "张三",
			Avatar:       "https://example.com/avatar.jpg",
			Bio:          "互联网创业者",
			City:         "杭州",
			Industry:     "互联网",
			LeftJobAt:    &leftJobAt,
			LastWorkDate: &lastWorkDate,
		}

		svc := NewCohortServiceWithRepo(repo)

		filter := CohortFilter{
			City:     "杭州",
			Page:     1,
			PageSize: 10,
		}
		peers, total, err := svc.ListPeers(filter)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		require.Len(t, peers, 1)
		assert.Equal(t, "https://example.com/avatar.jpg", peers[0].Avatar)
		assert.Equal(t, "互联网创业者", peers[0].Bio)
		assert.NotNil(t, peers[0].LeftJobAt)
		assert.NotNil(t, peers[0].FreedomDays)
	})

	t.Run("empty_list", func(t *testing.T) {
		repo := newMockCohortRepo()

		svc := NewCohortServiceWithRepo(repo)

		filter := CohortFilter{
			Page:     1,
			PageSize: 10,
		}
		peers, total, err := svc.ListPeers(filter)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, peers, 0)
	})
}

// ============================================================
// GetCities 测试
// ============================================================

func TestCohortService_GetCities(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockCohortRepo()
		svc := NewCohortServiceWithRepo(repo)

		cities, err := svc.GetCities()

		require.NoError(t, err)
		assert.Equal(t, []string{"北京", "上海", "深圳"}, cities)
	})

	t.Run("repo_error", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.citiesErr = errors.New("db connection lost")

		svc := NewCohortServiceWithRepo(repo)
		_, err := svc.GetCities()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db connection lost")
	})

	t.Run("empty_result", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.cities = []string{}

		svc := NewCohortServiceWithRepo(repo)
		cities, err := svc.GetCities()

		require.NoError(t, err)
		assert.Empty(t, cities)
	})
}

// ============================================================
// GetIndustries 测试
// ============================================================

func TestCohortService_GetIndustries(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockCohortRepo()
		svc := NewCohortServiceWithRepo(repo)

		industries, err := svc.GetIndustries()

		require.NoError(t, err)
		assert.Equal(t, []string{"互联网", "金融", "教育"}, industries)
	})

	t.Run("repo_error", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.industErr = errors.New("timeout")

		svc := NewCohortServiceWithRepo(repo)
		_, err := svc.GetIndustries()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "timeout")
	})

	t.Run("empty_result", func(t *testing.T) {
		repo := newMockCohortRepo()
		repo.industries = []string{}

		svc := NewCohortServiceWithRepo(repo)
		industries, err := svc.GetIndustries()

		require.NoError(t, err)
		assert.Empty(t, industries)
	})
}
