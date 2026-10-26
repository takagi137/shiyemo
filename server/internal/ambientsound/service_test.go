package ambientsound

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepo mock 仓储
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) ListActive() ([]AmbientSound, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]AmbientSound), args.Error(1)
}

func (m *MockRepo) List(page, pageSize int, onlyActive bool) ([]AmbientSound, int64, error) {
	args := m.Called(page, pageSize, onlyActive)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]AmbientSound), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepo) GetByID(id int64) (*AmbientSound, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AmbientSound), args.Error(1)
}

func (m *MockRepo) Create(s *AmbientSound) error {
	return m.Called(s).Error(0)
}

func (m *MockRepo) Update(s *AmbientSound) error {
	return m.Called(s).Error(0)
}

func (m *MockRepo) Delete(id int64) error {
	return m.Called(id).Error(0)
}

func (m *MockRepo) UpdateSortOrder(ids []int64) error {
	return m.Called(ids).Error(0)
}

func TestListActive_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewServiceWithDeps(mockRepo)

	sounds := []AmbientSound{
		{ID: 1, Key: "rain", Name: "雨声", URL: "https://x/rain.mp3"},
		{ID: 2, Key: "ocean", Name: "海浪", URL: "https://x/ocean.mp3"},
	}
	mockRepo.On("ListActive").Return(sounds, nil)

	result, err := svc.ListActive()
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "rain", result[0].Key)
}

func TestListActive_Error(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewServiceWithDeps(mockRepo)

	mockRepo.On("ListActive").Return(nil, errors.New("db error"))

	_, err := svc.ListActive()
	assert.Error(t, err)
}

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewServiceWithDeps(mockRepo)

	mockRepo.On("Create", mock.AnythingOfType("*ambientsound.AmbientSound")).Return(nil)

	sound, err := svc.Create("rain", "雨声", "🌧️", "https://x/rain.mp3", "BigSoundBank (CC0)", 1)
	assert.NoError(t, err)
	assert.Equal(t, "rain", sound.Key)
	assert.True(t, sound.IsActive)
	mockRepo.AssertCalled(t, "Create", mock.AnythingOfType("*ambientsound.AmbientSound"))
}

func TestCreate_EmptyKey(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewServiceWithDeps(mockRepo)

	_, err := svc.Create("", "雨声", "🌧️", "https://x/rain.mp3", "", 1)
	assert.Error(t, err)
	// 校验失败时不应触及仓储
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestCreate_EmptyURL(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewServiceWithDeps(mockRepo)

	_, err := svc.Create("rain", "雨声", "🌧️", "", "", 1)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestUpdate_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewServiceWithDeps(mockRepo)

	existing := &AmbientSound{ID: 1, Key: "rain", Name: "雨声", URL: "https://x/rain.mp3"}
	mockRepo.On("GetByID", int64(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*ambientsound.AmbientSound")).Return(nil)

	err := svc.Update(1, "rain", "雨声(新)", "🌧️", "https://x/rain2.mp3", "CC0", 2, true)
	assert.NoError(t, err)
	assert.Equal(t, "雨声(新)", existing.Name)
	assert.Equal(t, "https://x/rain2.mp3", existing.URL)
}

func TestUpdate_NotFound(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewServiceWithDeps(mockRepo)

	mockRepo.On("GetByID", int64(99)).Return(nil, errors.New("record not found"))

	err := svc.Update(99, "rain", "雨声", "🌧️", "https://x/rain.mp3", "", 1, true)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewServiceWithDeps(mockRepo)

	mockRepo.On("Delete", int64(1)).Return(nil)

	err := svc.Delete(1)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Delete", int64(1))
}

func TestReorder_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewServiceWithDeps(mockRepo)

	mockRepo.On("UpdateSortOrder", []int64{3, 1, 2}).Return(nil)

	err := svc.Reorder([]int64{3, 1, 2})
	assert.NoError(t, err)
}

func TestReorder_EmptyIDs(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewServiceWithDeps(mockRepo)

	err := svc.Reorder([]int64{})
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "UpdateSortOrder", mock.Anything)
}

func TestList_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewServiceWithDeps(mockRepo)

	sounds := []AmbientSound{{ID: 1, Key: "rain"}}
	mockRepo.On("List", 1, 20, false).Return(sounds, int64(1), nil)

	result, total, err := svc.List(1, 20)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
}
