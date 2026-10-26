package interaction

import (
	"testing"
	"time"

	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/post"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Mock 实现
// ============================================================

type mockInteractionRepo struct {
	likes     map[string]bool // "userID:postID" -> true
	bookmarks map[string]bool
	nextID    int64
}

func newMockInteractionRepo() *mockInteractionRepo {
	return &mockInteractionRepo{
		likes:     make(map[string]bool),
		bookmarks: make(map[string]bool),
		nextID:    1,
	}
}

func key(userID, postID int64) string {
	return string(rune(userID)) + ":" + string(rune(postID))
}

func (m *mockInteractionRepo) Like(userID, postID int64) error {
	m.likes[key(userID, postID)] = true
	m.nextID++
	return nil
}

func (m *mockInteractionRepo) Unlike(userID, postID int64) error {
	delete(m.likes, key(userID, postID))
	return nil
}

func (m *mockInteractionRepo) IsLiked(userID, postID int64) bool {
	return m.likes[key(userID, postID)]
}

func (m *mockInteractionRepo) Bookmark(userID, postID int64) error {
	m.bookmarks[key(userID, postID)] = true
	m.nextID++
	return nil
}

func (m *mockInteractionRepo) Unbookmark(userID, postID int64) error {
	delete(m.bookmarks, key(userID, postID))
	return nil
}

func (m *mockInteractionRepo) IsBookmarked(userID, postID int64) bool {
	return m.bookmarks[key(userID, postID)]
}

func (m *mockInteractionRepo) CountLikes(postID int64) (int64, error) {
	var count int64
	for k := range m.likes {
		// check if key ends with ":postID"
		if len(k) > 0 && int64(k[len(k)-1]) == postID {
			count++
		}
	}
	return count, nil
}

func (m *mockInteractionRepo) ListBookmarks(userID int64, page, size int) ([]int64, int64, error) {
	var postIDs []int64
	for k, v := range m.bookmarks {
		if v && len(k) > 0 && int64(k[0]) == userID {
			postIDs = append(postIDs, int64(k[len(k)-1]))
		}
	}
	total := int64(len(postIDs))
	start := (page - 1) * size
	if start >= len(postIDs) {
		return nil, total, nil
	}
	end := start + size
	if end > len(postIDs) {
		end = len(postIDs)
	}
	return postIDs[start:end], total, nil
}

type mockPostProvider struct {
	posts map[int64]*post.Post
}

func newMockPostProvider() *mockPostProvider {
	return &mockPostProvider{posts: map[int64]*post.Post{
		1: {ID: 1, UserID: 1, Title: "测试帖子", Content: "内容", Type: "emotion", Status: 1, CreatedAt: time.Now()},
	}}
}

func (m *mockPostProvider) GetByID(id int64) (*post.Post, error) {
	if p, ok := m.posts[id]; ok {
		return p, nil
	}
	return nil, assert.AnError
}

type mockUserProvider struct {
	users map[int64]*auth.User
}

func newMockUserProvider() *mockUserProvider {
	return &mockUserProvider{users: map[int64]*auth.User{
		1: {ID: 1, Nickname: "测试用户", Avatar: "avatar.png"},
	}}
}

func (m *mockUserProvider) GetByID(id int64) (*auth.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return &auth.User{ID: id, Nickname: "未知用户"}, nil
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

type mockLikeCountUpdater struct {
	lastPostID int64
	lastCount  int64
}

func (m *mockLikeCountUpdater) UpdateLikeCount(postID int64, count int64) {
	m.lastPostID = postID
	m.lastCount = count
}

// ============================================================
// 测试辅助
// ============================================================

type testDeps struct {
	repo      *mockInteractionRepo
	postRepo  *mockPostProvider
	userRepo  *mockUserProvider
	notifySvc *mockNotifier
	countUpd  *mockLikeCountUpdater
	service   *InteractionService
}

func newTestDeps() *testDeps {
	repo := newMockInteractionRepo()
	postRepo := newMockPostProvider()
	userRepo := newMockUserProvider()
	notifySvc := &mockNotifier{}
	countUpd := &mockLikeCountUpdater{}
	service := NewInteractionServiceWithDeps(repo, postRepo, userRepo, notifySvc, countUpd)
	return &testDeps{repo, postRepo, userRepo, notifySvc, countUpd, service}
}

// ============================================================
// Like 测试
// ============================================================

func TestInteractionService_Like(t *testing.T) {
	t.Run("like_success", func(t *testing.T) {
		d := newTestDeps()
		liked, count, err := d.service.Like(1, 1)

		require.NoError(t, err)
		assert.True(t, liked)
		assert.Equal(t, int64(1), count) // 返回权威计数
	})

	t.Run("unlike_success", func(t *testing.T) {
		d := newTestDeps()

		// 先点赞
		d.service.Like(1, 1)
		// 再取消
		liked, count, err := d.service.Like(1, 1)

		require.NoError(t, err)
		assert.False(t, liked)
		assert.Equal(t, int64(0), count) // 取消后计数归零
	})

	t.Run("post_not_found", func(t *testing.T) {
		d := newTestDeps()
		_, _, err := d.service.Like(1, 999)

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrPostNotFound, err)
	})

	t.Run("notification_on_like", func(t *testing.T) {
		d := newTestDeps()
		d.service.Like(1, 1)

		require.Len(t, d.notifySvc.calls, 1)
		assert.Equal(t, "like", d.notifySvc.calls[0].NType)
		assert.Equal(t, int64(1), d.notifySvc.calls[0].UserID) // 帖子作者
	})

	t.Run("count_updated", func(t *testing.T) {
		d := newTestDeps()
		d.service.Like(1, 1)

		assert.Equal(t, int64(1), d.countUpd.lastPostID)
	})

	t.Run("no_notification_on_unlike", func(t *testing.T) {
		d := newTestDeps()

		d.service.Like(1, 1) // 点赞 — 产生通知
		d.service.Like(1, 1) // 取消 — 不产生通知

		assert.Len(t, d.notifySvc.calls, 1)
	})
}

// ============================================================
// Bookmark 测试
// ============================================================

func TestInteractionService_Bookmark(t *testing.T) {
	t.Run("bookmark_success", func(t *testing.T) {
		d := newTestDeps()
		bookmarked, err := d.service.Bookmark(1, 1)

		require.NoError(t, err)
		assert.True(t, bookmarked)
	})

	t.Run("unbookmark_success", func(t *testing.T) {
		d := newTestDeps()

		d.service.Bookmark(1, 1)
		bookmarked, err := d.service.Bookmark(1, 1)

		require.NoError(t, err)
		assert.False(t, bookmarked)
	})

	t.Run("post_not_found", func(t *testing.T) {
		d := newTestDeps()
		_, err := d.service.Bookmark(1, 999)

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrPostNotFound, err)
	})
}

