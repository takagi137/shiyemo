package tag

import bizerr "github.com/shiyemo/server/internal/common/errors"

// TagService 标签服务
type TagService struct {
	repo TagRepository
}

// NewTagService 创建标签服务
func NewTagService() *TagService {
	return &TagService{
		repo: NewTagRepo(),
	}
}

// NewTagServiceWithDeps 创建标签服务（注入依赖，用于测试）
func NewTagServiceWithDeps(repo TagRepository) *TagService {
	return &TagService{
		repo: repo,
	}
}

// GetHotTags 获取热门标签
func (s *TagService) GetHotTags(limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 20
	}
	tags, err := s.repo.GetHotTags(limit)
	if err != nil {
		return nil, bizerr.ErrDB
	}
	return tags, nil
}

// UpdatePostTags 更新帖子标签
func (s *TagService) UpdatePostTags(postID int64, tagNames []string) ([]string, error) {
	var tagIDs []int64
	
	for _, name := range tagNames {
		tag, err := s.repo.GetOrCreate(name)
		if err != nil {
			return nil, bizerr.ErrDB
		}
		tagIDs = append(tagIDs, tag.ID)
	}
	
	if err := s.repo.UpdatePostTags(postID, tagIDs); err != nil {
		return nil, bizerr.ErrDB
	}
	
	return tagNames, nil
}

// GetPostTags 获取帖子标签
func (s *TagService) GetPostTags(postID int64) ([]string, error) {
	names, err := s.repo.GetPostTags(postID)
	if err != nil {
		return nil, bizerr.ErrDB
	}
	return names, nil
}
