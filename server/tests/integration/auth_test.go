package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================
// Auth 模块集成测试
// 验证：注册→登录→获取 Profile 端到端流程
// ============================================================

func TestAuth_RegisterAndLogin(t *testing.T) {
	srv := NewTestServer()

	t.Run("注册成功", func(t *testing.T) {
		CleanupTestData()
		email := "auth-integration@example.com"
		password := "testpass123"
		nickname := "集成测试用户"

		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/register/email", &RegisterRequest{
			Email:    email,
			Password: password,
			Nickname: nickname,
		}, nil)

		resp := assertOK(t, w)
		token := extractField(t, resp.Data, "token")
		assert.NotEmpty(t, token, "注册应返回 JWT token")

		userData := extractField(t, resp.Data, "user")
		assert.NotEmpty(t, userData, "注册应返回用户数据")
	})

	t.Run("登录成功", func(t *testing.T) {
		CleanupTestData()
		email := "login-test@example.com"
		password := "securepass"
		nickname := "登录测试用户"

		// 先注册
		token := authToken(t, srv.Router, email, password, nickname)
		assert.NotEmpty(t, token, "注册后应获得 token")

		// 再登录
		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/login/email", &LoginRequest{
			Email:    email,
			Password: password,
		}, nil)
		loginResp := assertOK(t, w)
		loginToken := extractField(t, loginResp.Data, "token")
		assert.NotEmpty(t, loginToken, "登录应返回 JWT token")
	})

	t.Run("获取 Profile", func(t *testing.T) {
		CleanupTestData()
		email := "profile-test@example.com"
		nickname := "Profile测试"
		token := authToken(t, srv.Router, email, "testpass123", nickname)

		w := doRequest(t, srv.Router, "GET", "/api/v1/auth/profile", nil, authHeader(token))
		resp := assertOK(t, w)

		// 验证返回了正确的用户数据
		profileID := extractID(t, resp.Data)
		assert.Greater(t, profileID, int64(0), "Profile 应包含有效的用户 ID")

		nick := extractField(t, resp.Data, "nickname")
		var nickStr string
		_ = jsonUnmarshal(t, nick, &nickStr)
		assert.Equal(t, nickname, nickStr)
	})
}

func TestAuth_RegisterErrors(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	// 先注册一个用户
	email := "existing@example.com"
	token := authToken(t, srv.Router, email, "testpass123", "已有用户")
	_ = token

	t.Run("重复邮箱注册", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/register/email", &RegisterRequest{
			Email:    email,
			Password: "testpass123",
			Nickname: "重复用户",
		}, nil)

		assertHTTPStatus(t, w, 400)
		assertAPIError(t, w, 40206) // ErrEmailExists
	})

	t.Run("密码太短", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/register/email", &RegisterRequest{
			Email:    "short@example.com",
			Password: "123",
			Nickname: "密码太短",
		}, nil)

		assertHTTPStatus(t, w, 400)
	})

	t.Run("邮箱格式错误", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/register/email", &RegisterRequest{
			Email:    "not-an-email",
			Password: "testpass123",
			Nickname: "邮箱错误",
		}, nil)

		assertHTTPStatus(t, w, 400)
	})

	t.Run("昵称太短", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/register/email", &RegisterRequest{
			Email:    "nickname@example.com",
			Password: "testpass123",
			Nickname: "A",
		}, nil)

		assertHTTPStatus(t, w, 400)
	})
}

func TestAuth_LoginErrors(t *testing.T) {
	srv := NewTestServer()
	CleanupTestData()

	email := "logintest@example.com"
	token := authToken(t, srv.Router, email, "correctpass", "登录测试")
	_ = token

	t.Run("错误密码", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/login/email", &LoginRequest{
			Email:    email,
			Password: "wrongpass",
		}, nil)

		assertHTTPStatus(t, w, 400)
		assertAPIError(t, w, 40204) // ErrPasswordWrong
	})

	t.Run("不存在的邮箱", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/login/email", &LoginRequest{
			Email:    "noexist@example.com",
			Password: "anypass",
		}, nil)

		assertHTTPStatus(t, w, 400)
	})

	t.Run("未认证访问 Profile", func(t *testing.T) {
		w := doRequest(t, srv.Router, "GET", "/api/v1/auth/profile", nil, nil)
		assertHTTPStatus(t, w, 401)
	})
}

func TestAuth_EndToEndFlow(t *testing.T) {
	// 端到端验证：注册 → 登录 → 获取 Profile
	srv := NewTestServer()
	CleanupTestData()

	email := "e2e@example.com"
	password := "e2epassword123"
	nickname := "E2E测试用户"
	var savedToken string

	t.Run("注册", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/register/email", &RegisterRequest{
			Email:    email,
			Password: password,
			Nickname: nickname,
		}, nil)
		resp := assertOK(t, w)

		uid := extractID(t, resp.Data)
		assert.Greater(t, uid, int64(0))
		savedToken = extractStringField(t, resp.Data, "token")
	})

	t.Run("登录", func(t *testing.T) {
		w := doRequest(t, srv.Router, "POST", "/api/v1/auth/login/email", &LoginRequest{
			Email:    email,
			Password: password,
		}, nil)
		resp := assertOK(t, w)
		assert.NotEmpty(t, extractStringField(t, resp.Data, "token"))
	})

	t.Run("获取 Profile", func(t *testing.T) {
		token := savedToken
		if token == "" {
			t.Fatal("上一步未获取到 token")
		}
		w := doRequest(t, srv.Router, "GET", "/api/v1/auth/profile", nil, authHeader(token))
		resp := assertOK(t, w)

		nickField := extractField(t, resp.Data, "nickname")
		var got string
		_ = jsonUnmarshal(t, nickField, &got)
		assert.Equal(t, nickname, got)
	})
}
