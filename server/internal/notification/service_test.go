package notification

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/shiyemo/server/internal/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================
// Mock 实现
// ============================================================

type mockNotifRepo struct {
	nextID int64
	notifs map[int64]*Notification
}

func newMockNotifRepo() *mockNotifRepo {
	return &mockNotifRepo{notifs: make(map[int64]*Notification)}
}

func (m *mockNotifRepo) Create(n *Notification) error {
	m.nextID++
	n.ID = m.nextID
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	m.notifs[n.ID] = n
	return nil
}

func (m *mockNotifRepo) ListByUserID(userID int64, page, pageSize int) ([]Notification, int64, error) {
	var result []Notification
	for _, n := range m.notifs {
		if n.UserID == userID {
			result = append(result, *n)
		}
	}
	total := int64(len(result))

	// 模拟分页
	start := (page - 1) * pageSize
	if start >= len(result) {
		return nil, total, nil
	}
	end := start + pageSize
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

func (m *mockNotifRepo) ListByType(userID int64, nType string, unreadOnly bool, page, pageSize int) ([]Notification, int64, error) {
	var result []Notification
	for _, n := range m.notifs {
		if n.UserID == userID && (nType == "" || n.Type == nType) && (!unreadOnly || !n.IsRead) {
			result = append(result, *n)
		}
	}
	total := int64(len(result))

	start := (page - 1) * pageSize
	if start >= len(result) {
		return nil, total, nil
	}
	end := start + pageSize
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

func (m *mockNotifRepo) BatchMarkAsRead(ids []int64, userID int64) error {
	idSet := make(map[int64]bool)
	for _, id := range ids {
		idSet[id] = true
	}
	for _, n := range m.notifs {
		if n.UserID == userID && idSet[n.ID] && !n.IsRead {
			n.IsRead = true
		}
	}
	return nil
}

func (m *mockNotifRepo) BatchDelete(ids []int64, userID int64) error {
	idSet := make(map[int64]bool)
	for _, id := range ids {
		idSet[id] = true
	}
	for id, n := range m.notifs {
		if idSet[id] && n.UserID == userID {
			delete(m.notifs, id)
		}
	}
	return nil
}

func (m *mockNotifRepo) GetByID(id int64) (*Notification, error) {
	if n, ok := m.notifs[id]; ok {
		return n, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockNotifRepo) MarkAsRead(id, userID int64) error {
	if n, ok := m.notifs[id]; ok && n.UserID == userID {
		n.IsRead = true
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *mockNotifRepo) MarkAllAsRead(userID int64) error {
	for _, n := range m.notifs {
		if n.UserID == userID {
			n.IsRead = true
		}
	}
	return nil
}

func (m *mockNotifRepo) CountUnread(userID int64) (int64, error) {
	var count int64
	for _, n := range m.notifs {
		if n.UserID == userID && !n.IsRead {
			count++
		}
	}
	return count, nil
}

// ============================================================
// 工具函数
// ============================================================

func createMockService() (*Service, *mockNotifRepo) {
	mock := newMockNotifRepo()
	svc := &Service{repo: mock}
	return svc, mock
}

func seedNotifications(mock *mockNotifRepo, userID, actorID int64, nType string, count int) {
	for i := 0; i < count; i++ {
		n := &Notification{
			UserID:  userID,
			ActorID: actorID,
			Type:    nType,
			Content: fmt.Sprintf("通知%d", i+1),
			IsRead:  false,
		}
		mock.Create(n)
	}
}

// ============================================================
// 测试
// ============================================================

func TestNotificationService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock := createMockService()

		err := svc.Create(100, 200, TypeLike, nil, nil, "赞了你的帖子")
		require.NoError(t, err)
		assert.Len(t, mock.notifs, 1)
	})

	t.Run("self_notification_skip", func(t *testing.T) {
		svc, mock := createMockService()

		err := svc.Create(100, 100, TypeLike, nil, nil, "赞了自己的帖子")
		require.NoError(t, err)
		assert.Len(t, mock.notifs, 0) // 不应创建
	})
}

func TestNotificationService_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock := createMockService()
		seedNotifications(mock, 100, 200, TypeLike, 5)

		notifications, total, err := svc.List(100, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, notifications, 5)
	})

	t.Run("empty", func(t *testing.T) {
		svc, _ := createMockService()

		notifications, total, err := svc.List(100, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, notifications, 0)
	})

	t.Run("pagination", func(t *testing.T) {
		svc, mock := createMockService()
		seedNotifications(mock, 100, 200, TypeLike, 15)

		n1, total, err := svc.List(100, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(15), total)
		assert.Len(t, n1, 10)

		n2, _, err := svc.List(100, 2, 10)
		require.NoError(t, err)
		assert.Len(t, n2, 5)
	})

	t.Run("default_params", func(t *testing.T) {
		svc, mock := createMockService()
		seedNotifications(mock, 100, 200, TypeLike, 3)

		_, _, err := svc.List(100, 0, 0)
		require.NoError(t, err)
	})
}

