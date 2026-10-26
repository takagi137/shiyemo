package reading

import (
	"errors"

	bizerr "github.com/shiyemo/server/internal/common/errors"
	"gorm.io/gorm"
)

// ReadingService 阅读服务
type ReadingService struct {
	repo ReadingRepoInterface
}

// NewReadingService 创建阅读服务
func NewReadingService() *ReadingService {
	return &ReadingService{
		repo: NewReadingRepo(),
	}
}

// NewReadingServiceWithDeps 创建阅读服务（注入依赖，用于测试）
func NewReadingServiceWithDeps(repo ReadingRepoInterface) *ReadingService {
	return &ReadingService{
		repo: repo,
	}
}

// GetToday 获取今日阅读
func (s *ReadingService) GetToday() (*Reading, error) {
	reading, err := s.repo.GetToday()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// readings 表为空，返回默认内容
			return &Reading{
				Source:         "古人智慧",
				Content:        "古人也曾在风雨中前行，你并不孤单。",
				Interpretation: "每一代人都有属于自己的风雨，而每一代人都走了过来。",
			}, nil
		}
		return nil, bizerr.ErrDB
	}
	return reading, nil
}

// GetRandom 获取随机阅读（排除指定 ID）
func (s *ReadingService) GetRandom(excludeIDs []int64) (*Reading, error) {
	total, err := s.repo.Count()
	if err != nil {
		return nil, bizerr.ErrDB
	}

	if total == 0 {
		return &Reading{
			Source:         "古人智慧",
			Content:        "古人也曾在风雨中前行，你并不孤单。",
			Interpretation: "每一代人都有属于自己的风雨，而每一代人都走了过来。",
		}, nil
	}

	reading, err := s.repo.GetRandom(excludeIDs)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 所有都读过了，重置排除列表
			return s.repo.GetRandom(nil)
		}
		return nil, bizerr.ErrDB
	}
	return reading, nil
}
