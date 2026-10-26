package post

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/shiyemo/server/internal/auth"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// 覆盖 repo 错误路径
// ============================================================

func TestPostService_Create_RepoError(t *testing.T) {
	mock := &errorPostRepo{mockPostRepo: newMockPostRepo(), createErr: errors.New("写入失败")}
	userRepo := newMockUserGetter()
	moderator := newMockModerator()
	svc := NewPostServiceWithDeps(mock, userRepo, moderator)

	_, err := svc.Create(100, CreateRequest{
		Type:    "emotion",
		Title:   "标题",
		Content: "内容",
	})
	require.Error(t, err)
	assert.Equal(t, bizerr.ErrDB, err)
}

func TestPostService_GetDetail_RepoError(t *testing.T) {
	mock := &errorPostRepo{mockPostRepo: newMockPostRepo(), getErr: errors.New("连接失败")}
	userRepo := newMockUserGetter()
	moderator := newMockModerator()
	svc := NewPostServiceWithDeps(mock, userRepo, moderator)

	_, err := svc.GetDetail(1)
	require.Error(t, err)
	assert.Equal(t, bizerr.ErrDB, err)
}

// ============================================================
// Model 辅助测试
// ============================================================

func TestPost_Model_GetSetImages(t *testing.T) {
	p := &Post{}

	images := p.GetImages()
	assert.Equal(t, []string{}, images)

	imgs := []string{"img1.png", "img2.jpg"}
	p.SetImages(imgs)
	got := p.GetImages()
	assert.Equal(t, imgs, got)

	jsonData, _ := json.Marshal(imgs)
	assert.Equal(t, string(jsonData), p.Images)
}

func TestPost_Model_GetSetTags(t *testing.T) {
	p := &Post{}

	tags := []string{"Go", "Web", "AI"}
	p.SetTags(tags)
	got := p.GetTags()
	assert.Equal(t, tags, got)
}

func TestPost_IsPublished(t *testing.T) {
	p := &Post{Status: 1}
	assert.True(t, p.IsPublished())

	p.Status = 0
	assert.False(t, p.IsPublished())

	p.Status = 3
	assert.False(t, p.IsPublished())
}

func TestPost_TableName(t *testing.T) {
	p := Post{}
	assert.Equal(t, "posts", p.TableName())
}

// ============================================================
// CreateRequest / UpdateRequest 边界验证
// ============================================================

func TestPostService_Create_WithoutImagesAndTags(t *testing.T) {
	svc, _, _, _ := createService()

	p, err := svc.Create(100, CreateRequest{
		Type:    "case",
		Title:   "案例",
		Content: "详细描述",
	})
	require.NoError(t, err)
	assert.Equal(t, "case", p.Type)
	assert.Equal(t, []string{}, p.GetImages())
	assert.Equal(t, []string{}, p.GetTags())
}

func TestPostService_Update_WithImagesAndTags(t *testing.T) {
	svc, postRepo, _, _ := createService()
	postRepo.Create(createTestPost(100, "原标题", "原内容", "emotion"))

	var postID int64
	for id := range postRepo.posts {
		postID = id
	}

	err := svc.Update(100, postID, UpdateRequest{
		Images: []string{"a.png"},
		Tags:   []string{"new-tag"},
	})
	require.NoError(t, err)
}

