package integration

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Post 模块集成测试
// 验证：帖子 CRUD 端到端流程
// ============================================================

// createPost 便捷函数：创建帖子并返回帖子 ID
func createPost(t *testing.T, srv *TestServer, token, postType, title, content string, tags []string) int64 {
	t.Helper()
	w := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{
		"type":    postType,
		"title":   title,
		"content": content,
		"tags":    tags,
	}, authHeader(token))
	resp := assertOK(t, w)
	postID := extractID(t, resp.Data)
	assert.Greater(t, postID, int64(0), "创建帖子应返回正数 ID")
	return postID
}

// ============================================================
// 核心 CRUD 测试
// ============================================================

func TestPost_CRUD(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	// 创建测试用户
	token := authToken(t, srv.Router, "post-crud@test.com", "testpass123", "帖子CRUD用户")

	var postID int64

	// ── Create ──────────────────────────────────────
	t.Run("创建帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{
			"type":    "emotion",
			"title":   "今天的感悟",
			"content": "失业并不可怕，重要的是重新认识自己",
			"images":  []string{"http://example.com/img1.jpg"},
			"tags":    []string{"感悟", "失业"},
		}, authHeader(token))

		resp := assertOK(t, w)
		postID = extractID(t, resp.Data)
		assert.Greater(t, postID, int64(0), "帖子 ID 应为正数")

		// 验证返回的数据包含预期字段
		title := extractField(t, resp.Data, "title")
		var titleStr string
		_ = jsonUnmarshal(t, title, &titleStr)
		assert.Equal(t, "今天的感悟", titleStr)

		userID := extractField(t, resp.Data, "user_id")
		assert.NotNil(t, userID)
		t.Logf("✅ 创建帖子成功，ID: %d", postID)
	})

	// ── Get Detail ──────────────────────────────────
	t.Run("获取帖子详情", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", fmt.Sprintf("/api/v1/posts/%d", postID), nil, nil)
		resp := assertOK(t, w)

		// 验证标题
		title := extractField(t, resp.Data, "title")
		var titleStr string
		_ = jsonUnmarshal(t, title, &titleStr)
		assert.Equal(t, "今天的感悟", titleStr)

		// 验证内容
		content := extractField(t, resp.Data, "content")
		var contentStr string
		_ = jsonUnmarshal(t, content, &contentStr)
		assert.Equal(t, "失业并不可怕，重要的是重新认识自己", contentStr)

		// 验证作者信息
		author := extractField(t, resp.Data, "author")
		assert.NotNil(t, author, "详情应包含 author 字段")
		t.Log("✅ 获取帖子详情成功")
	})

	// ── List ───────────────────────────────────────
	t.Run("帖子列表", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", "/api/v1/posts?page=1&size=20", nil, nil)
		resp := assertOK(t, w)

		// 分页数据
		list := extractField(t, resp.Data, "list")
		assert.NotNil(t, list, "列表应包含 list 字段")

		var items []map[string]interface{}
		_ = json.Unmarshal(list, &items)
		assert.GreaterOrEqual(t, len(items), 1, "列表应至少包含刚创建的帖子")
		t.Logf("✅ 帖子列表返回 %d 条", len(items))
	})

	t.Run("按类型筛选", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", "/api/v1/posts?type=emotion", nil, nil)
		resp := assertOK(t, w)

		var pageData struct {
			List []map[string]interface{} `json:"list"`
		}
		_ = json.Unmarshal(extractField(t, resp.Data, "list"), &pageData.List)
		t.Log("✅ 类型筛选请求成功")
	})

	// ── Update ─────────────────────────────────────
	t.Run("更新帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "PUT", fmt.Sprintf("/api/v1/posts/%d", postID), map[string]interface{}{
			"title":   "修改后的感悟",
			"content": "修改后的内容：重新出发",
		}, authHeader(token))
		assertOK(t, w)

		// 验证更新生效
		detailW := doRequest(t, srv.Router, "GET", fmt.Sprintf("/api/v1/posts/%d", postID), nil, nil)
		detailResp := assertOK(t, detailW)

		title := extractField(t, detailResp.Data, "title")
		var titleStr string
		_ = jsonUnmarshal(t, title, &titleStr)
		assert.Equal(t, "修改后的感悟", titleStr)

		content := extractField(t, detailResp.Data, "content")
		var contentStr string
		_ = jsonUnmarshal(t, content, &contentStr)
		assert.Equal(t, "修改后的内容：重新出发", contentStr)
		t.Log("✅ 更新帖子成功")
	})

	// ── Delete ─────────────────────────────────────
	t.Run("删除帖子", func(t *testing.T) {
		// 先创建一个帖子来删除
		deleteID := createPost(t, srv, token, "case", "待删除帖子", "这条帖子将被删除", nil)

		w := doRequest(t, srv.Router, "DELETE", fmt.Sprintf("/api/v1/posts/%d", deleteID), nil, authHeader(token))
		assertOK(t, w)

		// 验证帖子已无法访问
		detailW := doRequest(t, srv.Router, "GET", fmt.Sprintf("/api/v1/posts/%d", deleteID), nil, nil)
		assertHTTPStatus(t, detailW, 400)
		t.Log("✅ 删除帖子成功")
	})
}

