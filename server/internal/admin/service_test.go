package admin

import (
	"errors"
	"testing"

	"github.com/shiyemo/server/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================
// Mock 实现
// ============================================================

type mockAdminRepo struct {
	posts      []map[string]interface{}
	pending    []map[string]interface{}
	users      []auth.User
	stats      map[string]interface{}
	deletedIDs []int64
	topIDs     map[int64]bool
	featuredIDs map[int64]bool
	approvedIDs []int64
	rejectedIDs map[int64]string
	userMap    map[int64]*auth.User
	userStatus map[int64]int8
}

func newMockAdminRepo() *mockAdminRepo {
	return &mockAdminRepo{
		topIDs:      make(map[int64]bool),
		featuredIDs: make(map[int64]bool),
		rejectedIDs: make(map[int64]string),
		userMap:     make(map[int64]*auth.User),
		userStatus:  make(map[int64]int8),
	}
}

func (m *mockAdminRepo) ListPosts(page, pageSize int, status string) ([]map[string]interface{}, int64, error) {
	return m.posts, int64(len(m.posts)), nil
}

func (m *mockAdminRepo) DeletePost(postID int64) error {
	m.deletedIDs = append(m.deletedIDs, postID)
	return nil
}

func (m *mockAdminRepo) ListUsers(page, pageSize int, search string, status *int) ([]auth.User, int64, error) {
	return m.users, int64(len(m.users)), nil
}

func (m *mockAdminRepo) UpdateUserStatus(userID int64, status int8) error {
	m.userStatus[userID] = status
	return nil
}

func (m *mockAdminRepo) GetUserByID(userID int64) (*auth.User, error) {
	if u, ok := m.userMap[userID]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAdminRepo) GetStats() (map[string]interface{}, error) {
	return m.stats, nil
}

func (m *mockAdminRepo) SetPostTop(postID int64, isTop bool) error {
	m.topIDs[postID] = isTop
	return nil
}

func (m *mockAdminRepo) SetPostFeatured(postID int64, isFeatured bool) error {
	m.featuredIDs[postID] = isFeatured
	return nil
}

func (m *mockAdminRepo) ListPendingPosts(page, pageSize int) ([]map[string]interface{}, int64, error) {
	return m.pending, int64(len(m.pending)), nil
}

func (m *mockAdminRepo) ApprovePost(postID, reviewerID int64) error {
	for _, p := range m.pending {
		if pid, ok := p["id"]; ok && pid.(int64) == postID {
			m.approvedIDs = append(m.approvedIDs, postID)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockAdminRepo) RejectPost(postID, reviewerID int64, reason string) error {
	for _, p := range m.pending {
		if pid, ok := p["id"]; ok && pid.(int64) == postID {
			m.rejectedIDs[postID] = reason
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockAdminRepo) ListComments(page, pageSize int, postID int64) ([]map[string]interface{}, int64, error) {
	return []map[string]interface{}{}, 0, nil
}

func (m *mockAdminRepo) UpdateCommentStatus(commentID int64, status int) error {
	return nil
}

// ============================================================
// 测试
// ============================================================

func TestService_ListPendingPosts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := newMockAdminRepo()
		mock.pending = []map[string]interface{}{
			{"id": int64(1), "title": "待审核帖子1"},
			{"id": int64(2), "title": "待审核帖子2"},
		}
		svc := NewServiceWithRepo(mock)

		posts, total, err := svc.ListPendingPosts(1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, posts, 2)
	})

	t.Run("empty", func(t *testing.T) {
		mock := newMockAdminRepo()
		svc := NewServiceWithRepo(mock)

		posts, total, err := svc.ListPendingPosts(1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, posts, 0)
	})

	t.Run("default_params", func(t *testing.T) {
		mock := newMockAdminRepo()
		svc := NewServiceWithRepo(mock)

		_, _, err := svc.ListPendingPosts(0, 0)
		require.NoError(t, err)
	})
}

func TestService_ApprovePost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := newMockAdminRepo()
		mock.pending = []map[string]interface{}{
			{"id": int64(1), "title": "待审核帖子"},
		}
		svc := NewServiceWithRepo(mock)

		err := svc.ApprovePost(1, 100)
		require.NoError(t, err)
		assert.Contains(t, mock.approvedIDs, int64(1))
	})

	t.Run("post_not_found", func(t *testing.T) {
		mock := newMockAdminRepo()
		svc := NewServiceWithRepo(mock)

		err := svc.ApprovePost(99999, 1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "帖子不存在")
	})
}

func TestService_RejectPost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := newMockAdminRepo()
		mock.pending = []map[string]interface{}{
			{"id": int64(1), "title": "违规帖子"},
		}
		svc := NewServiceWithRepo(mock)

		err := svc.RejectPost(1, 100, "违规内容")
		require.NoError(t, err)
		assert.Equal(t, "违规内容", mock.rejectedIDs[int64(1)])
	})

	t.Run("post_not_found", func(t *testing.T) {
		mock := newMockAdminRepo()
		svc := NewServiceWithRepo(mock)

		err := svc.RejectPost(99999, 1, "违规")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "帖子不存在")
	})

	t.Run("empty_reason", func(t *testing.T) {
		mock := newMockAdminRepo()
		svc := NewServiceWithRepo(mock)

		err := svc.RejectPost(1, 100, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "拒绝原因不能为空")
	})
}

