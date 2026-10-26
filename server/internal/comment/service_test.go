package comment

import (
	"testing"
	"time"

	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/post"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================
// Mock 实现
// ============================================================

type mockCommentRepo struct {
	comments  map[int64]*Comment
	nextID    int64
	deleteErr error
}

func newMockCommentRepo() *mockCommentRepo {
	return &mockCommentRepo{comments: make(map[int64]*Comment), nextID: 1}
}

func (m *mockCommentRepo) Create(c *Comment) error {
	c.ID = m.nextID
	m.nextID++
	c.CreatedAt = time.Now()
	m.comments[c.ID] = c
	return nil
}

func (m *mockCommentRepo) GetByID(id int64) (*Comment, error) {
	if c, ok := m.comments[id]; ok {
		return c, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockCommentRepo) GetByPostID(postID int64, page, size int) ([]Comment, int64, error) {
	var result []Comment
	for _, c := range m.comments {
		if c.PostID == postID && c.ParentID == nil && c.Status == 1 {
			result = append(result, *c)
		}
	}
	total := int64(len(result))
	start := (page - 1) * size
	if start >= len(result) {
		return nil, total, nil
	}
	end := start + size
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], total, nil
}

func (m *mockCommentRepo) GetReplies(parentID int64) ([]Comment, error) {
	var replies []Comment
	for _, c := range m.comments {
		if c.ParentID != nil && *c.ParentID == parentID && c.Status == 1 {
			replies = append(replies, *c)
		}
	}
	return replies, nil
}

func (m *mockCommentRepo) Delete(id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if c, ok := m.comments[id]; ok {
		c.Status = 3
	}
	return nil
}

func (m *mockCommentRepo) CountByPostID(postID int64) (int64, error) {
	var count int64
	for _, c := range m.comments {
		if c.PostID == postID && c.Status == 1 {
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
		1: {ID: 1, Nickname: "测试用户", Avatar: "avatar.png"},
	}}
}

func (m *mockUserProvider) GetByID(id int64) (*auth.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return &auth.User{ID: id, Nickname: "未知用户"}, nil
}

type mockPostProvider struct {
	posts map[int64]*post.Post
}

func newMockPostProvider() *mockPostProvider {
	return &mockPostProvider{posts: map[int64]*post.Post{
		1: {ID: 1, UserID: 1, Title: "测试帖子", Status: 1},
	}}
}

func (m *mockPostProvider) GetByID(id int64) (*post.Post, error) {
	if p, ok := m.posts[id]; ok {
		return p, nil
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

type mockCountUpdater struct {
	lastPostID int64
	lastCount  int64
}

func (m *mockCountUpdater) UpdateCommentCount(postID int64, count int64) {
	m.lastPostID = postID
	m.lastCount = count
}

// ============================================================
// 测试辅助
// ============================================================

type testDeps struct {
	repo      *mockCommentRepo
	userRepo  *mockUserProvider
	postRepo  *mockPostProvider
	notifySvc *mockNotifier
	countUpd  *mockCountUpdater
	service   *CommentService
}

func newTestDeps() *testDeps {
	repo := newMockCommentRepo()
	userRepo := newMockUserProvider()
	postRepo := newMockPostProvider()
	notifySvc := &mockNotifier{}
	countUpd := &mockCountUpdater{}
	service := NewCommentServiceWithDeps(repo, userRepo, postRepo, notifySvc, countUpd)
	return &testDeps{repo, userRepo, postRepo, notifySvc, countUpd, service}
}

// ============================================================
// Create 测试
// ============================================================

func TestCommentService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		d := newTestDeps()
		result, err := d.service.Create(1, 1, CreateRequest{Content: "测试评论"})

		require.NoError(t, err)
		assert.Equal(t, "测试评论", result["content"])
		assert.Equal(t, int64(1), result["post_id"])
		assert.Equal(t, int64(1), result["user_id"])
		assert.Nil(t, result["parent_id"])
	})

	t.Run("reply_success", func(t *testing.T) {
		d := newTestDeps()

		// 先创建主评论
		main, err := d.service.Create(1, 1, CreateRequest{Content: "主评论"})
		require.NoError(t, err)
		parentID := main["id"].(int64)

		// 创建回复
		reply, err := d.service.Create(1, 1, CreateRequest{Content: "回复", ParentID: &parentID})
		require.NoError(t, err)
		assert.Equal(t, "回复", reply["content"])
		assert.Equal(t, parentID, *reply["parent_id"].(*int64))
	})

	t.Run("post_not_found", func(t *testing.T) {
		d := newTestDeps()
		_, err := d.service.Create(1, 999, CreateRequest{Content: "评论"})

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrPostNotFound, err)
	})

	t.Run("parent_not_found", func(t *testing.T) {
		d := newTestDeps()
		parentID := int64(999)
		_, err := d.service.Create(1, 1, CreateRequest{Content: "回复", ParentID: &parentID})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "父评论不存在")
	})

	t.Run("parent_not_same_post", func(t *testing.T) {
		d := newTestDeps()

		// 创建主评论在帖子1
		main, _ := d.service.Create(1, 1, CreateRequest{Content: "主评论"})
		parentID := main["id"].(int64)

		// 尝试在帖子2下回复帖子1的评论 — mock post 2 也存在
		d.postRepo.posts[2] = &post.Post{ID: 2, UserID: 1, Title: "帖子2"}
		_, err := d.service.Create(1, 2, CreateRequest{Content: "跨帖回复", ParentID: &parentID})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "父评论不属于该帖子")
	})

	t.Run("max_two_levels", func(t *testing.T) {
		d := newTestDeps()

		// 主评论
		main, _ := d.service.Create(1, 1, CreateRequest{Content: "主评论"})
		parentID1 := main["id"].(int64)

		// 回复（二级）
		reply, _ := d.service.Create(1, 1, CreateRequest{Content: "回复", ParentID: &parentID1})
		parentID2 := reply["id"].(int64)

		// 回复的回复（三级 — 应失败）
		_, err := d.service.Create(1, 1, CreateRequest{Content: "三级", ParentID: &parentID2})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "最多支持两级评论")
	})

	t.Run("notification_sent", func(t *testing.T) {
		d := newTestDeps()
		d.service.Create(1, 1, CreateRequest{Content: "评论"})

		require.Len(t, d.notifySvc.calls, 1)
		assert.Equal(t, "comment", d.notifySvc.calls[0].NType)
		assert.Equal(t, int64(1), d.notifySvc.calls[0].UserID) // 帖子作者
	})

	t.Run("count_updated", func(t *testing.T) {
		d := newTestDeps()
		d.service.Create(1, 1, CreateRequest{Content: "评论"})

		assert.Equal(t, int64(1), d.countUpd.lastPostID)
		assert.Equal(t, int64(1), d.countUpd.lastCount)
	})
}