// ============================================================
// 错误场景测试
// ============================================================

func TestPost_AuthErrors(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	t.Run("未认证创建帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{
			"type":    "emotion",
			"title":   "未认证",
			"content": "应该失败",
		}, nil)
		assertHTTPStatus(t, w, 401)
		t.Log("✅ 未认证请求正确返回 401")
	})

	t.Run("未认证更新帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "PUT", "/api/v1/posts/1", map[string]interface{}{
			"title": "修改",
		}, nil)
		assertHTTPStatus(t, w, 401)
		t.Log("✅ 未认证更新正确返回 401")
	})

	t.Run("未认证删除帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "DELETE", "/api/v1/posts/1", nil, nil)
		assertHTTPStatus(t, w, 401)
		t.Log("✅ 未认证删除正确返回 401")
	})
}

func TestPost_NotFound(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	token := authToken(t, srv.Router, "notfound@test.com", "testpass123", "NotFound用户")

	t.Run("获取不存在的帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", "/api/v1/posts/999999", nil, nil)
		// 服务层返回 ErrPostNotFound，handler 映射为 400 + code 40301
		assertHTTPStatus(t, w, 400)
		assertAPIError(t, w, 40301)
		t.Log("✅ 不存在的帖子正确返回错误")
	})

	t.Run("更新不存在的帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "PUT", "/api/v1/posts/999999", map[string]interface{}{
			"title": "不应该成功",
		}, authHeader(token))
		assertHTTPStatus(t, w, 400)
		assertAPIError(t, w, 40301)
		t.Log("✅ 更新不存在帖子正确返回错误")
	})

	t.Run("删除不存在的帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "DELETE", "/api/v1/posts/999999", nil, authHeader(token))
		assertHTTPStatus(t, w, 400)
		assertAPIError(t, w, 40301)
		t.Log("✅ 删除不存在帖子正确返回错误")
	})

	t.Run("无效的帖子ID", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", "/api/v1/posts/abc", nil, nil)
		assertHTTPStatus(t, w, 400)
		t.Log("✅ 无效帖子 ID 正确返回 400")
	})
}

func TestPost_Forbidden(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	// 用户 A 创建帖子
	tokenA := authToken(t, srv.Router, "userA@forbidden.test.com", "pass1234", "用户A")
	postID := createPost(t, srv, tokenA, "insight", "A的帖子", "只有A能修改", nil)

	// 用户 B 尝试操作
	tokenB := authToken(t, srv.Router, "userB@forbidden.test.com", "pass1234", "用户B")

	t.Run("其他用户更新帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "PUT", fmt.Sprintf("/api/v1/posts/%d", postID), map[string]interface{}{
			"title": "B的修改",
		}, authHeader(tokenB))
		assertHTTPStatus(t, w, 400)
		assertAPIError(t, w, 40302)
		t.Log("✅ 非作者更新正确返回 forbidden")
	})

	t.Run("其他用户删除帖子", func(t *testing.T) {
		w := doRequest(t, srv.Router, "DELETE", fmt.Sprintf("/api/v1/posts/%d", postID), nil, authHeader(tokenB))
		assertHTTPStatus(t, w, 400)
		assertAPIError(t, w, 40302)
		t.Log("✅ 非作者删除正确返回 forbidden")
	})

	t.Run("作者仍可操作", func(t *testing.T) {
		w := doRequest(t, srv.Router, "PUT", fmt.Sprintf("/api/v1/posts/%d", postID), map[string]interface{}{
			"title": "A修改自己的帖子",
		}, authHeader(tokenA))
		assertOK(t, w)
		t.Log("✅ 作者更新自己的帖子成功")
	})
}

// ============================================================
// 参数校验测试
// ============================================================

func TestPost_Validation(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	token := authToken(t, srv.Router, "validate@test.com", "testpass123", "参数校验用户")

	t.Run("缺少必填 title", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{
			"type":    "emotion",
			"content": "有内容但没有标题",
		}, authHeader(token))
		assertHTTPStatus(t, w, 400)
		t.Log("✅ 缺少 title 正确返回 400")
	})

	t.Run("缺少必填 content", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{
			"type":  "emotion",
			"title": "有标题但没有内容",
		}, authHeader(token))
		assertHTTPStatus(t, w, 400)
		t.Log("✅ 缺少 content 正确返回 400")
	})

	t.Run("缺少必填 type", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{
			"title":   "有标题",
			"content": "有内容",
		}, authHeader(token))
		assertHTTPStatus(t, w, 400)
		t.Log("✅ 缺少 type 正确返回 400")
	})

	t.Run("无效的 type 值", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{
			"type":    "invalid_type",
			"title":   "测试",
			"content": "测试内容",
		}, authHeader(token))
		assertHTTPStatus(t, w, 400)
		t.Log("✅ 无效 type 正确返回 400")
	})

	t.Run("空 body", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{}, authHeader(token))
		assertHTTPStatus(t, w, 400)
		t.Log("✅ 空 body 正确返回 400")
	})
}