// ============================================================
// ListBookmarks 测试
// ============================================================

func TestInteractionService_ListBookmarks(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		d := newTestDeps()

		d.service.Bookmark(1, 1)
		posts, total, err := d.service.ListBookmarks(1, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, posts, 1)
		assert.Equal(t, "测试帖子", posts[0]["title"])
	})

	t.Run("empty", func(t *testing.T) {
		d := newTestDeps()

		posts, total, err := d.service.ListBookmarks(1, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, posts, 0)
	})
}

// ============================================================
// GetUserPostState 测试
// ============================================================

func TestInteractionService_GetUserPostState(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		d := newTestDeps()
		liked, bookmarked := d.service.GetUserPostState(1, 1)
		assert.False(t, liked)
		assert.False(t, bookmarked)
	})

	t.Run("liked_and_bookmarked", func(t *testing.T) {
		d := newTestDeps()
		d.service.Like(1, 1)
		d.service.Bookmark(1, 1)

		liked, bookmarked := d.service.GetUserPostState(1, 1)
		assert.True(t, liked)
		assert.True(t, bookmarked)
	})

	t.Run("only_liked", func(t *testing.T) {
		d := newTestDeps()
		d.service.Like(1, 1)

		liked, bookmarked := d.service.GetUserPostState(1, 1)
		assert.True(t, liked)
		assert.False(t, bookmarked)
	})
}
