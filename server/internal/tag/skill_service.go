package tag

import (
	bizerr "github.com/shiyemo/server/internal/common/errors"
)

// SkillTagService 用户技能标签服务
type SkillTagService struct {
	repo TagRepository
}

// NewSkillTagService 创建技能标签服务
func NewSkillTagService() *SkillTagService {
	return &SkillTagService{
		repo: NewTagRepo(),
	}
}

// NewSkillTagServiceWithDeps 创建技能标签服务（注入依赖，用于测试）
func NewSkillTagServiceWithDeps(repo TagRepository) *SkillTagService {
	return &SkillTagService{
		repo: repo,
	}
}

// AddUserSkill 添加用户技能标签
func (s *SkillTagService) AddUserSkill(userID int64, tagName string) (*SkillTag, error) {
	// 获取或创建标签
	tag, err := s.repo.GetOrCreate(tagName)
	if err != nil {
		return nil, bizerr.ErrDB
	}

	// 关联用户
	if err := s.repo.AddUserSkill(userID, tag.ID); err != nil {
		return nil, bizerr.ErrDB
	}

	return &SkillTag{
		ID:   tag.ID,
		Name: tag.Name,
	}, nil
}

// RemoveUserSkill 移除用户技能标签
func (s *SkillTagService) RemoveUserSkill(userID int64, tagName string) error {
	tag, err := s.repo.GetByName(tagName)
	if err != nil {
		return bizerr.ErrDB
	}
	if tag == nil {
		return bizerr.ErrTagNotFound
	}

	if err := s.repo.RemoveUserSkill(userID, tag.ID); err != nil {
		return bizerr.ErrDB
	}

	return nil
}

// GetUserSkills 获取用户技能标签列表
func (s *SkillTagService) GetUserSkills(userID int64) ([]SkillTag, error) {
	tags, err := s.repo.GetUserSkills(userID)
	if err != nil {
		return nil, bizerr.ErrDB
	}

	var result []SkillTag
	for _, tag := range tags {
		result = append(result, SkillTag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return result, nil
}

// BatchAddUserSkills 批量添加用户技能标签
func (s *SkillTagService) BatchAddUserSkills(userID int64, tagNames []string) ([]SkillTag, error) {
	var result []SkillTag

	for _, name := range tagNames {
		tag, bizErr := s.AddUserSkill(userID, name)
		if bizErr != nil {
			return nil, bizErr
		}
		result = append(result, *tag)
	}

	return result, nil
}

// SearchBySkill 根据技能标签搜索用户
func (s *SkillTagService) SearchBySkill(tagName string) ([]int64, error) {
	tag, err := s.repo.GetByName(tagName)
	if err != nil {
		return nil, bizerr.ErrDB
	}
	if tag == nil {
		return []int64{}, nil
	}

	userIDs, err := s.repo.GetUserIDsByTag(tag.ID)
	if err != nil {
		return nil, bizerr.ErrDB
	}

	return userIDs, nil
}
