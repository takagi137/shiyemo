package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/storage"
)

// UploadService 文件上传服务
type UploadService struct {
	storage storage.Storage
	maxSize int64
}

// NewUploadService 创建上传服务（使用本地存储）
func NewUploadService() *UploadService {
	return &UploadService{
		storage: storage.NewLocalStorage(storage.LocalConfig{
			UploadDir: "uploads",
			BaseURL:   "/uploads",
		}),
		maxSize: 10 * 1024 * 1024, // 10MB
	}
}

// NewUploadServiceWithStorage 创建上传服务（指定存储后端）
func NewUploadServiceWithStorage(s storage.Storage) *UploadService {
	return &UploadService{
		storage: s,
		maxSize: 10 * 1024 * 1024, // 10MB
	}
}

// UploadResult 上传结果
type UploadResult struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// Upload 上传文件
func (s *UploadService) Upload(userID int64, file io.Reader, filename string, category string) (*UploadResult, error) {
	// 检查文件格式
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	if !allowedExts[ext] {
		return nil, bizerr.New(40001, "不支持的文件格式，仅支持 jpg/jpeg/png/gif/webp")
	}

	// 读取文件内容到缓冲区（用于检查大小）
	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, file)
	if err != nil {
		return nil, bizerr.ErrInternal
	}

	// 检查文件大小
	if size > s.maxSize {
		return nil, bizerr.New(40002, fmt.Sprintf("文件大小超过限制（最大 %dMB）", s.maxSize/1024/1024))
	}

	// 生成存储 key
	newFilename := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixMilli(), ext)
	key := fmt.Sprintf("%s/%s", category, newFilename)

	// 根据文件扩展名确定 Content-Type
	contentType := getContentType(ext)

	// 上传到存储
	url, err := s.storage.Upload(context.Background(), key, buf, size, contentType)
	if err != nil {
		return nil, bizerr.ErrInternal
	}

	return &UploadResult{
		URL:      url,
		Filename: newFilename,
		Size:     size,
	}, nil
}

// Delete 删除文件
func (s *UploadService) Delete(url string) error {
	// 从存储中删除
	return s.storage.Delete(context.Background(), url)
}

// getContentType 根据扩展名获取 Content-Type
func getContentType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
