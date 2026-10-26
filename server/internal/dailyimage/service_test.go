package dailyimage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockImageRepo mock 仓储
type MockImageRepo struct {
	mock.Mock
}

func (m *MockImageRepo) RandomActive() (*DailyImage, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DailyImage), args.Error(1)
}

func (m *MockImageRepo) List(page, pageSize int, onlyActive bool) ([]DailyImage, int64, error) {
	args := m.Called(page, pageSize, onlyActive)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]DailyImage), args.Get(1).(int64), args.Error(2)
}

func (m *MockImageRepo) GetByID(id int64) (*DailyImage, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DailyImage), args.Error(1)
}

func (m *MockImageRepo) Create(img *DailyImage) error {
	return m.Called(img).Error(0)
}

func (m *MockImageRepo) Update(img *DailyImage) error {
	return m.Called(img).Error(0)
}

func (m *MockImageRepo) Delete(id int64) error {
	return m.Called(id).Error(0)
}

func (m *MockImageRepo) UpdateSortOrder(ids []int64) error {
	return m.Called(ids).Error(0)
}

// RandomImageFresh 是「换一张」使用的无缓存随机取图，须每次直达仓储
func TestRandomImageFresh_Success(t *testing.T) {
	mockRepo := new(MockImageRepo)
	svc := NewServiceWithDeps(mockRepo)

	img := &DailyImage{ID: 3, URL: "https://x/3.jpg", Description: "山谷"}
	mockRepo.On("RandomActive").Return(img, nil)

	result, err := svc.RandomImageFresh()
	assert.NoError(t, err)
	assert.Equal(t, "https://x/3.jpg", result.URL)
	mockRepo.AssertCalled(t, "RandomActive")
}

func TestRandomImageFresh_Error(t *testing.T) {
	mockRepo := new(MockImageRepo)
	svc := NewServiceWithDeps(mockRepo)

	mockRepo.On("RandomActive").Return(nil, errors.New("empty table"))

	_, err := svc.RandomImageFresh()
	assert.Error(t, err)
}

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockImageRepo)
	svc := NewServiceWithDeps(mockRepo)

	mockRepo.On("Create", mock.AnythingOfType("*dailyimage.DailyImage")).Return(nil)

	img, err := svc.Create("https://x/1.jpg", "湖泊", "Unsplash", 1)
	assert.NoError(t, err)
	assert.Equal(t, "https://x/1.jpg", img.URL)
	assert.True(t, img.IsActive)
}

func TestCreate_EmptyURL(t *testing.T) {
	mockRepo := new(MockImageRepo)
	svc := NewServiceWithDeps(mockRepo)

	_, err := svc.Create("", "湖泊", "Unsplash", 1)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestUpdate_Success(t *testing.T) {
	mockRepo := new(MockImageRepo)
	svc := NewServiceWithDeps(mockRepo)

	existing := &DailyImage{ID: 1, URL: "https://x/1.jpg"}
	mockRepo.On("GetByID", int64(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*dailyimage.DailyImage")).Return(nil)

	err := svc.Update(1, "https://x/1b.jpg", "新描述", "Credit", 2, true)
	assert.NoError(t, err)
	assert.Equal(t, "https://x/1b.jpg", existing.URL)
	assert.Equal(t, "新描述", existing.Description)
}

func TestUpdate_NotFound(t *testing.T) {
	mockRepo := new(MockImageRepo)
	svc := NewServiceWithDeps(mockRepo)

	mockRepo.On("GetByID", int64(99)).Return(nil, errors.New("not found"))

	err := svc.Update(99, "https://x/1.jpg", "d", "c", 1, true)
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockImageRepo)
	svc := NewServiceWithDeps(mockRepo)

	mockRepo.On("Delete", int64(1)).Return(nil)

	err := svc.Delete(1)
	assert.NoError(t, err)
}

func TestReorder_Success(t *testing.T) {
	mockRepo := new(MockImageRepo)
	svc := NewServiceWithDeps(mockRepo)

	mockRepo.On("UpdateSortOrder", []int64{2, 1}).Return(nil)

	err := svc.Reorder([]int64{2, 1})
	assert.NoError(t, err)
}

func TestReorder_EmptyIDs(t *testing.T) {
	mockRepo := new(MockImageRepo)
	svc := NewServiceWithDeps(mockRepo)

	err := svc.Reorder([]int64{})
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "UpdateSortOrder", mock.Anything)
}

func TestList_Success(t *testing.T) {
	mockRepo := new(MockImageRepo)
	svc := NewServiceWithDeps(mockRepo)

	imgs := []DailyImage{{ID: 1, URL: "https://x/1.jpg"}}
	mockRepo.On("List", 1, 20, false).Return(imgs, int64(1), nil)

	result, total, err := svc.List(1, 20)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
}