func TestService_GetStats(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := newMockAdminRepo()
		mock.stats = map[string]interface{}{
			"total_users":    10,
			"total_posts":    25,
			"total_comments": 50,
		}
		svc := NewServiceWithRepo(mock)

		stats, err := svc.GetStats()

		require.NoError(t, err)
		assert.Equal(t, 10, stats["total_users"])
		assert.Equal(t, 25, stats["total_posts"])
		assert.Equal(t, 50, stats["total_comments"])
	})
}

func TestService_SetPostTop(t *testing.T) {
	t.Run("set_top", func(t *testing.T) {
		mock := newMockAdminRepo()
		svc := NewServiceWithRepo(mock)

		err := svc.SetPostTop(1, true)
		require.NoError(t, err)
		assert.True(t, mock.topIDs[int64(1)])
	})

	t.Run("unset_top", func(t *testing.T) {
		mock := newMockAdminRepo()
		mock.topIDs[int64(1)] = true
		svc := NewServiceWithRepo(mock)

		err := svc.SetPostTop(1, false)
		require.NoError(t, err)
		assert.False(t, mock.topIDs[int64(1)])
	})
}

func TestService_SetPostFeatured(t *testing.T) {
	t.Run("set_featured", func(t *testing.T) {
		mock := newMockAdminRepo()
		svc := NewServiceWithRepo(mock)

		err := svc.SetPostFeatured(1, true)
		require.NoError(t, err)
		assert.True(t, mock.featuredIDs[int64(1)])
	})
}

func TestService_DeletePost(t *testing.T) {
	mock := newMockAdminRepo()
	svc := NewServiceWithRepo(mock)

	err := svc.DeletePost(42)
	require.NoError(t, err)
	assert.Contains(t, mock.deletedIDs, int64(42))
}

func TestService_ListPosts(t *testing.T) {
	t.Run("with_status", func(t *testing.T) {
		mock := newMockAdminRepo()
		mock.posts = []map[string]interface{}{
			{"id": int64(1), "title": "帖子1"},
		}
		svc := NewServiceWithRepo(mock)

		posts, total, err := svc.ListPosts(1, 10, "active")

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, posts, 1)
	})

	t.Run("default_params", func(t *testing.T) {
		mock := newMockAdminRepo()
		svc := NewServiceWithRepo(mock)

		_, _, err := svc.ListPosts(0, 0, "")
		require.NoError(t, err)
	})
}

func TestService_ListUsers(t *testing.T) {
	mock := newMockAdminRepo()
	mock.users = []auth.User{
		{ID: 1, Nickname: "用户1"},
		{ID: 2, Nickname: "用户2"},
	}
	svc := NewServiceWithRepo(mock)

	users, total, err := svc.ListUsers(1, 10, "", nil)

	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, users, 2)
}

func TestService_UpdateUserStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := newMockAdminRepo()
		mock.userMap[int64(1)] = &auth.User{ID: 1, Nickname: "test"}
		svc := NewServiceWithRepo(mock)

		err := svc.UpdateUserStatus(1, 0)
		require.NoError(t, err)
		assert.Equal(t, int8(0), mock.userStatus[int64(1)])
	})

	t.Run("user_not_found", func(t *testing.T) {
		mock := newMockAdminRepo()
		svc := NewServiceWithRepo(mock)

		err := svc.UpdateUserStatus(99999, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "用户不存在")
	})

	t.Run("db_error", func(t *testing.T) {
		mock := newMockAdminRepo()
		mock.userMap[int64(1)] = &auth.User{ID: 1}

		// 模拟 repo 返回非 NotFound 错误的场景
		// 通过覆盖 GetUserByID 行为
		customMock := &errorMockAdminRepo{
			mockAdminRepo: mock,
			getUserErr:    errors.New("数据库连接失败"),
		}
		svc2 := NewServiceWithRepo(customMock)
		err := svc2.UpdateUserStatus(1, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "数据库连接失败")
	})
}

// ============================================================
// ListComments / HideComment / RestoreComment 测试
// ============================================================

func TestService_ListComments(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := newMockAdminRepo()
		svc := NewServiceWithRepo(mock)

		comments, total, err := svc.ListComments(1, 10, 0)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, comments, 0)
	})

	t.Run("default_params", func(t *testing.T) {
		mock := newMockAdminRepo()
		svc := NewServiceWithRepo(mock)

		_, _, err := svc.ListComments(0, 0, 0)
		require.NoError(t, err)
	})
}

func TestService_HideComment(t *testing.T) {
	mock := newMockAdminRepo()
	svc := NewServiceWithRepo(mock)

	err := svc.HideComment(1)
	require.NoError(t, err)
}

func TestService_RestoreComment(t *testing.T) {
	mock := newMockAdminRepo()
	svc := NewServiceWithRepo(mock)

	err := svc.RestoreComment(1)
	require.NoError(t, err)
}

// errorMockAdminRepo 用于模拟各类错误的 mock
type errorMockAdminRepo struct {
	*mockAdminRepo
	getUserErr error
}

func (m *errorMockAdminRepo) GetUserByID(userID int64) (*auth.User, error) {
	if m.getUserErr != nil {
		return nil, m.getUserErr
	}
	return m.mockAdminRepo.GetUserByID(userID)
}
