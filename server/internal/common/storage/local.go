package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage 本地文件存储实现
type LocalStorage struct {
	uploadDir string
	baseURL   string
}

// NewLocalStorage 创建本地存储
func NewLocalStorage(cfg LocalConfig) *LocalStorage {
	return &LocalStorage{
		uploadDir: cfg.UploadDir,
		baseURL:   cfg.BaseURL,
	}
}

// Upload 上传文件到本地
func (s *LocalStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
	// 确保 key 不以 / 开头
	key = strings.TrimPrefix(key, "/")
	filePath := filepath.Join(s.uploadDir, key)

	// 创建目录
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 创建文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return s.GetURL(key), nil
}

// Delete 删除本地文件
func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	// 确保 key 不以 / 开头
	key = strings.TrimPrefix(key, "/")
	filePath := filepath.Join(s.uploadDir, key)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // 文件不存在，视为成功
	}

	return os.Remove(filePath)
}

// GetURL 获取文件访问 URL
func (s *LocalStorage) GetURL(key string) string {
	// 确保 key 不以 / 开头（URL 拼接）
	key = strings.TrimPrefix(key, "/")
	return fmt.Sprintf("%s/%s", s.baseURL, key)
}
