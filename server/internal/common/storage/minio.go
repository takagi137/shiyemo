package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOStorage MinIO 对象存储实现
type MinIOStorage struct {
	client  *minio.Client
	bucket  string
	baseURL string
}

// NewMinIOStorage 创建 MinIO 存储
func NewMinIOStorage(cfg MinIOConfig) (*MinIOStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// 确保存储桶存在
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MinIOStorage{
		client:  client,
		bucket:  cfg.Bucket,
		baseURL: cfg.BaseURL,
	}, nil
}

// Upload 上传文件到 MinIO
func (s *MinIOStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
	// 确保 key 以 / 开头
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}

	_, err := s.client.PutObject(ctx, s.bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	return s.GetURL(key), nil
}

// Delete 从 MinIO 删除文件
func (s *MinIOStorage) Delete(ctx context.Context, key string) error {
	// 确保 key 以 / 开头
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}

	return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}

// GetURL 获取文件访问 URL
func (s *MinIOStorage) GetURL(key string) string {
	// 确保 key 不以 / 开头（URL 拼接）
	key = strings.TrimPrefix(key, "/")
	return fmt.Sprintf("%s/%s/%s", s.baseURL, s.bucket, key)
}