// ============================================================
// 列表分页测试
// ============================================================

func TestPost_ListPagination(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	token := authToken(t, srv.Router, "pagination@test.com", "testpass123", "分页测试用户")

	// 创建 5 个帖子
	for i := 1; i <= 5; i++ {
		createPost(t, srv, token, "emotion",
			fmt.Sprintf("分页帖子 %d", i),
			fmt.Sprintf("这是第 %d 条帖子的内容", i),
			nil,
		)
	}

	t.Run("第一页", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", "/api/v1/posts?page=1&size=3", nil, nil)
		resp := assertOK(t, w)

		// 验证分页结构
		list := extractField(t, resp.Data, "list")
		assert.NotNil(t, list)

		sizeField := extractField(t, resp.Data, "size")
		var sizeVal int
		_ = jsonUnmarshal(t, sizeField, &sizeVal)
		assert.Equal(t, 3, sizeVal, "page size 应为 3")

		hasMore := extractField(t, resp.Data, "has_more")
		var hasMoreVal bool
		_ = jsonUnmarshal(t, hasMore, &hasMoreVal)
		assert.True(t, hasMoreVal, "第一页应有 has_more=true")
		t.Log("✅ 分页第一页返回正确")
	})

	t.Run("第二页", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", "/api/v1/posts?page=2&size=3", nil, nil)
		resp := assertOK(t, w)

		hasMore := extractField(t, resp.Data, "has_more")
		var hasMoreVal bool
		_ = jsonUnmarshal(t, hasMore, &hasMoreVal)
		assert.False(t, hasMoreVal, "最后一页应 has_more=false")
		t.Log("✅ 分页第二页返回正确")
	})
}

// ============================================================
// 帖子类型覆盖测试
// ============================================================

func TestPost_Types(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	token := authToken(t, srv.Router, "types@test.com", "testpass123", "类型测试用户")

	types := map[string]string{
		"emotion":  "情绪分享",
		"insight":  "经验洞察",
		"case":     "真实案例",
		"action":   "行动记录",
	}

	for postType, desc := range types {
		t.Run(fmt.Sprintf("创建 %s 类型帖子", postType), func(t *testing.T) {
			w := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{
				"type":    postType,
				"title":   desc + "测试",
				"content": "这是一条 " + postType + " 类型的帖子",
			}, authHeader(token))
			resp := assertOK(t, w)

			// 验证类型返回正确
			typeField := extractField(t, resp.Data, "type")
			var typeStr string
			_ = jsonUnmarshal(t, typeField, &typeStr)
			assert.Equal(t, postType, typeStr)
			t.Logf("✅ %s 类型帖子创建成功", postType)
		})
	}
}

// ============================================================
// 幂等性测试 — 验证数据库层面创建
// ============================================================

func TestPost_CreateReturnsValidData(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	token := authToken(t, srv.Router, "datainteg@test.com", "testpass123", "数据完整性用户")

	w := doRequest(t, srv.Router, "POST", "/api/v1/posts", map[string]interface{}{
		"type":    "emotion",
		"title":   "数据完整性测试",
		"content": "验证返回数据是否完整",
		"images":  []string{"http://example.com/a.jpg", "http://example.com/b.jpg"},
		"tags":    []string{"标签A", "标签B"},
	}, authHeader(token))

	resp := assertOK(t, w)

	// 验证返回的帖子数据包含所有必要字段
	dataBytes, _ := json.Marshal(resp.Data)
	var post map[string]interface{}
	err := json.Unmarshal(dataBytes, &post)
	require.NoError(t, err)

	// 基本字段
	assert.Greater(t, post["id"], float64(0), "应有 id")
	assert.Equal(t, "emotion", post["type"])
	assert.Equal(t, "数据完整性测试", post["title"])
	assert.Equal(t, "验证返回数据是否完整", post["content"])
	assert.Equal(t, float64(0), post["like_count"], "新帖子 like_count 应为 0")
	assert.Equal(t, float64(0), post["comment_count"], "新帖子 comment_count 应为 0")

	// 图片和标签（存储为 JSON 字符串）
	if imagesRaw, ok := post["images"].(string); ok {
		var images []interface{}
		err = json.Unmarshal([]byte(imagesRaw), &images)
		assert.NoError(t, err, "images 应为有效的 JSON")
		assert.Equal(t, 2, len(images), "应有 2 张图片")
	} else {
		t.Logf("images 类型: %T, 值: %v", post["images"], post["images"])
		assert.True(t, post["images"] != nil, "images 不应为 nil")
	}

	if tagsRaw, ok := post["tags"].(string); ok {
		var tags []interface{}
		err = json.Unmarshal([]byte(tagsRaw), &tags)
		assert.NoError(t, err, "tags 应为有效的 JSON")
		assert.Equal(t, 2, len(tags), "应有 2 个标签")
	} else {
		t.Logf("tags 类型: %T, 值: %v", post["tags"], post["tags"])
		assert.True(t, post["tags"] != nil, "tags 不应为 nil")
	}

	t.Log("✅ 创建返回数据完整性验证通过")
}
