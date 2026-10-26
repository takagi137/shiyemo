package reading

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockReadingRepo mock 仓储
type MockReadingRepo struct {
	mock.Mock
}

func (m *MockReadingRepo) GetByID(id int64) (*Reading, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Reading), args.Error(1)
}

func (m *MockReadingRepo) GetToday() (*Reading, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Reading), args.Error(1)
}

func (m *MockReadingRepo) GetRandom(excludeIDs []int64) (*Reading, error) {
	args := m.Called(excludeIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Reading), args.Error(1)
}

func (m *MockReadingRepo) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func TestGetToday_Success(t *testing.T) {
	mockRepo := new(MockReadingRepo)
	svc := NewReadingServiceWithDeps(mockRepo)

	reading := &Reading{ID: 1, Source: "苏轼《定风波》", Content: "莫听穿林打叶声"}
	mockRepo.On("GetToday").Return(reading, nil)

	result, err := svc.GetToday()
	assert.NoError(t, err)
	assert.Equal(t, "苏轼《定风波》", result.Source)
}

func TestGetToday_EmptyTable(t *testing.T) {
	mockRepo := new(MockReadingRepo)
	svc := NewReadingServiceWithDeps(mockRepo)

	mockRepo.On("GetToday").Return(nil, gorm.ErrRecordNotFound)

	result, err := svc.GetToday()
	assert.NoError(t, err)
	assert.Equal(t, "古人也曾在风雨中前行，你并不孤单。", result.Content)
}

func TestGetRandom_Success(t *testing.T) {
	mockRepo := new(MockReadingRepo)
	svc := NewReadingServiceWithDeps(mockRepo)

	mockRepo.On("Count").Return(int64(10), nil)
	reading := &Reading{ID: 3, Content: "静夜思"}
	mockRepo.On("GetRandom", []int64{1, 2}).Return(reading, nil)

	result, err := svc.GetRandom([]int64{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, "静夜思", result.Content)
}

func TestGetRandom_EmptyTable(t *testing.T) {
	mockRepo := new(MockReadingRepo)
	svc := NewReadingServiceWithDeps(mockRepo)

	mockRepo.On("Count").Return(int64(0), nil)

	result, err := svc.GetRandom(nil)
	assert.NoError(t, err)
	assert.Equal(t, "古人也曾在风雨中前行，你并不孤单。", result.Content)
}
