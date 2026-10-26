package tips

import (
	"errors"

	bizerr "github.com/shiyemo/server/internal/common/errors"
	"gorm.io/gorm"
)

// TipsService 贴士服务
type TipsService struct {
	repo TipsRepoInterface
}

// NewTipsService 创建贴士服务
func NewTipsService() *TipsService {
	return &TipsService{
		repo: NewTipsRepo(),
	}
}

// NewTipsServiceWithDeps 创建贴士服务（注入依赖，用于测试）
func NewTipsServiceWithDeps(repo TipsRepoInterface) *TipsService {
	return &TipsService{
		repo: repo,
	}
}

// GetRandom 获取随机贴士
func (s *TipsService) GetRandom() (*SurvivalTip, error) {
	tip, err := s.repo.GetRandom()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerr.ErrParamInvalid
		}
		return nil, bizerr.ErrDB
	}
	return tip, nil
}

// GetList 获取贴士列表
func (s *TipsService) GetList(page, pageSize int) ([]SurvivalTip, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	return s.repo.GetList(page, pageSize)
}

// Create 创建贴士（管理员）
func (s *TipsService) Create(content, category string) (*SurvivalTip, error) {
	tip := &SurvivalTip{
		Content:  content,
		Category: category,
		IsActive: true,
	}

	if err := s.repo.Create(tip); err != nil {
		return nil, bizerr.ErrDB
	}

	return tip, nil
}

// Update 更新贴士（管理员）
func (s *TipsService) Update(id int64, content, category string) error {
	// TODO: 先查询是否存在
	tip := &SurvivalTip{
		ID:       id,
		Content:  content,
		Category: category,
	}

	if err := s.repo.Update(tip); err != nil {
		return bizerr.ErrDB
	}

	return nil
}

// Delete 删除贴士（管理员）
func (s *TipsService) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return bizerr.ErrDB
	}
	return nil
}
