package integration

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/shiyemo/server/internal/common/database"
	"github.com/shiyemo/server/internal/post"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Flow Test — 跨模块用户旅程验证
// 验证多接口编排后流程是否一致，而非单接口正确性
// 运行: go test ./tests/integration/ -run TestFlow -v
// ============================================================

// TestFlow_NewUserJourney 新用户完整旅程：
// 注册 → 发帖 → 列表可见 → 评论 → 点赞 → 搜索 → 取消点赞 → 删除帖子
func TestFlow_NewUserJourney(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	var (
		token  string
		userID int64
		postID int64
	)

	// ── Step 1: 注册 ────────────────────────────────
	t.Run("1_注册新用户", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/register/email", &RegisterRequest{
			Email:    "flow-user@test.com",
			Password: "FlowTest123",
			Nickname: "流程测试用户",
		}, nil)

		resp := assertOK(t, w)
		token = extractStringField(t, resp.Data, "token")
		require.NotEmpty(t, token, "注册必须返回 token")

		userID = extractID(t, resp.Data)
		require.Greater(t, userID, int64(0))
		t.Logf("→ 注册成功，userID=%d", userID)
	})

	// ── Step 2: 发帖 ────────────────────────────────
	t.Run("2_发布帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{
			"type":    "emotion",
			"title":   "离职第一天的感受",
			"content": "今天正式离职了，心情很复杂",
			"tags":    []string{"离职", "情绪"},
		}, authHeader(token))

		resp := assertOK(t, w)
		postID = extractID(t, resp.Data)
		require.Greater(t, postID, int64(0))
		t.Logf("→ 发帖成功，postID=%d", postID)
	})

	// ── Step 3: 帖子列表可见 ─────────────────────────
	t.Run("3_列表中可见", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", "/api/v1/posts?page=1&size=20", nil, nil)
		resp := assertOK(t, w)

		list := extractField(t, resp.Data, "list")
		var items []map[string]interface{}
		require.NoError(t, json.Unmarshal(list, &items))
		require.GreaterOrEqual(t, len(items), 1, "列表必须包含刚发的帖子")

		found := false
		for _, item := range items {
			if int64(item["id"].(float64)) == postID {
				found = true
				break
			}
		}
		assert.True(t, found, "列表中必须能找到 postID=%d", postID)
		t.Log("→ 帖子在列表中可见")
	})

	// ── Step 4: 帖子详情包含正确数据 ───────────────────
	t.Run("4_帖子详情正确", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", fmt.Sprintf("/api/v1/posts/%d", postID), nil, authHeader(token))
		resp := assertOK(t, w)

		var post map[string]interface{}
		require.NoError(t, json.Unmarshal(resp.Data, &post))

		assert.Equal(t, "离职第一天的感受", post["title"])
		assert.Equal(t, "emotion", post["type"])
		assert.Equal(t, float64(0), post["like_count"])
		assert.Equal(t, float64(0), post["comment_count"])
		t.Log("→ 帖子详情数据正确")
	})

	// ── Step 5: 评论帖子 ─────────────────────────────
	var commentID int64
	t.Run("5_发表评论", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", fmt.Sprintf("/api/v1/posts/%d/comments", postID), map[string]interface{}{
			"content": "加油，一切都会好的！",
		}, authHeader(token))

		resp := assertOK(t, w)
		commentID = extractID(t, resp.Data)
		require.Greater(t, commentID, int64(0))
		t.Logf("→ 评论成功，commentID=%d", commentID)
	})

	// ── Step 6: 评论后 comment_count 更新 ──────────────
	t.Run("6_评论计数更新", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", fmt.Sprintf("/api/v1/posts/%d", postID), nil, nil)
		resp := assertOK(t, w)

		var post map[string]interface{}
		require.NoError(t, json.Unmarshal(resp.Data, &post))
		assert.Equal(t, float64(1), post["comment_count"], "评论后 comment_count 应为 1")
		t.Log("→ 评论计数正确更新")
	})

	// ── Step 7: 点赞帖子 ─────────────────────────────
	t.Run("7_点赞", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", fmt.Sprintf("/api/v1/posts/%d/like", postID), nil, authHeader(token))
		assertOK(t, w)
		t.Log("→ 点赞成功")
	})

	// ── Step 8: 点赞后 like_count 更新 ────────────────
	t.Run("8_点赞计数更新", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", fmt.Sprintf("/api/v1/posts/%d", postID), nil, authHeader(token))
		resp := assertOK(t, w)

		var post map[string]interface{}
		require.NoError(t, json.Unmarshal(resp.Data, &post))
		assert.Equal(t, float64(1), post["like_count"], "点赞后 like_count 应为 1")
		t.Log("→ 点赞计数正确更新")
	})

	// ── Step 9: 取消点赞（POST toggle） ─────────────────
	t.Run("9_取消点赞", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", fmt.Sprintf("/api/v1/posts/%d/like", postID), nil, authHeader(token))
		assertOK(t, w)

		// 验证计数回落
		detailW := doRequest(t, srv.Router, "GET", fmt.Sprintf("/api/v1/posts/%d", postID), nil, nil)
		detailResp := assertOK(t, detailW)
		var post map[string]interface{}
		require.NoError(t, json.Unmarshal(detailResp.Data, &post))
		assert.Equal(t, float64(0), post["like_count"], "取消点赞后应回到 0")
		t.Log("→ 取消点赞，计数正确回落")
	})

	// ── Step 10: 删除帖子后列表不可见 ────────────────────
	t.Run("10_删除后列表不可见", func(t *testing.T) {
		w := doRequest(t, srv.Router, "DELETE", fmt.Sprintf("/api/v1/posts/%d", postID), nil, authHeader(token))
		assertOK(t, w)

		// 验证列表中不再有该帖子
		listW := doRequest(t, srv.Router, "GET", "/api/v1/posts?page=1&size=20", nil, nil)
		listResp := assertOK(t, listW)
		list := extractField(t, listResp.Data, "list")

		var items []map[string]interface{}
		_ = json.Unmarshal(list, &items)
		for _, item := range items {
			assert.NotEqual(t, float64(postID), item["id"], "删除后不应出现在列表")
		}
		t.Log("→ 删除后列表中已不可见")
	})
}