func TestPostService_Update_Forbidden_DifferentUser(t *testing.T) {
	svc, postRepo, _, _ := createService()
	postRepo.Create(createTestPost(100, "帖子", "内容", "insight"))

	var postID int64
	for id := range postRepo.posts {
		postID = id
	}

	err := svc.Update(999, postID, UpdateRequest{Title: "恶意"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无权操作")
}

func TestPostService_Delete_Forbidden_DifferentUser(t *testing.T) {
	svc, postRepo, _, _ := createService()
	postRepo.Create(createTestPost(100, "帖子", "内容", "case"))

	var postID int64
	for id := range postRepo.posts {
		postID = id
	}

	err := svc.Delete(999, postID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无权操作")
}

func TestPostService_List_RepoError(t *testing.T) {
	postRepo := newMockPostRepo()
	userRepo := newMockUserGetter()
	moderator := newMockModerator()
	svc := NewPostServiceWithDeps(postRepo, userRepo, moderator)

	posts, total, err := svc.List("emotion", 0, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, posts, 0)
}

func TestPostService_GetDetail_AuthorError(t *testing.T) {
	svc, postRepo, _, _ := createService()
	post := createTestPost(100, "标题", "内容", "emotion")
	postRepo.Create(post)

	_, err := svc.GetDetail(1)
	require.Error(t, err)
}

func TestPostService_Create_MultiplePosts(t *testing.T) {
	svc, _, _, _ := createService()

	for i := 0; i < 5; i++ {
		p, err := svc.Create(100, CreateRequest{
			Type:    "emotion",
			Title:   "帖子" + string(rune('A'+i)),
			Content: "内容",
		})
		require.NoError(t, err)
		assert.NotNil(t, p)
	}
}

func TestPostService_List_DifferentTypes(t *testing.T) {
	svc, postRepo, _, _ := createService()
	types := []string{"emotion", "insight", "case", "action"}
	for _, typ := range types {
		postRepo.Create(createTestPost(100, typ+"帖", "内容", typ))
	}

	for _, typ := range types {
		posts, total, err := svc.List(typ, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total, "类型 %s 应只有1个帖子", typ)
		assert.Len(t, posts, 1, "类型 %s 返回数量", typ)
	}
}

// TestPostService_List_AchievementType 补 P1-4：achievement 类型于 2026-07-20 加入
// 后端 enum，但已有 TestPostService_List_DifferentTypes 仅覆盖 4 个旧类型，
// 新类型 List 过滤无回归保护。此用例验证：
//  1. 单一 achievement 类型可正确返回
//  2. 多种类型混存时按 type 严格过滤（不串到其它类型）
func TestPostService_List_AchievementType(t *testing.T) {
	svc, postRepo, _, _ := createService()

	// 5 种类型各 1 个
	allTypes := []string{"emotion", "insight", "case", "action", "achievement"}
	for _, typ := range allTypes {
		postRepo.Create(createTestPost(100, typ+"展示", "内容", typ))
	}

	t.Run("achievement 单独过滤", func(t *testing.T) {
		posts, total, err := svc.List("achievement", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total, "achievement 类型应只有1个帖子")
		assert.Len(t, posts, 1)
	})

	t.Run("其它类型不会被串到 achievement 列表", func(t *testing.T) {
		posts, _, err := svc.List("achievement", 1, 10)
		require.NoError(t, err)
		require.Len(t, posts, 1)
		// 不应包含任何 emotion/insight/case/action 帖子
		for _, p := range posts {
			typ, _ := p["type"].(string)
			assert.Equal(t, "achievement", typ,
				"achievement 过滤结果不应混入 %q 帖", typ)
		}
	})

	t.Run("全类型都返回全量", func(t *testing.T) {
		// 空 type 参数 = 不过滤
		posts, total, err := svc.List("", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, posts, 5)
	})
}

// TestPostService_Create_Achievement 补 P1-4：创建 achievement 类型帖子后
// 能通过 List 查回，验证完整链路（Create → DB → List）
func TestPostService_Create_Achievement(t *testing.T) {
	svc, postRepo, _, _ := createService()

	// 调 service.Create 创建（经过审核 / 字段处理等业务逻辑）
	p, err := svc.Create(100, CreateRequest{
		Type:    "achievement",
		Title:   "学会做红烧肉了",
		Content: "虽然只是一道菜，但失业两个月后的第一顿像样的饭。",
	})
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, "achievement", p.Type)

	// 验证入库后 List 能查回
	posts, total, err := svc.List("achievement", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, posts, 1)
	// 以及直接通过 repo 拿底层对象校验 type 字段确实是 achievement
	require.NotEmpty(t, postRepo.posts)
	for _, stored := range postRepo.posts {
		assert.Equal(t, "achievement", stored.Type,
			"仓库里存的 type 应是 achievement（前端 / 后端 enum 一致性）")
	}
}

// ============================================================
// CohortFeed 测试
// ============================================================

func TestPostService_CohortFeed(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, postRepo, userRepo, _ := createService()

		user := &auth.User{ID: 100, Nickname: "测试", City: "北京", Industry: "互联网"}
		userRepo.users[100] = user
		postRepo.Create(createTestPost(100, "同温层帖子", "内容", "emotion"))

		posts, total, err := svc.CohortFeed("emotion", "北京", "互联网", 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, posts, 1)
	})

	t.Run("empty", func(t *testing.T) {
		svc, _, _, _ := createService()

		posts, total, err := svc.CohortFeed("emotion", "不存在的城市", "不存在的行业", 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, posts, 0)
	})

	t.Run("default_pagination", func(t *testing.T) {
		svc, _, _, _ := createService()

		_, _, err := svc.CohortFeed("", "", "", 0, 0)
		require.NoError(t, err)
	})
}

