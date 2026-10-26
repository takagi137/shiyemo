package journal

import (
	"errors"

	bizerr "github.com/shiyemo/server/internal/common/errors"
	"gorm.io/gorm"
)

// JournalService 日记服务
type JournalService struct {
	repo JournalRepoInterface
}

// NewJournalService 创建日记服务
func NewJournalService() *JournalService {
	return &JournalService{
		repo: NewJournalRepo(),
	}
}

// NewJournalServiceWithDeps 创建日记服务（注入依赖，用于测试）
func NewJournalServiceWithDeps(repo JournalRepoInterface) *JournalService {
	return &JournalService{
		repo: repo,
	}
}

// Create 创建日记
func (s *JournalService) Create(userID int64, content string) (*Journal, error) {
	if content == "" {
		return nil, bizerr.ErrContentEmpty
	}

	journal := &Journal{
		UserID:  userID,
		Content: content,
	}

	if err := s.repo.Create(journal); err != nil {
		return nil, bizerr.ErrDB
	}

	return journal, nil
}

// GetList 获取用户日记列表
func (s *JournalService) GetList(userID int64, page, pageSize int) ([]Journal, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	return s.repo.GetByUserID(userID, page, pageSize)
}

// Delete 删除日记
func (s *JournalService) Delete(id, userID int64) error {
	// 先检查是否存在且属于当前用户
	_, err := s.repo.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerr.ErrPostNotFound // 复用"不存在"错误码
		}
		return bizerr.ErrDB
	}

	return s.repo.Delete(id, userID)
}
