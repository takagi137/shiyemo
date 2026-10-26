package tips

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestRouter 创建测试路由（mock 模式，不连数据库）
func setupTestRouter(handler *TipsHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/daily/quote", handler.RandomQuote)
	return r
}

func TestRandomQuote_BDD(t *testing.T) {
	// Given: tips 表中有文案
	repo := newMockTipsRepo()
	repo.tips = append(repo.tips, SurvivalTip{
		ID: 1, Content: "慢慢来，你已经在路上了。", Category: "encourage", IsActive: true,
	})
	handler := &TipsHandler{service: NewTipsServiceWithDeps(repo)}
	r := setupTestRouter(handler)

	// When: 用户请求 GET /daily/quote
	req, _ := http.NewRequest("GET", "/daily/quote", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then: 返回 200，包含文案
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "慢慢来，你已经在路上了。", data["text"])
}