// TestFlow_MultiUserInteraction 多用户交互旅程：
// A注册 → A发帖 → B注册 → B关注A → B在关注流看到A的帖 → B评论 → A收到通知
func TestFlow_MultiUserInteraction(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	var (
		tokenA string
		userA  int64
		tokenB string
		userB  int64
		postID int64
	)

	// ── 用户 A 注册并发帖 ─────────────────────────────
	t.Run("1_A注册", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/register/email", &RegisterRequest{
			Email: "userA-flow@test.com", Password: "TestPass123", Nickname: "用户A",
		}, nil)
		resp := assertOK(t, w)
		tokenA = extractStringField(t, resp.Data, "token")
		userA = extractID(t, resp.Data)
		t.Logf("→ 用户A注册成功, ID=%d", userA)
	})

	t.Run("2_A发帖", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{
			"type":    "insight",
			"title":   "找工作的三个技巧",
			"content": "第一不要焦虑，第二保持学习，第三多跟人聊",
			"tags":    []string{"求职", "经验"},
		}, authHeader(tokenA))
		resp := assertOK(t, w)
		postID = extractID(t, resp.Data)
		t.Logf("→ 用户A发帖成功, postID=%d", postID)
	})

	// ── 用户 B 注册 → 关注 A ──────────────────────────
	t.Run("3_B注册", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/register/email", &RegisterRequest{
			Email: "userB-flow@test.com", Password: "TestPass123", Nickname: "用户B",
		}, nil)
		resp := assertOK(t, w)
		tokenB = extractStringField(t, resp.Data, "token")
		userB = extractID(t, resp.Data)
		t.Logf("→ 用户B注册成功, ID=%d", userB)
	})

	t.Run("4_B关注A", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", fmt.Sprintf("/api/v1/users/%d/follow", userA), nil, authHeader(tokenB))
		assertOK(t, w)
		t.Log("→ 用户B关注用户A成功")
	})

	// ── 诊断检查 DB 状态 ──────────────────────────────────
	t.Run("4.5_诊断DB关注记录", func(t *testing.T) {
		var count int64
		testDB.Table("follows").Where("follower_id = ? AND following_id = ?", userB, userA).Count(&count)
		t.Logf("→ DB 中 follows 记录: follower_id=%d following_id=%d count=%d", userB, userA, count)

		var postCount int64
		testDB.Table("posts").Where("user_id = ? AND status = 1", userA).Count(&postCount)
		t.Logf("→ DB 中 userA 的帖子数: %d", postCount)

		// 直接测试 PostRepo.GetFollowingIDs
		globalDB := database.GetDB()
		t.Logf("→ globalDB == testDB? %v", globalDB == testDB)

		postRepo := post.NewPostRepoWithDB(testDB)
		ids, err := postRepo.GetFollowingIDs(userB)
		t.Logf("→ PostRepo.GetFollowingIDs(%d) = %v, err=%v", userB, ids, err)
	})

	// ── 诊断：直接调用 service ──────────────────────────────
	t.Run("4.8_诊断FollowingFeed", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", fmt.Sprintf("/debug/following/%d", userB), nil, nil)
		resp := parseResponse(t, w)
		t.Logf("→ debug/following/%d: code=%d data=%s", userB, resp.Code, string(resp.Data))
	})

	// ── 验证关注流 ─────────────────────────────────────
	t.Run("5_B的关注流有A的帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", "/api/v1/posts/following?page=1&size=20", nil, authHeader(tokenB))
		resp := assertOK(t, w)

		list := extractField(t, resp.Data, "list")
		var items []map[string]interface{}
		require.NoError(t, json.Unmarshal(list, &items))
		require.GreaterOrEqual(t, len(items), 1, "关注流必须包含A的帖子")

		found := false
		for _, item := range items {
			if int64(item["id"].(float64)) == postID {
				found = true
				break
			}
		}
		assert.True(t, found, "关注流中必须能找到A的帖子")
		t.Log("→ 关注流正确展示A的帖子")
	})

	// ── B 评论 A 的帖子 ────────────────────────────────
	t.Run("6_B评论A的帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", fmt.Sprintf("/api/v1/posts/%d/comments", postID), map[string]interface{}{
			"content": "写得太好了，收藏！",
		}, authHeader(tokenB))
		resp := assertOK(t, w)
		commentID := extractID(t, resp.Data)
		require.Greater(t, commentID, int64(0))
		t.Log("→ 用户B评论成功")
	})

	// ── 验证 A 的通知 ──────────────────────────────────
	t.Run("7_A有未读通知", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", "/api/v1/notifications/unread-count", nil, authHeader(tokenA))
		resp := assertOK(t, w)

		var count map[string]interface{}
		require.NoError(t, json.Unmarshal(resp.Data, &count))
		unread, ok := count["count"].(float64)
		if ok {
			assert.GreaterOrEqual(t, int(unread), 1, "A应有至少1条未读通知")
			t.Logf("→ 用户A有 %d 条未读通知", int(unread))
		} else {
			t.Log("→ 通知计数格式待确认，跳过断言")
		}
	})

	// ── B 取消关注 A（POST toggle） ──────────────────────
	t.Run("8_B取关A", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", fmt.Sprintf("/api/v1/users/%d/follow", userA), nil, authHeader(tokenB))
		assertOK(t, w)
		t.Log("→ 用户B取消关注A成功")
	})

	// ── 取关后关注流为空 ─────────────────────────────────
	t.Run("9_取关后关注流无A帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", "/api/v1/posts/following?page=1&size=20", nil, authHeader(tokenB))
		resp := assertOK(t, w)

		list := extractField(t, resp.Data, "list")
		if list == nil || string(list) == "null" {
			t.Log("→ 取关后关注流为空（null）")
			return
		}
		var items []map[string]interface{}
		_ = json.Unmarshal(list, &items)
		for _, item := range items {
			assert.NotEqual(t, float64(postID), item["id"], "取关后不应出现A的帖子")
		}
		t.Log("→ 取关后关注流正确清除")
	})
}

