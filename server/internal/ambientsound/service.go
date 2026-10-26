package ambientsound

import "fmt"

// Service 环境音服务
type Service struct {
	repo RepoInterface
}

// NewService 创建服务
func NewService() *Service {
	return &Service{repo: NewRepo()}
}

// NewServiceWithDeps 创建服务（注入仓储，用于测试）
func NewServiceWithDeps(repo RepoInterface) *Service {
	return &Service{repo: repo}
}

// ListActive 获取启用音效列表（公开）
func (s *Service) ListActive() ([]AmbientSound, error) {
	return s.repo.ListActive()
}

// List 获取音效列表（管理后台）
func (s *Service) List(page, pageSize int) ([]AmbientSound, int64, error) {
	return s.repo.List(page, pageSize, false)
}

// Create 创建音效
func (s *Service) Create(key, name, emoji, url, credit string, sortOrder int) (*AmbientSound, error) {
	if key == "" || url == "" {
		return nil, fmt.Errorf("key 与 url 不能为空")
	}
	sound := &AmbientSound{
		Key:       key,
		Name:      name,
		Emoji:     emoji,
		URL:       url,
		Credit:    credit,
		SortOrder: sortOrder,
		IsActive:  true,
	}
	if err := s.repo.Create(sound); err != nil {
		return nil, err
	}
	return sound, nil
}

// Update 更新音效
func (s *Service) Update(id int64, key, name, emoji, url, credit string, sortOrder int, isActive bool) error {
	sound, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	sound.Key = key
	sound.Name = name
	sound.Emoji = emoji
	sound.URL = url
	sound.Credit = credit
	sound.SortOrder = sortOrder
	sound.IsActive = isActive
	return s.repo.Update(sound)
}

// Delete 删除音效
func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Reorder 批量调整排序
func (s *Service) Reorder(ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("音效ID列表不能为空")
	}
	return s.repo.UpdateSortOrder(ids)
}
