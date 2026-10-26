package storage

import "fmt"

// NewStorage 根据配置创建存储实例
func NewStorage(cfg Config) (Storage, error) {
	switch cfg.Type {
	case "minio":
		return NewMinIOStorage(cfg.MinIO)
	case "local", "":
		return NewLocalStorage(cfg.Local), nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}