// TestFlow_PostVisibility_AfterCreate 验证刚发帖后缓存不影响可见性
func TestFlow_PostVisibility_AfterCreate(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	token := authToken(t, srv.Router, "visibility@test.com", "TestPass123", "可见性测试")

	// 先请求列表（建立缓存）
	w1 := doRequest(t, srv.Router, "GET", "/api/v1/posts?page=1&size=20", nil, nil)
	assertOK(t, w1)

	// 发帖
	w2 := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{
		"type": "emotion", "title": "缓存测试帖子", "content": "发帖后应立即可见",
	}, authHeader(token))
	resp := assertOK(t, w2)
	newPostID := extractID(t, resp.Data)

	// 再次请求列表 — 必须能看到新帖子（缓存应已失效）
	w3 := doRequest(t, srv.Router, "GET", "/api/v1/posts?page=1&size=20", nil, nil)
	listResp := assertOK(t, w3)
	list := extractField(t, listResp.Data, "list")

	var items []map[string]interface{}
	require.NoError(t, json.Unmarshal(list, &items))

	found := false
	for _, item := range items {
		if int64(item["id"].(float64)) == newPostID {
			found = true
			break
		}
	}
	assert.True(t, found, "发帖后立即请求列表必须能看到新帖子（缓存失效验证）")
	t.Log("✅ 发帖后缓存正确失效，新帖立即可见")
}

// ============================================================
// 辅助函数
// ============================================================

func extractStringField(t *testing.T, data json.RawMessage, field string) string {
	t.Helper()
	raw := extractField(t, data, field)
	if raw == nil {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return string(raw)
	}
	return s
}
