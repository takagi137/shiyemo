package log

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogService(t *testing.T) {
	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "log-test")
	defer os.RemoveAll(tempDir)

	// 创建日志服务
	service, err := NewLogService(tempDir, DEBUG)
	assert.NoError(t, err)
	assert.NotNil(t, service)

	defer service.Close()

	// 测试日志记录
	service.Info("测试消息", map[string]interface{}{
		"key": "value",
	})

	// 验证日志文件创建
	files, err := os.ReadDir(tempDir)
	assert.NoError(t, err)
	assert.True(t, len(files) > 0)

	// 验证日志文件内容
	logFile := filepath.Join(tempDir, files[0].Name())
	content, err := os.ReadFile(logFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "测试消息")
	assert.Contains(t, string(content), "key=value")
}

func TestRequestLogger(t *testing.T) {
	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "log-test-request")
	defer os.RemoveAll(tempDir)

	// 创建日志服务
	service, err := NewLogService(tempDir, DEBUG)
	assert.NoError(t, err)
	assert.NotNil(t, service)

	defer service.Close()

	// 创建请求日志记录器
	reqLogger := service.WithRequest("test-request-id", 12345)
	assert.NotNil(t, reqLogger)

	// 测试各种日志级别
	reqLogger.Debug("调试消息")
	reqLogger.Info("信息消息")
	reqLogger.Warn("警告消息")
	reqLogger.Error("错误消息")

	// 验证日志文件
	files, err := os.ReadDir(tempDir)
	assert.NoError(t, err)
	assert.True(t, len(files) > 0)

	logFile := filepath.Join(tempDir, files[0].Name())
	content, err := os.ReadFile(logFile)
	assert.NoError(t, err)

	// 验证请求ID和用户ID
	assert.Contains(t, string(content), "request_id=test-request-id")
	assert.Contains(t, string(content), "user_id=12345")

	// 验证不同日志级别
	assert.Contains(t, string(content), "[DEBUG] 调试消息")
	assert.Contains(t, string(content), "[INFO] 信息消息")
	assert.Contains(t, string(content), "[WARN] 警告消息")
	assert.Contains(t, string(content), "[ERROR] 错误消息")
}

func TestLogLevelFiltering(t *testing.T) {
	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "log-test-level")
	defer os.RemoveAll(tempDir)

	// 创建日志服务，级别为 WARN
	service, err := NewLogService(tempDir, WARN)
	assert.NoError(t, err)
	assert.NotNil(t, service)

	defer service.Close()

	// 记录不同级别的日志
	service.Debug("调试消息")
	service.Info("信息消息")
	service.Warn("警告消息")
	service.Error("错误消息")

	// 验证只有 WARN 和 ERROR 被记录
	files, err := os.ReadDir(tempDir)
	assert.NoError(t, err)
	assert.True(t, len(files) > 0)

	logFile := filepath.Join(tempDir, files[0].Name())
	content, err := os.ReadFile(logFile)
	assert.NoError(t, err)

	assert.NotContains(t, string(content), "调试消息")
	assert.NotContains(t, string(content), "信息消息")
	assert.Contains(t, string(content), "警告消息")
	assert.Contains(t, string(content), "错误消息")
}