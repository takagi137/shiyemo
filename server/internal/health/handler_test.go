package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// 设置 Gin 测试模式
	gin.SetMode(gin.TestMode)

	// 创建处理器
	handler := NewHandler()

	// 创建测试路由
	r := gin.New()
	r.GET("/health", handler.HealthCheck)

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
	assert.NotNil(t, response["timestamp"])
}

func TestDetailedHealthCheck(t *testing.T) {
	// 设置 Gin 测试模式
	gin.SetMode(gin.TestMode)

	// 创建处理器
	handler := NewHandler()

	// 创建测试路由
	r := gin.New()
	r.GET("/health/detailed", handler.DetailedHealthCheck)

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/health/detailed", nil)
	w := httptest.NewRecorder()

	// 执行请求
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 验证响应结构
	assert.Contains(t, response, "status")
	assert.Contains(t, response, "timestamp")
	assert.Contains(t, response, "checks")

	// 验证 checks 结构
	checks, ok := response["checks"].(map[string]interface{})
	assert.True(t, ok)

	// 验证系统信息
	system, ok := checks["system"].(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, system, "goroutines")
	assert.Contains(t, system, "heap_alloc_mb")
	assert.Contains(t, system, "sys_mb")
	assert.Contains(t, system, "gc_cycles")

	// 验证运行时间
	uptime, ok := checks["uptime"].(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, uptime, "seconds")
	assert.Contains(t, uptime, "human")
}