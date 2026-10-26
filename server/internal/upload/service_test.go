package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Mock Storage
// ============================================================

// MockStorage 模拟存储实现
type MockStorage struct {
	uploadFunc func(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error)
	deleteFunc func(ctx context.Context, key string) error
	uploaded   []mockUpload
	deleted    []string
}

type mockUpload struct {
	key         string
	size        int64
	contentType string
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		uploadFunc: func(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
			return fmt.Sprintf("http://mock.storage/%s", key), nil
		},
		deleteFunc: func(ctx context.Context, key string) error {
			return nil
		},
	}
}

func (m *MockStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
	m.uploaded = append(m.uploaded, mockUpload{key: key, size: size, contentType: contentType})
	return m.uploadFunc(ctx, key, reader, size, contentType)
}

func (m *MockStorage) Delete(ctx context.Context, key string) error {
	m.deleted = append(m.deleted, key)
	return m.deleteFunc(ctx, key)
}

func (m *MockStorage) GetURL(key string) string {
	return fmt.Sprintf("http://mock.storage/%s", key)
}

// ============================================================
// 测试辅助函数
// ============================================================

func setupTestService(t *testing.T) (*UploadService, *MockStorage) {
	t.Helper()
	mock := NewMockStorage()
	svc := NewUploadServiceWithStorage(mock)
	return svc, mock
}

// ============================================================
// Upload 测试
// ============================================================

func TestUpload_Success(t *testing.T) {
	svc, mock := setupTestService(t)

	tests := []struct {
		name        string
		filename    string
		content     string
		category    string
		contentType string
	}{
		{"jpg文件", "test.jpg", "fake jpg content", "images", "image/jpeg"},
		{"jpeg文件", "test.jpeg", "fake jpeg content", "images", "image/jpeg"},
		{"png文件", "test.png", "fake png content", "avatars", "image/png"},
		{"gif文件", "test.gif", "fake gif content", "images", "image/gif"},
		{"webp文件", "test.webp", "fake webp content", "images", "image/webp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置 mock
			mock.uploaded = nil

			reader := bytes.NewReader([]byte(tt.content))
			result, err := svc.Upload(123, reader, tt.filename, tt.category)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.URL)
			assert.NotEmpty(t, result.Filename)
			assert.Equal(t, int64(len(tt.content)), result.Size)

			// 验证 storage 被正确调用
			require.Len(t, mock.uploaded, 1)
			assert.True(t, strings.HasPrefix(mock.uploaded[0].key, tt.category+"/"))
			assert.Equal(t, tt.contentType, mock.uploaded[0].contentType)
			assert.Equal(t, int64(len(tt.content)), mock.uploaded[0].size)
		})
	}
}

func TestUpload_UnsupportedFormat(t *testing.T) {
	svc, _ := setupTestService(t)

	tests := []struct {
		name     string
		filename string
	}{
		{"txt文件", "test.txt"},
		{"pdf文件", "test.pdf"},
		{"mp4文件", "test.mp4"},
		{"无扩展名", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte("content"))
			result, err := svc.Upload(123, reader, tt.filename, "images")

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "不支持的文件格式")
		})
	}
}

func TestUpload_FileSizeLimit(t *testing.T) {
	svc, mock := setupTestService(t)

	// 创建超过限制的内容（10MB + 1）
	largeContent := make([]byte, 10*1024*1024+1)
	reader := bytes.NewReader(largeContent)

	result, err := svc.Upload(123, reader, "large.jpg", "images")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "文件大小超过限制")

	// 验证 storage 未被调用
	assert.Empty(t, mock.uploaded)
}

func TestUpload_StorageError(t *testing.T) {
	mock := NewMockStorage()
	mock.uploadFunc = func(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
		return "", fmt.Errorf("storage error")
	}
	svc := NewUploadServiceWithStorage(mock)

	reader := bytes.NewReader([]byte("content"))
	result, err := svc.Upload(123, reader, "test.jpg", "images")

	assert.Error(t, err)
	assert.Nil(t, result)
	// 错误会被包装为内部错误
	assert.NotNil(t, err)
}

func TestUpload_FilenameContainsUserID(t *testing.T) {
	svc, mock := setupTestService(t)

	reader := bytes.NewReader([]byte("content"))
	result, err := svc.Upload(456, reader, "test.jpg", "images")

	require.NoError(t, err)

	// 验证文件名包含用户 ID
	assert.True(t, strings.HasPrefix(result.Filename, "456_"))

	// 验证 key 包含用户 ID
	assert.True(t, strings.Contains(mock.uploaded[0].key, "456_"))
}

func TestUpload_ContentTypeDetection(t *testing.T) {
	svc, mock := setupTestService(t)

	tests := []struct {
		ext         string
		contentType string
	}{
		{".jpg", "image/jpeg"},
		{".jpeg", "image/jpeg"},
		{".png", "image/png"},
		{".gif", "image/gif"},
		{".webp", "image/webp"},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			mock.uploaded = nil
			reader := bytes.NewReader([]byte("content"))
			_, err := svc.Upload(123, reader, "test"+tt.ext, "images")

			require.NoError(t, err)
			require.Len(t, mock.uploaded, 1)
			assert.Equal(t, tt.contentType, mock.uploaded[0].contentType)
		})
	}
}

// ============================================================
// Delete 测试
// ============================================================

func TestDelete_Success(t *testing.T) {
	svc, mock := setupTestService(t)

	// 设置 mock 记录删除调用
	mock.deleteFunc = func(ctx context.Context, key string) error {
		return nil
	}

	err := svc.Delete("images/123_456.jpg")
	assert.NoError(t, err)

	// 验证 storage.Delete 被调用
	assert.Len(t, mock.deleted, 1)
	assert.Equal(t, "images/123_456.jpg", mock.deleted[0])
}

func TestDelete_StorageError(t *testing.T) {
	mock := NewMockStorage()
	mock.deleteFunc = func(ctx context.Context, key string) error {
		return fmt.Errorf("delete failed")
	}
	svc := NewUploadServiceWithStorage(mock)

	err := svc.Delete("images/test.jpg")
	assert.Error(t, err)
}
