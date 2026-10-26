package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Mock 实现（tag 模块测试共享）
// ============================================================

// MockTagRepo mock 标签仓储
type MockTagRepo struct {
	tags       map[string]*Tag
	postTags   map[int64][]int64
	userSkills map[int64][]int64
	nextID     int64
	hotTags    []map[string]interface{}
	getErr     error
	createErr  error
	updateErr  error
	deleteErr  error
}

// NewMockTagRepo 创建 mock 标签仓储
func NewMockTagRepo() *MockTagRepo {
	return &MockTagRepo{
		tags:       make(map[string]*Tag),
		postTags:   make(map[int64][]int64),
		userSkills: make(map[int64][]int64),
		nextID:     1,
	}
}

// SetHotTags 设置热门标签数据
func (m *MockTagRepo) SetHotTags(tags []map[string]interface{}) {
	m.hotTags = tags
}

// SetGetError 设置查询错误
func (m *MockTagRepo) SetGetError(err error) {
	m.getErr = err
}

// SetCreateError 设置创建错误
func (m *MockTagRepo) SetCreateError(err error) {
	m.createErr = err
}

// SetUpdateError 设置更新错误
func (m *MockTagRepo) SetUpdateError(err error) {
	m.updateErr = err
}

// SetDeleteError 设置删除错误
func (m *MockTagRepo) SetDeleteError(err error) {
	m.deleteErr = err
}

// GetHotTags 获取热门标签
func (m *MockTagRepo) GetHotTags(limit int) ([]map[string]interface{}, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if limit > len(m.hotTags) {
		return m.hotTags, nil
	}
	return m.hotTags[:limit], nil
}

// GetOrCreate 获取或创建标签
func (m *MockTagRepo) GetOrCreate(name string) (*Tag, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.createErr != nil {
		return nil, m.createErr
	}

	if tag, ok := m.tags[name]; ok {
		return tag, nil
	}

	tag := &Tag{ID: m.nextID, Name: name}
	m.nextID++
	m.tags[name] = tag
	return tag, nil
}

// GetByName 根据名称获取标签
func (m *MockTagRepo) GetByName(name string) (*Tag, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if tag, ok := m.tags[name]; ok {
		return tag, nil
	}
	return nil, nil
}

// AddUserSkill 添加用户技能标签
func (m *MockTagRepo) AddUserSkill(userID int64, tagID int64) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.userSkills[userID] = append(m.userSkills[userID], tagID)
	return nil
}

// RemoveUserSkill 移除用户技能标签
func (m *MockTagRepo) RemoveUserSkill(userID int64, tagID int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	tags := m.userSkills[userID]
	for i, id := range tags {
		if id == tagID {
			m.userSkills[userID] = append(tags[:i], tags[i+1:]...)
			break
		}
	}
	return nil
}

// GetUserSkills 获取用户技能标签
func (m *MockTagRepo) GetUserSkills(userID int64) ([]Tag, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var result []Tag
	for _, tagID := range m.userSkills[userID] {
		for _, tag := range m.tags {
			if tag.ID == tagID {
				result = append(result, *tag)
				break
			}
		}
	}
	return result, nil
}

// GetUserIDsByTag 根据标签获取用户ID列表
func (m *MockTagRepo) GetUserIDsByTag(tagID int64) ([]int64, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var result []int64
	for userID, tagIDs := range m.userSkills {
		for _, id := range tagIDs {
			if id == tagID {
				result = append(result, userID)
				break
			}
		}
	}
	return result, nil
}

// UpdatePostTags 更新帖子标签
func (m *MockTagRepo) UpdatePostTags(postID int64, tagIDs []int64) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.postTags[postID] = tagIDs
	return nil
}

// GetPostTags 获取帖子标签
func (m *MockTagRepo) GetPostTags(postID int64) ([]string, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}

	tagIDs, ok := m.postTags[postID]
	if !ok {
		return nil, nil
	}

	var names []string
	for _, tagID := range tagIDs {
		for _, tag := range m.tags {
			if tag.ID == tagID {
				names = append(names, tag.Name)
				break
			}
		}
	}
	return names, nil
}

// ============================================================
// MockTagRepo 基础测试
// ============================================================

func TestMockTagRepo(t *testing.T) {
	t.Run("GetOrCreate", func(t *testing.T) {
		repo := NewMockTagRepo()

		tag1, err := repo.GetOrCreate("Go")
		require.NoError(t, err)
		assert.Equal(t, "Go", tag1.Name)
		assert.Equal(t, int64(1), tag1.ID)

		tag2, err := repo.GetOrCreate("Go")
		require.NoError(t, err)
		assert.Equal(t, tag1.ID, tag2.ID)
	})

	t.Run("GetHotTags", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.SetHotTags([]map[string]interface{}{
			{"id": 1, "name": "Go", "count": 10},
			{"id": 2, "name": "Python", "count": 8},
		})

		tags, err := repo.GetHotTags(10)
		require.NoError(t, err)
		assert.Len(t, tags, 2)
	})

	t.Run("UpdatePostTags", func(t *testing.T) {
		repo := NewMockTagRepo()

		err := repo.UpdatePostTags(1, []int64{1, 2})
		require.NoError(t, err)
		assert.Len(t, repo.postTags[1], 2)
	})

	t.Run("GetPostTags", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.GetOrCreate("Go")
		repo.GetOrCreate("Python")
		repo.UpdatePostTags(1, []int64{1, 2})

		names, err := repo.GetPostTags(1)
		require.NoError(t, err)
		assert.Len(t, names, 2)
	})

	t.Run("GetByName", func(t *testing.T) {
		repo := NewMockTagRepo()
		repo.GetOrCreate("Go")

		tag, err := repo.GetByName("Go")
		require.NoError(t, err)
		assert.NotNil(t, tag)
		assert.Equal(t, "Go", tag.Name)

		nilTag, err := repo.GetByName("NotExist")
		require.NoError(t, err)
		assert.Nil(t, nilTag)
	})

	t.Run("UserSkills", func(t *testing.T) {
		repo := NewMockTagRepo()
		tag, _ := repo.GetOrCreate("Go")

		err := repo.AddUserSkill(1, tag.ID)
		require.NoError(t, err)

		skills, err := repo.GetUserSkills(1)
		require.NoError(t, err)
		assert.Len(t, skills, 1)

		userIDs, err := repo.GetUserIDsByTag(tag.ID)
		require.NoError(t, err)
		assert.Contains(t, userIDs, int64(1))

		err = repo.RemoveUserSkill(1, tag.ID)
		require.NoError(t, err)
		skills, _ = repo.GetUserSkills(1)
		assert.Len(t, skills, 0)
	})
}
