package moderation

import bizerr "github.com/shiyemo/server/internal/common/errors"

// Service 审核管理服务
type Service struct {
	repo ModerationRepository
}

// NewService 创建审核管理服务
func NewService() *Service {
	return &Service{
		repo: NewModerationRepo(),
	}
}

// NewServiceWithDeps 创建审核管理服务（注入依赖，用于测试）
func NewServiceWithDeps(repo ModerationRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// AddWord 添加敏感词
func (s *Service) AddWord(word string) (*SensitiveWord, error) {
	// 检查是否已存在
	existing, err := s.repo.GetWordByName(word)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	// 创建新敏感词
	newWord := &SensitiveWord{
		Word: normalize(word),
	}
	if err := s.repo.CreateWord(newWord); err != nil {
		return nil, bizerr.ErrDB
	}

	return newWord, nil
}

// DeleteWord 删除敏感词
func (s *Service) DeleteWord(id int64) error {
	return s.repo.DeleteWord(id)
}

// ListWords 获取敏感词列表
func (s *Service) ListWords(page, pageSize int) ([]SensitiveWord, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 50
	}

	return s.repo.ListWords(page, pageSize)
}