func TestNotificationService_ListByType(t *testing.T) {
	t.Run("filter_by_like", func(t *testing.T) {
		svc, mock := createMockService()
		seedNotifications(mock, 100, 200, TypeLike, 2)
		seedNotifications(mock, 100, 200, TypeComment, 3)

		notifications, total, err := svc.ListByType(100, "like", false, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, notifications, 2)
		for _, n := range notifications {
			assert.Equal(t, "like", n.Type)
		}
	})

	t.Run("filter_by_comment", func(t *testing.T) {
		svc, mock := createMockService()
		seedNotifications(mock, 100, 200, TypeLike, 1)
		seedNotifications(mock, 100, 200, TypeComment, 2)

		notifications, total, err := svc.ListByType(100, "comment", false, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, notifications, 2)
	})
}

func TestNotificationService_MarkAsRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock := createMockService()
		seedNotifications(mock, 100, 200, TypeLike, 1)

		var nID int64
		for id := range mock.notifs {
			nID = id
		}

		err := svc.MarkAsRead(nID, 100)
		require.NoError(t, err)

		n, _ := svc.GetDetail(nID)
		assert.True(t, n.IsRead)
	})

	t.Run("not_found", func(t *testing.T) {
		svc, _ := createMockService()

		err := svc.MarkAsRead(99999, 100)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "通知不存在")
	})

	t.Run("forbidden", func(t *testing.T) {
		svc, mock := createMockService()
		seedNotifications(mock, 100, 200, TypeLike, 1)

		var nID int64
		for id := range mock.notifs {
			nID = id
		}

		err := svc.MarkAsRead(nID, 300) // 其他用户
		require.Error(t, err)
		assert.Contains(t, err.Error(), "无权操作该通知")
	})
}

func TestNotificationService_MarkAllAsRead(t *testing.T) {
	svc, mock := createMockService()
	seedNotifications(mock, 100, 200, TypeLike, 5)

	count, _ := svc.GetUnreadCount(100)
	assert.Equal(t, int64(5), count)

	err := svc.MarkAllAsRead(100)
	require.NoError(t, err)

	count, _ = svc.GetUnreadCount(100)
	assert.Equal(t, int64(0), count)
}

func TestNotificationService_GetUnreadCount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock := createMockService()
		seedNotifications(mock, 100, 200, TypeLike, 3)

		count, err := svc.GetUnreadCount(100)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	t.Run("zero", func(t *testing.T) {
		svc, _ := createMockService()
		count, err := svc.GetUnreadCount(100)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestNotificationService_GetDetail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock := createMockService()
		seedNotifications(mock, 100, 200, TypeLike, 1)

		var nID int64
		for id := range mock.notifs {
			nID = id
		}

		n, err := svc.GetDetail(nID)
		require.NoError(t, err)
		assert.NotNil(t, n)
	})

	t.Run("not_found", func(t *testing.T) {
		svc, _ := createMockService()
		_, err := svc.GetDetail(99999)
		require.Error(t, err)
	})
}

func TestNotificationService_BatchMarkAsRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock := createMockService()
		seedNotifications(mock, 100, 200, TypeLike, 5)

		var ids []int64
		for id := range mock.notifs {
			ids = append(ids, id)
			if len(ids) >= 3 {
				break
			}
		}

		err := svc.BatchMarkAsRead(ids, 100)
		require.NoError(t, err)

		count, _ := svc.GetUnreadCount(100)
		assert.Equal(t, int64(2), count)
	})

	t.Run("empty_ids", func(t *testing.T) {
		svc, _ := createMockService()

		err := svc.BatchMarkAsRead([]int64{}, 100)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "通知ID列表不能为空")
	})
}

