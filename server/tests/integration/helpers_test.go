package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ============================================================
// HTTP 请求工具
// ============================================================

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// APIResponse 统一响应结构
type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// AuthData 登录成功后返回的数据
type AuthData struct {
	Token string `json:"token"`
	User  struct {
		ID       int64  `json:"id"`
		Nickname string `json:"nickname"`
	} `json:"user"`
}

// doRequest 发起 HTTP 请求
func doRequest(t *testing.T, router http.Handler, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("序列化请求体失败: %v", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// doRequestRaw 发起原生 HTTP 请求（不自动设置 Content-Type）
func doRequestRaw(t *testing.T, router http.Handler, method, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ============================================================
// 认证工具
// ============================================================

// authToken 通过真实 API 注册并登录，返回 JWT Token
func authToken(t *testing.T, router http.Handler, email, password, nickname string) string {
	t.Helper()

	// 注册
	regData := &RegisterRequest{Email: email, Password: password, Nickname: nickname}
	w := doRequest(t, router, "POST", "/api/v1/auth/register/email", regData, nil)

	if w.Code != http.StatusOK {
		t.Fatalf("注册用户 %s 失败: HTTP %d, body: %s", email, w.Code, w.Body.String())
	}

	// 从注册响应提取 token
	resp := parseResponse(t, w)
	var authData AuthData
	if err := json.Unmarshal(resp.Data, &authData); err != nil {
		t.Fatalf("解析注册响应失败: %v, body: %s", err, w.Body.String())
	}

	if authData.Token != "" {
		return authData.Token
	}

	// 注册没返回 token（如已有该用户），则登录
	loginW := authTokenLogin(t, router, email, password)
	return loginW
}

// authTokenLogin 直接登录，返回 token 字符串
func authTokenLogin(t *testing.T, router http.Handler, email, password string) string {
	t.Helper()

	loginData := &LoginRequest{Email: email, Password: password}
	w := doRequest(t, router, "POST", "/api/v1/auth/login/email", loginData, nil)

	if w.Code != http.StatusOK {
		t.Fatalf("登录用户 %s 失败: HTTP %d, body: %s", email, w.Code, w.Body.String())
	}

	resp := parseResponse(t, w)
	var authData AuthData
	if err := json.Unmarshal(resp.Data, &authData); err != nil {
		t.Fatalf("解析登录响应失败: %v, body: %s", err, w.Body.String())
	}

	return authData.Token
}

// authHeader 返回带 JWT 认证的请求头
func authHeader(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// ============================================================
// 响应断言工具
// ============================================================

// assertOK 断言 HTTP 200 且 code=0
func assertOK(t *testing.T, w *httptest.ResponseRecorder) APIResponse {
	t.Helper()
	if w.Code != http.StatusOK {
		t.Fatalf("期望 HTTP 200, 实际 %d, body: %s", w.Code, w.Body.String())
	}
	return parseResponse(t, w)
}

// assertHTTPStatus 断言 HTTP 状态码
func assertHTTPStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) APIResponse {
	t.Helper()
	if w.Code != expected {
		t.Fatalf("期望 HTTP %d, 实际 %d, body: %s", expected, w.Code, w.Body.String())
	}
	return parseResponse(t, w)
}

// assertAPIError 断言 API 错误码
func assertAPIError(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) APIResponse {
	t.Helper()
	resp := parseResponse(t, w)
	if resp.Code != expectedCode {
		t.Fatalf("期望 API code=%d, 实际 code=%d, message: %s", expectedCode, resp.Code, resp.Message)
	}
	return resp
}

// parseResponse 解析 JSON 响应
func parseResponse(t *testing.T, w *httptest.ResponseRecorder) APIResponse {
	t.Helper()
	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v, body: %s", err, w.Body.String())
	}
	return resp
}

// extractField 从 JSON Data 中提取指定字段的值
func extractField(t *testing.T, data json.RawMessage, field string) json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("解析 data 失败: %v", err)
	}
	v, ok := m[field]
	if !ok {
		t.Fatalf("data 中未找到字段 %q，可用字段: %v", field, fieldKeys(m))
	}
	return v
}

// extractID 从 JSON Data 中提取 id 字段
// 先尝试顶层 id（帖子/评论/用户等响应），再尝试 data.user.id（认证响应）
func extractID(t *testing.T, data json.RawMessage) int64 {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("解析 data 失败: %v", err)
	}
	// 顶层 id
	if raw, ok := m["id"]; ok {
		var id int64
		if err := json.Unmarshal(raw, &id); err == nil && id > 0 {
			return id
		}
	}
	// 认证响应: data.user.id
	if userRaw, ok := m["user"]; ok {
		var user struct {
			ID int64 `json:"id"`
		}
		if err := json.Unmarshal(userRaw, &user); err == nil {
			return user.ID
		}
	}
	// 无法提取时返回 0（让调用方决定是否 fatalf）
	return 0
}

// fieldKeys 返回 map 的所有 key
func fieldKeys(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// jsonUnmarshal 将 json.Unmarshal 包装为断言友好的辅助函数
func jsonUnmarshal(t *testing.T, data []byte, v interface{}) error {
	t.Helper()
	return json.Unmarshal(data, v)
}

// ============================================================
// 冒烟测试（验证集成测试基础设施可用）
// ============================================================

func TestIntegration_Infrastructure_SmokeTest(t *testing.T) {
	srv := NewTestServer()

	// 1. 健康检查 — 不需要数据库表，纯通路验证
	w := doRequest(t, srv.Router, "GET", "/health", nil, nil)
	assertOK(t, w)
	t.Log("✅ 健康检查端点可用")

	// 2. 注册 → 登录 → 获取 Profile（真实端到端）
	token := authToken(t, srv.Router, "smoke@test.com", "123456", "冒烟用户")
	if token == "" {
		t.Fatal("期望返回非空 token")
	}
	t.Log("✅ 注册+登录成功，已获取 token")

	// 3. 带 token 请求 /auth/profile
	profileW := doRequest(t, srv.Router, "GET", "/api/v1/auth/profile", nil, authHeader(token))
	assertOK(t, profileW)
	t.Log("✅ 获取 Profile 成功")

	// 4. 创建帖子
	postData := map[string]interface{}{
		"type":    "emotion",
		"title":   "集成测试帖子",
		"content": "这是一条集成测试自动创建的帖子",
	}
	createPostW := doRequest(t, srv.Router, "POST", "/api/v1/posts", postData, authHeader(token))
	resp := assertOK(t, createPostW)
	postID := extractID(t, resp.Data)
	if postID == 0 {
		t.Fatal("期望返回非零帖子 ID")
	}
	t.Logf("✅ 创建帖子成功，ID: %d", postID)

	// 5. 查询帖子详情（公开接口）
	detailW := doRequest(t, srv.Router, "GET", fmt.Sprintf("/api/v1/posts/%d", postID), nil, nil)
	detailResp := assertOK(t, detailW)
	title := extractField(t, detailResp.Data, "title")
	var titleStr string
	json.Unmarshal(title, &titleStr)
	if titleStr != "集成测试帖子" {
		t.Fatalf("期望标题 '集成测试帖子', 实际 '%s'", titleStr)
	}
	t.Log("✅ 查询帖子详情成功")

	// 6. 清理后验证
	CleanupTestData()
	t.Log("✅ CleanupTestData 执行成功，基础设施验收通过")
}