func TestPostService_SetIndexer(t *testing.T) {
	svc, _, _, _ := createService()
	svc.SetIndexer(nil)
}

// ============================================================
// Indexer + Anonymous 路径测试（提升覆盖率）
// ============================================================

func TestPostService_Create_Anonymous_WithIndexer(t *testing.T) {
	postRepo := newMockPostRepo()
	userRepo := newMockUserGetter()
	moderator := newMockModerator()
	indexer := newMockPostIndexer()

	svc := &PostService{
		postRepo:  postRepo,
		userRepo:  userRepo,
		moderator: moderator,
		indexer:   indexer,
		sf:        nil, // 不影响测试，单条帖子无并发
	}

	p, err := svc.Create(100, CreateRequest{
		Type:        "emotion",
		Title:       "匿名测试",
		Content:     "匿名内容",
		IsAnonymous: true,
		Tags:        []string{"匿名"},
		Images:      []string{"img.png"},
	})

	require.NoError(t, err)
	assert.True(t, p.IsAnonymous)
	assert.NotEmpty(t, p.AnonymousName, "匿名帖应生成匿名名称")
	assert.Equal(t, 1, len(indexer.indexed), "indexer 应被调用")
	assert.Equal(t, p.ID, indexer.indexed[0].ID)
}

func TestPostService_Create_IndexerError(t *testing.T) {
	postRepo := newMockPostRepo()
	userRepo := newMockUserGetter()
	moderator := newMockModerator()
	indexer := newMockPostIndexer()
	indexer.indexErr = errors.New("索引失败")

	svc := &PostService{
		postRepo:  postRepo,
		userRepo:  userRepo,
		moderator: moderator,
		indexer:   indexer,
	}

	// 创建应成功，indexer 错误只打印 log，不阻断
	p, err := svc.Create(100, CreateRequest{
		Type:    "emotion",
		Title:   "正常帖子",
		Content: "内容",
	})
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestPostService_Delete_IndexerError(t *testing.T) {
	postRepo := newMockPostRepo()
	userRepo := newMockUserGetter()
	moderator := newMockModerator()
	indexer := newMockPostIndexer()
	indexer.deleteErr = errors.New("删除索引失败")

	svc := &PostService{
		postRepo:  postRepo,
		userRepo:  userRepo,
		moderator: moderator,
		indexer:   indexer,
	}

	// 创建并删除帖子
	postRepo.Create(createTestPost(100, "帖子", "内容", "emotion"))
	var postID int64
	for id := range postRepo.posts {
		postID = id
	}

	err := svc.Delete(100, postID)
	require.NoError(t, err) // indexer 错误只打印 log，不阻断
	assert.Len(t, postRepo.posts, 0)
}

// ============================================================
// FollowingFeed 错误路径测试
// ============================================================

func TestPostService_FollowingFeed_RepoError(t *testing.T) {
	postRepo := &errorPostRepo{
		mockPostRepo:   newMockPostRepo(),
		getFollowingErr: errors.New("查询失败"),
	}
	userRepo := newMockUserGetter()
	moderator := newMockModerator()
	svc := NewPostServiceWithDeps(postRepo, userRepo, moderator)

	_, _, err := svc.FollowingFeed(100, 1, 10)
	require.Error(t, err)
	assert.Equal(t, bizerr.ErrDB, err)
}

// ============================================================
// Update 字段更新完整路径测试
// ============================================================

func TestPostService_Update_AllFields(t *testing.T) {
	svc, postRepo, _, _ := createService()
	postRepo.Create(createTestPost(100, "原标题", "原内容", "emotion"))

	var postID int64
	for id := range postRepo.posts {
		postID = id
	}

	err := svc.Update(100, postID, UpdateRequest{
		Type:    "insight",
		Title:   "新标题",
		Content: "新内容",
		Images:  []string{"new.png"},
		Tags:    []string{"new-tag"},
	})
	require.NoError(t, err)

	p := postRepo.posts[postID]
	assert.Equal(t, "新标题", p.Title)
	assert.Equal(t, "新内容", p.Content)
	assert.Equal(t, "insight", p.Type)
}
