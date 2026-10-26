package user

import (
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

type mockUserRepo struct {
	users  map[int64]*auth.User
	nextID int64
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[int64]*auth.User)}
}

func (m *mockUserRepo) GetByID(id int64) (*auth.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) Update(user *auth.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) UpdateFields(id int64, fields map[string]interface{}) error {
	if u, ok := m.users[id]; ok {
		if v, ok := fields["nickname"]; ok {
			u.Nickname = v.(string)
		}
		if v, ok := fields["bio"]; ok {
			u.Bio = v.(string)
		}
		if v, ok := fields["company"]; ok {
			u.Company = v.(string)
		}
		if v, ok := fields["occupation"]; ok {
			u.Occupation = v.(string)
		}
		if v, ok := fields["city"]; ok {
			u.City = v.(string)
		}
		if v, ok := fields["industry"]; ok {
			u.Industry = v.(string)
		}
		if v, ok := fields["left_job_at"]; ok {
			if v == nil {
				u.LeftJobAt = nil
			} else {
				t := v.(time.Time)
				u.LeftJobAt = &t
			}
		}
		if v, ok := fields["avatar"]; ok {
			u.Avatar = v.(string)
		}
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *mockUserRepo) GetUserPosts(userID int64, page, pageSize int) ([]map[string]interface{}, int64, error) {
	var result []map[string]interface{}
	for _, u := range m.users {
		if u.ID == userID {
			break
		}
	}
	return result, 0, nil
}

func (m *mockUserRepo) CountUserPosts(userID int64) (int64, error) {
	return 5, nil
}

func (m *mockUserRepo) CountUserLikesReceived(userID int64) (int64, error) {
	return 10, nil
}

type mockUserRepoWithPosts struct {
	*mockUserRepo
	posts map[int64][]map[string]interface{}
}

func (m *mockUserRepoWithPosts) GetUserPosts(userID int64, page, pageSize int) ([]map[string]interface{}, int64, error) {
	posts, ok := m.posts[userID]
	if !ok {
		return nil, 0, nil
	}
	total := int64(len(posts))
	start := (page - 1) * pageSize
	if start >= int(total) {
		return nil, total, nil
	}
	end := start + pageSize
	if end > int(total) {
		end = int(total)
	}
	return posts[start:end], total, nil
}

func (m *mockUserRepoWithPosts) CountUserPosts(userID int64) (int64, error) {
	return int64(len(m.posts[userID])), nil
}

// ============================================================
// 工具函数
// ============================================================

func createMockService() (*UserService, *mockUserRepo) {
	mock := newMockUserRepo()
	svc := NewUserServiceWithRepo(mock)
	return svc, mock
}

func seedUsers(mock *mockUserRepo, count int) {
	for i := 1; i <= count; i++ {
		mock.users[int64(i)] = &auth.User{
			ID:       int64(i),
			Nickname: "用户" + string(rune('A'+i-1)),
			Email:    strPtr("user" + string(rune('a'+i-1)) + "@test.com"),
			Bio:      "简介",
		}
	}
}

func strPtr(s string) *string {
	return &s
}

// ============================================================
// 核心功能测试
// ============================================================

func TestUserService_GetProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock := createMockService()
		mock.users[100] = &auth.User{ID: 100, Nickname: "张三"}

		user, err := svc.GetProfile(100)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "张三", user.Nickname)
	})

	t.Run("not_found", func(t *testing.T) {
		svc, _ := createMockService()

		_, err := svc.GetProfile(99999)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "用户不存在")
	})
}

func TestUserService_GetPublicProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock := createMockService()
		mock.users[100] = &auth.User{
			ID:         100,
			Nickname:   "李四",
			Avatar:     "avatar.png",
			Bio:        "热爱编程",
			Company:    "科技公司",
			Occupation: "工程师",
			CreatedAt:  time.Now(),
		}

		profile, err := svc.GetPublicProfile(100)
		require.NoError(t, err)
		assert.Equal(t, "李四", profile["nickname"])
		assert.Equal(t, "avatar.png", profile["avatar"])
		assert.Equal(t, "热爱编程", profile["bio"])
	})

	t.Run("not_found", func(t *testing.T) {
		svc, _ := createMockService()

		_, err := svc.GetPublicProfile(99999)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "用户不存在")
	})
}

func TestUserService_UpdateProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock := createMockService()
		mock.users[100] = &auth.User{ID: 100, Nickname: "原昵称"}

		err := svc.UpdateProfile(100, UpdateProfileRequest{
			Nickname: "新昵称",
			Bio:      "新简介",
		})
		require.NoError(t, err)
		assert.Equal(t, "新昵称", mock.users[100].Nickname)
		assert.Equal(t, "新简介", mock.users[100].Bio)
	})

	t.Run("not_found", func(t *testing.T) {
		svc, _ := createMockService()

		err := svc.UpdateProfile(99999, UpdateProfileRequest{Nickname: "新昵称"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "用户不存在")
	})

	t.Run("empty_fields", func(t *testing.T) {
		svc, mock := createMockService()
		mock.users[100] = &auth.User{ID: 100, Nickname: "原昵称"}

		err := svc.UpdateProfile(100, UpdateProfileRequest{})
		require.NoError(t, err)
		assert.Equal(t, "原昵称", mock.users[100].Nickname)
	})

	t.Run("partial_fields", func(t *testing.T) {
		svc, mock := createMockService()
		mock.users[100] = &auth.User{ID: 100, Nickname: "原名", Bio: "旧简介"}

		err := svc.UpdateProfile(100, UpdateProfileRequest{Bio: "新简介"})
		require.NoError(t, err)
		assert.Equal(t, "原名", mock.users[100].Nickname)
		assert.Equal(t, "新简介", mock.users[100].Bio)
	})

	t.Run("all_fields", func(t *testing.T) {
		svc, mock := createMockService()
		mock.users[100] = &auth.User{ID: 100}

		err := svc.UpdateProfile(100, UpdateProfileRequest{
			Nickname:   "昵称",
			Bio:        "简介",
			Company:    "公司",
			Occupation: "职位",
		})
		require.NoError(t, err)
		assert.Equal(t, "昵称", mock.users[100].Nickname)
		assert.Equal(t, "公司", mock.users[100].Company)
	})

	t.Run("new_profile_fields", func(t *testing.T) {
		svc, mock := createMockService()
		mock.users[100] = &auth.User{ID: 100}

		leftJobAt := "2024-01-15"
		err := svc.UpdateProfile(100, UpdateProfileRequest{
			City:      "北京",
			Industry:  "互联网",
			LeftJobAt: &leftJobAt,
		})
		require.NoError(t, err)
		assert.Equal(t, "北京", mock.users[100].City)
		assert.Equal(t, "互联网", mock.users[100].Industry)
		assert.NotNil(t, mock.users[100].LeftJobAt)
	})
}

func TestUserService_GetUserStats(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock := createMockService()
		mock.users[100] = &auth.User{ID: 100, FollowerCount: 15, FollowingCount: 8}

		stats, err := svc.GetUserStats(100)
		require.NoError(t, err)
		assert.Equal(t, int64(5), stats["post_count"])
		assert.Equal(t, int64(10), stats["like_received"])
		assert.Equal(t, 15, stats["follower_count"])
		assert.Equal(t, 8, stats["following_count"])
	})

	t.Run("not_found", func(t *testing.T) {
		svc, _ := createMockService()
		_, err := svc.GetUserStats(99999)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "用户不存在")
	})
}

