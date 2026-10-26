package follow

import (
	"testing"

	"github.com/shiyemo/server/internal/auth"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================
// Mock 实现
// ============================================================

type mockFollowRepo struct {
	following map[int64]map[int64]bool // followerID -> set of followingID
}

func newMockFollowRepo() *mockFollowRepo {
	return &mockFollowRepo{following: make(map[int64]map[int64]bool)}
}

func (m *mockFollowRepo) Follow(followerID, followingID int64) error {
	if m.following[followerID] == nil {
		m.following[followerID] = make(map[int64]bool)
	}
	m.following[followerID][followingID] = true
	return nil
}

func (m *mockFollowRepo) Unfollow(followerID, followingID int64) error {
	delete(m.following[followerID], followingID)
	return nil
}

func (m *mockFollowRepo) IsFollowing(followerID, followingID int64) bool {
	return m.following[followerID] != nil && m.following[followerID][followingID]
}

func (m *mockFollowRepo) ListFollowing(userID int64, page, size int) ([]int64, int64, error) {
	var ids []int64
	for fid := range m.following[userID] {
		ids = append(ids, fid)
	}
	total := int64(len(ids))
	start := (page - 1) * size
	if start >= len(ids) {
		return nil, total, nil
	}
	end := start + size
	if end > len(ids) {
		end = len(ids)
	}
	return ids[start:end], total, nil
}

func (m *mockFollowRepo) ListFollowers(userID int64, page, size int) ([]int64, int64, error) {
	var ids []int64
	for followerID, targets := range m.following {
		if targets[userID] {
			ids = append(ids, followerID)
		}
	}
	total := int64(len(ids))
	start := (page - 1) * size
	if start >= len(ids) {
		return nil, total, nil
	}
	end := start + size
	if end > len(ids) {
		end = len(ids)
	}
	return ids[start:end], total, nil
}

func (m *mockFollowRepo) CountFollowing(userID int64) (int64, error) {
	return int64(len(m.following[userID])), nil
}

func (m *mockFollowRepo) CountFollowers(userID int64) (int64, error) {
	var count int64
	for _, targets := range m.following {
		if targets[userID] {
			count++
		}
	}
	return count, nil
}

type mockUserProvider struct {
	users map[int64]*auth.User
}

func newMockUserProvider() *mockUserProvider {
	return &mockUserProvider{users: map[int64]*auth.User{
		1: {ID: 1, Nickname: "用户A", Avatar: "a.png", Bio: "bio A"},
		2: {ID: 2, Nickname: "用户B", Avatar: "b.png", Bio: "bio B"},
	}}
}

func (m *mockUserProvider) GetByID(id int64) (*auth.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

type mockNotifier struct {
	calls []notifyCall
}

type notifyCall struct {
	UserID, ActorID int64
	NType           string
	TargetID        *int64
	TargetType      *string
	Content         string
}

func (m *mockNotifier) Create(userID, actorID int64, nType string, targetID *int64, targetType *string, content string) error {
	m.calls = append(m.calls, notifyCall{userID, actorID, nType, targetID, targetType, content})
	return nil
}

type mockFollowCountUpdater struct {
	lastFollowerID  int64
	lastFollowingID int64
	lastFC          int64
	lastFG          int64
}

func (m *mockFollowCountUpdater) UpdateCounts(followerID, followingID int64, followingCount, followerCount int64) {
	m.lastFollowerID = followerID
	m.lastFollowingID = followingID
	m.lastFC = followingCount
	m.lastFG = followerCount
}

// ============================================================
// 测试辅助
// ============================================================

type testDeps struct {
	repo      *mockFollowRepo
	userRepo  *mockUserProvider
	notifySvc *mockNotifier
	countUpd  *mockFollowCountUpdater
	service   *FollowService
}

func newTestDeps() *testDeps {
	repo := newMockFollowRepo()
	userRepo := newMockUserProvider()
	notifySvc := &mockNotifier{}
	countUpd := &mockFollowCountUpdater{}
	service := NewFollowServiceWithDeps(repo, userRepo, notifySvc, countUpd)
	return &testDeps{repo, userRepo, notifySvc, countUpd, service}
}

// ============================================================
// Follow 测试
// ============================================================

func TestFollowService_Follow(t *testing.T) {
	t.Run("follow_success", func(t *testing.T) {
		d := newTestDeps()
		following, err := d.service.Follow(1, 2)

		require.NoError(t, err)
		assert.True(t, following)
	})

	t.Run("unfollow_success", func(t *testing.T) {
		d := newTestDeps()

		d.service.Follow(1, 2)
		following, err := d.service.Follow(1, 2)

		require.NoError(t, err)
		assert.False(t, following)
	})

	t.Run("cannot_follow_self", func(t *testing.T) {
		d := newTestDeps()
		_, err := d.service.Follow(1, 1)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "不能关注自己")
	})

	t.Run("user_not_found", func(t *testing.T) {
		d := newTestDeps()
		_, err := d.service.Follow(1, 999)

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrUserNotFound, err)
	})

	t.Run("notification_sent", func(t *testing.T) {
		d := newTestDeps()
		d.service.Follow(1, 2)

		require.Len(t, d.notifySvc.calls, 1)
		assert.Equal(t, "follow", d.notifySvc.calls[0].NType)
		assert.Equal(t, int64(2), d.notifySvc.calls[0].UserID) // 被关注者
		assert.Equal(t, int64(1), d.notifySvc.calls[0].ActorID)
	})

	t.Run("no_notification_on_unfollow", func(t *testing.T) {
		d := newTestDeps()

		d.service.Follow(1, 2) // 关注 — 通知
		d.service.Follow(1, 2) // 取消 — 无通知

		assert.Len(t, d.notifySvc.calls, 1)
	})

	t.Run("counts_updated", func(t *testing.T) {
		d := newTestDeps()
		d.service.Follow(1, 2)

		assert.Equal(t, int64(1), d.countUpd.lastFollowerID)
		assert.Equal(t, int64(2), d.countUpd.lastFollowingID)
		assert.Equal(t, int64(1), d.countUpd.lastFC) // 关注数
		assert.Equal(t, int64(1), d.countUpd.lastFG) // 粉丝数
	})
}

// ============================================================
// ListFollowing 测试
// ============================================================

func TestFollowService_ListFollowing(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		d := newTestDeps()

		d.service.Follow(1, 2)
		users, total, err := d.service.ListFollowing(1, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, users, 1)
		assert.Equal(t, "用户B", users[0]["nickname"])
	})

	t.Run("empty", func(t *testing.T) {
		d := newTestDeps()

		users, total, err := d.service.ListFollowing(1, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, users, 0)
	})
}

// ============================================================
// ListFollowers 测试
// ============================================================

func TestFollowService_ListFollowers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		d := newTestDeps()

		d.service.Follow(1, 2) // 1 关注 2
		users, total, err := d.service.ListFollowers(2, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, users, 1)
		assert.Equal(t, "用户A", users[0]["nickname"])
	})

	t.Run("empty", func(t *testing.T) {
		d := newTestDeps()

		users, total, err := d.service.ListFollowers(1, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, users, 0)
	})
}
