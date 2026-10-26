package moderation

import (
	"errors"
	"strings"
)

// MockModerationRepo 审核仓储Mock实现
type MockModerationRepo struct {
	words    map[int64]*SensitiveWord
	nextID   int64
	getErr   error
	createErr error
	deleteErr error
}

// NewMockModerationRepo 创建Mock仓储
func NewMockModerationRepo() *MockModerationRepo {
	return &MockModerationRepo{
		words:  make(map[int64]*SensitiveWord),
		nextID: 1,
	}
}

// SetGetError 设置查询错误
func (m *MockModerationRepo) SetGetError(err error) {
	m.getErr = err
}

// SetCreateError 设置创建错误
func (m *MockModerationRepo) SetCreateError(err error) {
	m.createErr = err
}

// SetDeleteError 设置删除错误
func (m *MockModerationRepo) SetDeleteError(err error) {
	m.deleteErr = err
}

// CreateWord 创建敏感词
func (m *MockModerationRepo) CreateWord(word *SensitiveWord) error {
	if m.createErr != nil {
		return m.createErr
	}
	word.ID = m.nextID
	m.nextID++
	m.words[word.ID] = word
	return nil
}

// GetWordByID 根据ID获取敏感词
func (m *MockModerationRepo) GetWordByID(id int64) (*SensitiveWord, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if word, ok := m.words[id]; ok {
		return word, nil
	}
	return nil, errors.New("not found")
}

// GetWordByName 根据名称获取敏感词
func (m *MockModerationRepo) GetWordByName(name string) (*SensitiveWord, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	for _, word := range m.words {
		if strings.EqualFold(word.Word, name) {
			return word, nil
		}
	}
	return nil, nil
}

// ListWords 获取敏感词列表
func (m *MockModerationRepo) ListWords(page, size int) ([]SensitiveWord, int64, error) {
	if m.getErr != nil {
		return nil, 0, m.getErr
	}
	total := int64(len(m.words))
	start := (page - 1) * size
	end := start + size
	if start > len(m.words) {
		return []SensitiveWord{}, total, nil
	}
	if end > len(m.words) {
		end = len(m.words)
	}
	var result []SensitiveWord
	i := 0
	for _, word := range m.words {
		if i >= start && i < end {
			result = append(result, *word)
		}
		i++
	}
	return result, total, nil
}

// DeleteWord 删除敏感词
func (m *MockModerationRepo) DeleteWord(id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.words[id]; !ok {
		return errors.New("not found")
	}
	delete(m.words, id)
	return nil
}

// GetAllWords 获取所有敏感词
func (m *MockModerationRepo) GetAllWords() ([]string, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var result []string
	for _, word := range m.words {
		result = append(result, word.Word)
	}
	return result, nil
}
