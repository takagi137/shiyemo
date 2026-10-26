package post

import (
	"testing"

	"github.com/shiyemo/server/internal/auth"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/moderation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, postRepo, _, _ := createService()

		p, err := svc.Create(100, CreateRequest{
			Type:    "emotion",
			Title:   "测试帖子",
			Content: "内容",
			Tags:    []string{"Go", "开发"},
		})

		require.NoError(t, err)
		assert.NotZero(t, p.ID)
		assert.Equal(t, "测试帖子", p.Title)
		assert.Len(t, postRepo.posts, 1)
	})

	t.Run("content_blocked", func(t *testing.T) {
		svc, _, _, moderator := createService()
		moderator.checkErr = &moderation.CheckResult{
			Passed:  false,
			Words:   []string{"违规词"},
			Message: "内容包含违规词汇",
		}

		_, err := svc.Create(100, CreateRequest{
			Type:    "emotion",
			Title:   "测试帖子",
			Content: "包含违规内容",
		})

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrContentBlocked, err)
	})
}

func TestPostService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, postRepo, _, _ := createService()
		postRepo.Create(createTestPost(100, "原标题", "原内容", "emotion"))

		var postID int64
		for id := range postRepo.posts {
			postID = id
		}

		err := svc.Update(100, postID, UpdateRequest{
			Title: "新标题",
		})

		require.NoError(t, err)
		assert.Equal(t, "新标题", postRepo.posts[postID].Title)
	})

	t.Run("not_found", func(t *testing.T) {
		svc, _, _, _ := createService()

		err := svc.Update(100, 99999, UpdateRequest{Title: "新标题"})
		require.Error(t, err)
		assert.Equal(t, bizerr.ErrPostNotFound, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		svc, postRepo, _, _ := createService()
		postRepo.Create(createTestPost(100, "标题", "内容", "emotion"))

		var postID int64
		for id := range postRepo.posts {
			postID = id
		}

		err := svc.Update(200, postID, UpdateRequest{Title: "恶意修改"})
		require.Error(t, err)
		assert.Equal(t, bizerr.ErrPostForbidden, err)
	})

	t.Run("empty_update_no_fields", func(t *testing.T) {
		svc, postRepo, _, _ := createService()
		postRepo.Create(createTestPost(100, "标题", "内容", "emotion"))

		var postID int64
		for id := range postRepo.posts {
			postID = id
		}

		err := svc.Update(100, postID, UpdateRequest{})
		require.NoError(t, err)
	})
}

func TestPostService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, postRepo, _, _ := createService()
		postRepo.Create(createTestPost(100, "待删帖子", "内容", "emotion"))

		var postID int64
		for id := range postRepo.posts {
			postID = id
		}

		err := svc.Delete(100, postID)
		require.NoError(t, err)
		assert.Len(t, postRepo.posts, 0)
	})

	t.Run("not_found", func(t *testing.T) {
		svc, _, _, _ := createService()

		err := svc.Delete(100, 99999)
		require.Error(t, err)
		assert.Equal(t, bizerr.ErrPostNotFound, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		svc, postRepo, _, _ := createService()
		postRepo.Create(createTestPost(100, "帖子", "内容", "emotion"))

		var postID int64
		for id := range postRepo.posts {
			postID = id
		}

		err := svc.Delete(200, postID)
		require.Error(t, err)
		assert.Equal(t, bizerr.ErrPostForbidden, err)
	})
}

func TestPostService_List(t *testing.T) {
	t.Run("with_type", func(t *testing.T) {
		svc, postRepo, _, _ := createService()
		postRepo.Create(createTestPost(100, "emotion帖", "内容", "emotion"))
		postRepo.Create(createTestPost(100, "insight帖", "内容", "insight"))

		posts, total, err := svc.List("emotion", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, posts, 1)
		assert.Equal(t, "emotion", posts[0]["type"])
	})

	t.Run("all_types", func(t *testing.T) {
		svc, postRepo, _, _ := createService()
		postRepo.Create(createTestPost(100, "帖1", "内容", "emotion"))
		postRepo.Create(createTestPost(100, "帖2", "内容", "insight"))

		posts, total, err := svc.List("", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, posts, 2)
	})

	t.Run("pagination", func(t *testing.T) {
		svc, postRepo, _, _ := createService()
		for i := 0; i < 15; i++ {
			postRepo.Create(createTestPost(100, "帖", "内容", "emotion"))
		}

		posts1, total, _ := svc.List("", 1, 10)
		assert.Equal(t, int64(15), total)
		assert.Len(t, posts1, 10)

		posts2, _, _ := svc.List("", 2, 10)
		assert.Len(t, posts2, 5)
	})
}

func TestPostService_GetDetail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, postRepo, userRepo, _ := createService()
		userRepo.users[100] = &auth.User{ID: 100, Nickname: "测试用户", Avatar: "avatar.png"}

		post := createTestPost(100, "标题", "内容", "emotion")
		postRepo.Create(post)

		result, err := svc.GetDetail(1)
		require.NoError(t, err)
		assert.Equal(t, "标题", result["title"])
		assert.NotNil(t, result["author"])
	})

	t.Run("not_found", func(t *testing.T) {
		svc, _, _, _ := createService()

		_, err := svc.GetDetail(99999)
		require.Error(t, err)
		assert.Equal(t, bizerr.ErrPostNotFound, err)
	})
}

func TestPostService_FollowingFeed(t *testing.T) {
	t.Run("empty_following", func(t *testing.T) {
		svc, _, _, _ := createService()

		posts, total, err := svc.FollowingFeed(100, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Len(t, posts, 0)
	})
}

func TestPostService_ListPostTypes(t *testing.T) {
	svc, postRepo, _, _ := createService()
	postRepo.Create(createTestPost(100, "帖1", "内容", "emotion"))
	postRepo.Create(createTestPost(100, "帖2", "内容", "insight"))
	postRepo.Create(createTestPost(100, "帖3", "内容", "case"))

	posts, _, err := svc.List("case", 1, 10)
	require.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, "case", posts[0]["type"])
}