// ============================================================
// List 测试
// ============================================================

func TestCommentService_List(t *testing.T) {
	t.Run("success_with_replies", func(t *testing.T) {
		d := newTestDeps()

		// 创建主评论
		main, _ := d.service.Create(1, 1, CreateRequest{Content: "主评论"})
		parentID := main["id"].(int64)

		// 创建回复
		d.service.Create(1, 1, CreateRequest{Content: "回复", ParentID: &parentID})

		comments, total, err := d.service.List(1, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total) // 只有主评论计入 total
		assert.Len(t, comments, 1)

		replies := comments[0]["replies"].([]map[string]interface{})
		assert.Len(t, replies, 1)
		assert.Equal(t, "回复", replies[0]["content"])
	})

	t.Run("empty", func(t *testing.T) {
		d := newTestDeps()
		comments, total, err := d.service.List(1, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, comments, 0)
	})

	t.Run("pagination", func(t *testing.T) {
		d := newTestDeps()

		for i := 0; i < 5; i++ {
			d.service.Create(1, 1, CreateRequest{Content: "评论"})
		}

		// 第一页 2 条
		page1, total, _ := d.service.List(1, 1, 2)
		assert.Equal(t, int64(5), total)
		assert.Len(t, page1, 2)

		// 第三页 2 条 — 只剩 1 条
		page3, _, _ := d.service.List(1, 3, 2)
		assert.Len(t, page3, 1)
	})

	t.Run("author_info_present", func(t *testing.T) {
		d := newTestDeps()
		d.service.Create(1, 1, CreateRequest{Content: "评论"})

		comments, _, _ := d.service.List(1, 1, 10)
		author := comments[0]["author"].(map[string]interface{})
		assert.Equal(t, "测试用户", author["nickname"])
		assert.Equal(t, "avatar.png", author["avatar"])
	})
}

// ============================================================
// Delete 测试
// ============================================================

func TestCommentService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		d := newTestDeps()

		result, _ := d.service.Create(1, 1, CreateRequest{Content: "待删除"})
		commentID := result["id"].(int64)

		err := d.service.Delete(1, commentID)
		require.NoError(t, err)

		// 验证已标记删除
		assert.Equal(t, 3, d.repo.comments[commentID].Status)
		// 验证计数更新
		assert.Equal(t, int64(0), d.countUpd.lastCount)
	})

	t.Run("not_found", func(t *testing.T) {
		d := newTestDeps()
		err := d.service.Delete(1, 999)

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrPostNotFound, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		d := newTestDeps()

		result, _ := d.service.Create(1, 1, CreateRequest{Content: "我的评论"})
		commentID := result["id"].(int64)

		err := d.service.Delete(2, commentID) // 其他人
		require.Error(t, err)
		assert.Equal(t, bizerr.ErrPostForbidden, err)
	})

	t.Run("repo_delete_error", func(t *testing.T) {
		d := newTestDeps()

		result, _ := d.service.Create(1, 1, CreateRequest{Content: "评论"})
		commentID := result["id"].(int64)

		d.repo.deleteErr = assert.AnError
		err := d.service.Delete(1, commentID)

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrDB, err)
	})
}
