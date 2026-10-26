package setting

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/shiyemo/server/internal/common/config"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/email"
)

// ============================================================
// Mock SettingServiceInterface
// ============================================================

type mockSettingService struct {
	mock.Mock
}

func (m *mockSettingService) GetSMTP() (*SMTPSettings, error) {
	args := m.Called()
	if v := args.Get(0); v != nil {
		return v.(*SMTPSettings), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockSettingService) GetSMTPWithPassword() (*SMTPSettings, error) {
	args := m.Called()
	if v := args.Get(0); v != nil {
		return v.(*SMTPSettings), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockSettingService) UpdateSMTP(req *UpdateSMTPRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *mockSettingService) SMTPConfigProvider() email.ConfigProvider {
	// 不实现 —— handler 三个端点没用这个方法
	// 仍需返回 nil / panic 都要实现接口签名
	return func() *config.SMTPConfig { return nil }
}

// TestSend 方法签名带 (to string) error，handler 调的就是它
func (m *mockSettingService) TestSend(to string) error {
	args := m.Called(to)
	return args.Error(0)
}

func (m *mockSettingService) IsEmailVerifyRequired() bool {
	args := m.Called()
	return args.Bool(0)
}

// ============================================================
// Helpers
// ============================================================

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/admin/settings/smtp", h.GetSMTP)
	r.PUT("/api/v1/admin/settings/smtp", h.UpdateSMTP)
	r.POST("/api/v1/admin/settings/smtp/test", h.TestSMTP)
	return r
}

func doRequest(t *testing.T, r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var reqBody *bytes.Buffer
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v\nbody=%s", err, w.Body.String())
	}
	return resp
}

// ============================================================
// GetSMTP
// ============================================================

func TestHandler_GetSMTP_Success(t *testing.T) {
	svc := new(mockSettingService)
	// service 层负责抹除 password（合约），handler 不应再加一层
	svc.On("GetSMTP").Return(&SMTPSettings{
		Host:                "smtp.example.com",
		Port:                587,
		Username:            "admin",
		Password:            "",
		PasswordSet:         true, // 标识曾设置过
		From:                "noreply@example.com",
		Enabled:             true,
		EmailVerifyRequired: true,
	}, nil)

	h := NewHandlerWithService(svc)
	r := setupRouter(h)
	w := doRequest(t, r, "GET", "/api/v1/admin/settings/smtp", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w)
	assert.EqualValues(t, 0, resp["code"])
	assert.Equal(t, "ok", resp["message"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "smtp.example.com", data["host"])
	assert.EqualValues(t, 587, data["port"])
	assert.Equal(t, true, data["enabled"])
	assert.Equal(t, true, data["password_set"])
	assert.Equal(t, "", data["password"])
	svc.AssertExpectations(t)
}

func TestHandler_GetSMTP_DBError(t *testing.T) {
	svc := new(mockSettingService)
	svc.On("GetSMTP").Return(nil, errors.New("db down"))

	h := NewHandlerWithService(svc)
	r := setupRouter(h)
	w := doRequest(t, r, "GET", "/api/v1/admin/settings/smtp", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	resp := decodeResponse(t, w)
	assert.EqualValues(t, bizerr.ErrDB.Code, resp["code"])
	assert.Equal(t, bizerr.ErrDB.Message, resp["message"])
	svc.AssertExpectations(t)
}

// ============================================================
// UpdateSMTP
// ============================================================

func TestHandler_UpdateSMTP_Success(t *testing.T) {
	svc := new(mockSettingService)
	req := &UpdateSMTPRequest{
		Enabled:             true,
		Host:                "smtp.example.com",
		Port:                587,
		Username:            "admin",
		From:                "noreply@example.com",
		Password:            "newpwd",
		EmailVerifyRequired: true,
	}
	svc.On("UpdateSMTP", req).Return(nil)

	h := NewHandlerWithService(svc)
	r := setupRouter(h)
	w := doRequest(t, r, "PUT", "/api/v1/admin/settings/smtp", req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, true, data["updated"])
	svc.AssertExpectations(t)
}

func TestHandler_UpdateSMTP_BadJSON(t *testing.T) {
	svc := new(mockSettingService)
	// 不期望 service 被调用

	h := NewHandlerWithService(svc)
	r := setupRouter(h)

	// 直接发非法 JSON
	req, _ := http.NewRequest("PUT", "/api/v1/admin/settings/smtp",
		strings.NewReader("{not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := decodeResponse(t, w)
	// 业务码 40001（参数错误），不是 bizerr.ErrDB
	assert.EqualValues(t, 40001, resp["code"])
	// service 不应被调用
	svc.AssertNotCalled(t, "UpdateSMTP")
}

func TestHandler_UpdateSMTP_InvalidEmail(t *testing.T) {
	svc := new(mockSettingService)

	h := NewHandlerWithService(svc)
	r := setupRouter(h)

	w := doRequest(t, r, "PUT", "/api/v1/admin/settings/smtp", map[string]interface{}{
		"enabled": true,
		"from":    "not-an-email",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := decodeResponse(t, w)
	assert.EqualValues(t, 40001, resp["code"])
	svc.AssertNotCalled(t, "UpdateSMTP")
}

func TestHandler_UpdateSMTP_ServiceError(t *testing.T) {
	svc := new(mockSettingService)
	req := &UpdateSMTPRequest{Enabled: true, Host: "x", Port: 587, Username: "u", From: "f@example.com"}
	svc.On("UpdateSMTP", req).Return(errors.New("write fail"))

	h := NewHandlerWithService(svc)
	r := setupRouter(h)
	w := doRequest(t, r, "PUT", "/api/v1/admin/settings/smtp", req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	resp := decodeResponse(t, w)
	assert.EqualValues(t, bizerr.ErrDB.Code, resp["code"])
	svc.AssertExpectations(t)
}

// ============================================================
// TestSMTP
// ============================================================

func TestHandler_TestSMTP_Success(t *testing.T) {
	svc := new(mockSettingService)
	svc.On("TestSend", "admin@example.com").Return(nil)

	h := NewHandlerWithService(svc)
	r := setupRouter(h)
	w := doRequest(t, r, "POST", "/api/v1/admin/settings/smtp/test",
		map[string]string{"to": "admin@example.com"})

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, true, data["sent"])
	svc.AssertExpectations(t)
}

func TestHandler_TestSMTP_BadJSON(t *testing.T) {
	svc := new(mockSettingService)

	h := NewHandlerWithService(svc)
	r := setupRouter(h)
	req, _ := http.NewRequest("POST", "/api/v1/admin/settings/smtp/test",
		strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "TestSend")
}

func TestHandler_TestSMTP_SendError_BubblesUp(t *testing.T) {
	svc := new(mockSettingService)
	// service 抛 SMTP 错误，handler 应把原始错误透出，否则管理员无从排查
	svc.On("TestSend", "admin@example.com").Return(errors.New("auth failed: 535"))

	h := NewHandlerWithService(svc)
	r := setupRouter(h)
	w := doRequest(t, r, "POST", "/api/v1/admin/settings/smtp/test",
		map[string]string{"to": "admin@example.com"})

	assert.Equal(t, http.StatusBadGateway, w.Code)
	resp := decodeResponse(t, w)
	assert.EqualValues(t, 50002, resp["code"])
	msg := resp["message"].(string)
	// 原始错误信息应包含在 message 里，方便排查
	assert.Contains(t, msg, "auth failed")
	assert.Contains(t, msg, "535")
	svc.AssertExpectations(t)
}