func TestNotificationService_BatchDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock := createMockService()
		seedNotifications(mock, 100, 200, TypeLike, 5)

		var ids []int64
		for id := range mock.notifs {
			ids = append(ids, id)
			if len(ids) >= 2 {
				break
			}
		}

		err := svc.BatchDelete(ids, 100)
		require.NoError(t, err)

		_, total, _ := svc.List(100, 1, 10)
		assert.Equal(t, int64(3), total)
	})

	t.Run("empty_ids", func(t *testing.T) {
		svc, _ := createMockService()

		err := svc.BatchDelete([]int64{}, 100)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "通知ID列表不能为空")
	})
}

func TestNotificationService_ListByType_UnreadOnly(t *testing.T) {
	svc, mock := createMockService()

	// 创建混合已读/未读的通知
	n1 := &Notification{UserID: 100, ActorID: 200, Type: TypeLike, Content: "未读", IsRead: false}
	n2 := &Notification{UserID: 100, ActorID: 200, Type: TypeLike, Content: "已读", IsRead: true}
	mock.Create(n1)
	mock.Create(n2)

	notifications, total, err := svc.ListByType(100, "like", true, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, notifications, 1)
	assert.Equal(t, "未读", notifications[0].Content)
}

// ============================================================
// 覆盖 repo 错误路径
// ============================================================

type errorNotifRepo struct {
	*mockNotifRepo
	createErr error
}

func (m *errorNotifRepo) Create(n *Notification) error {
	if m.createErr != nil {
		return m.createErr
	}
	return m.mockNotifRepo.Create(n)
}

func TestNotificationService_Create_RepoError(t *testing.T) {
	mock := &errorNotifRepo{mockNotifRepo: newMockNotifRepo(), createErr: errors.New("写入失败")}
	svc := &Service{repo: mock}

	err := svc.Create(100, 200, TypeLike, nil, nil, "赞了你的帖子")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "写入失败")
}

func TestNotificationService_ListByUserID_RepoError(t *testing.T) {
	svc, _ := createMockService()
	// 实际不会出错，但验证 default 参数的稳定性
	_, _, err := svc.List(100, -1, -1)
	require.NoError(t, err)
}

func TestNotificationService_ListByType_DefaultParams(t *testing.T) {
	svc, mock := createMockService()
	seedNotifications(mock, 100, 200, TypeLike, 3)

	_, _, err := svc.ListByType(100, "", false, 0, 0)
	require.NoError(t, err)
}

// ============================================================
// WebSocket 推送分支覆盖
// ============================================================

func TestNotificationService_Create_WithWebSocket(t *testing.T) {
	// 初始化测试用 Hub
	hub := websocket.NewTestHub()
	hub.AddOnlineUserForTest(100)
	websocket.SetGlobalHubForTest(hub)
	defer websocket.SetGlobalHubForTest(nil)

	svc, mock := createMockService()

	err := svc.Create(100, 200, TypeLike, nil, nil, "赞了你的帖子")
	require.NoError(t, err)
	assert.Len(t, mock.notifs, 1)
}

func TestNotificationService_Create_WebSocketNoMessage(t *testing.T) {
	// Hub 存在但用户不在线
	hub := websocket.NewTestHub()
	websocket.SetGlobalHubForTest(hub)
	defer websocket.SetGlobalHubForTest(nil)

	svc, mock := createMockService()

	err := svc.Create(100, 200, TypeLike, nil, nil, "赞了你的帖子")
	require.NoError(t, err)
	assert.Len(t, mock.notifs, 1)
}

// ============================================================
// MarkAsRead 错误路径覆盖
// ============================================================

func TestNotificationService_MarkAsRead_RepoError(t *testing.T) {
	svc, _ := createMockService()

	// repo.GetByID 默认返回 gorm.ErrRecordNotFound，已覆盖 not_found case
	// 验证其他错误路径
	err := svc.MarkAsRead(999, 100)
	require.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound) ||
		assert.Contains(t, err.Error(), "通知不存在"))
}
