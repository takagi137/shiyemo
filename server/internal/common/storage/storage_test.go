package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage_Local(t *testing.T) {
	cfg := Config{
		Type: "local",
		Local: LocalConfig{
			UploadDir: "/tmp/uploads",
			BaseURL:   "/uploads",
		},
	}

	store, err := NewStorage(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, store)

	// 验证是 LocalStorage 类型
	_, ok := store.(*LocalStorage)
	assert.True(t, ok)
}

func TestNewStorage_LocalDefault(t *testing.T) {
	cfg := Config{
		Type: "", // 空值应该默认为 local
		Local: LocalConfig{
			UploadDir: "/tmp/uploads",
			BaseURL:   "/uploads",
		},
	}

	store, err := NewStorage(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, store)
}

func TestNewStorage_UnsupportedType(t *testing.T) {
	cfg := Config{
		Type: "unsupported",
	}

	store, err := NewStorage(cfg)
	assert.Error(t, err)
	assert.Nil(t, store)
	assert.Contains(t, err.Error(), "unsupported storage type")
}

func TestLocalStorage_GetURL(t *testing.T) {
	store := NewLocalStorage(LocalConfig{
		UploadDir: "/tmp/uploads",
		BaseURL:   "/uploads",
	})

	tests := []struct {
		key      string
		expected string
	}{
		{"images/test.jpg", "/uploads/images/test.jpg"},
		{"/images/test.jpg", "/uploads/images/test.jpg"},
		{"test.jpg", "/uploads/test.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			url := store.GetURL(tt.key)
			assert.Equal(t, tt.expected, url)
		})
	}
}

func TestLocalStorage_Delete_NonExistent(t *testing.T) {
	store := NewLocalStorage(LocalConfig{
		UploadDir: t.TempDir(),
		BaseURL:   "/uploads",
	})

	// 删除不存在的文件应该返回 nil
	err := store.Delete(context.Background(), "nonexistent.jpg")
	assert.NoError(t, err)
}

func TestMinIOStorage_GetURL(t *testing.T) {
	store := &MinIOStorage{
		bucket:  "test-bucket",
		baseURL: "http://minio.example.com",
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"images/test.jpg", "http://minio.example.com/test-bucket/images/test.jpg"},
		{"/images/test.jpg", "http://minio.example.com/test-bucket/images/test.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			url := store.GetURL(tt.key)
			assert.Equal(t, tt.expected, url)
		})
	}
}
