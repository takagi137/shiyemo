package journal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockJournalRepo mock 仓储
type MockJournalRepo struct {
	mock.Mock
}

func (m *MockJournalRepo) Create(journal *Journal) error {
	args := m.Called(journal)
	return args.Error(0)
}

func (m *MockJournalRepo) GetByUserID(userID int64, page, pageSize int) ([]Journal, int64, error) {
	args := m.Called(userID, page, pageSize)
	return args.Get(0).([]Journal), args.Get(1).(int64), args.Error(2)
}

func (m *MockJournalRepo) GetByID(id, userID int64) (*Journal, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Journal), args.Error(1)
}

func (m *MockJournalRepo) Delete(id, userID int64) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockJournalRepo)
	svc := NewJournalServiceWithDeps(mockRepo)

	mockRepo.On("Create", mock.AnythingOfType("*journal.Journal")).Return(nil)

	journal, err := svc.Create(1, "今天感觉好多了")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), journal.UserID)
	assert.Equal(t, "今天感觉好多了", journal.Content)
	mockRepo.AssertExpectations(t)
}

func TestCreate_EmptyContent(t *testing.T) {
	mockRepo := new(MockJournalRepo)
	svc := NewJournalServiceWithDeps(mockRepo)

	journal, err := svc.Create(1, "")
	assert.Error(t, err)
	assert.Nil(t, journal)
}

func TestGetList_Success(t *testing.T) {
	mockRepo := new(MockJournalRepo)
	svc := NewJournalServiceWithDeps(mockRepo)

	journals := []Journal{{ID: 1, Content: "日记1"}, {ID: 2, Content: "日记2"}}
	mockRepo.On("GetByUserID", int64(1), 1, 20).Return(journals, int64(2), nil)

	result, total, err := svc.GetList(1, 1, 20)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestGetList_DefaultParams(t *testing.T) {
	mockRepo := new(MockJournalRepo)
	svc := NewJournalServiceWithDeps(mockRepo)

	mockRepo.On("GetByUserID", int64(1), 1, 20).Return([]Journal{}, int64(0), nil)

	result, total, err := svc.GetList(1, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, result)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockJournalRepo)
	svc := NewJournalServiceWithDeps(mockRepo)

	existing := &Journal{ID: 1, UserID: 1, Content: "test"}
	mockRepo.On("GetByID", int64(1), int64(1)).Return(existing, nil)
	mockRepo.On("Delete", int64(1), int64(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDelete_NotFound(t *testing.T) {
	mockRepo := new(MockJournalRepo)
	svc := NewJournalServiceWithDeps(mockRepo)

	mockRepo.On("GetByID", int64(999), int64(1)).Return(nil, gorm.ErrRecordNotFound)

	err := svc.Delete(999, 1)
	assert.Error(t, err)
}
