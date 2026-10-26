package storage

import (
	"context"
	"io"
)

// Storage 对象存储接口
type Storage interface {
	// Upload 上传文件，返回访问 URL
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error)
	// Delete 删除文件
	Delete(ctx context.Context, key string) error
	// GetURL 获取文件访问 URL
	GetURL(key string) string
}

// Config 存储配置
type Config struct {
	Type      string // local / minio
	Local     LocalConfig
	MinIO     MinIOConfig
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	UploadDir string
	BaseURL   string
}

// MinIOConfig MinIO 配置
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	BaseURL   string // 外部访问基础 URL
}
