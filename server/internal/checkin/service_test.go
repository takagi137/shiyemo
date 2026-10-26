package checkin

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockCheckinRepo 模拟打卡仓储
type MockCheckinRepo struct {
	mock.Mock
}

func (m *MockCheckinRepo) Create(checkin *DailyCheckin) error {
	args := m.Called(checkin)
	return args.Error(0)
}

func (m *MockCheckinRepo) GetTodayByUserID(userID int64) (*DailyCheckin, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DailyCheckin), args.Error(1)
}

func (m *MockCheckinRepo) GetHistoryByUserID(userID int64, page, pageSize int) ([]DailyCheckin, int64, error) {
	args := m.Called(userID, page, pageSize)
	return args.Get(0).([]DailyCheckin), args.Get(1).(int64), args.Error(2)
}

func (m *MockCheckinRepo) GetMoodCurve(userID int64, days int) ([]DailyCheckin, error) {
	args := m.Called(userID, days)
	return args.Get(0).([]DailyCheckin), args.Error(1)
}

func (m *MockCheckinRepo) GetPublicTimeline(page, pageSize int) ([]DailyCheckin, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]DailyCheckin), args.Get(1).(int64), args.Error(2)
}

func TestCheckinService_Create_Success(t *testing.T) {
	mockRepo := new(MockCheckinRepo)
	service := NewCheckinServiceWithDeps(mockRepo)

	userID := int64(1)
	req := &CreateRequest{
		Mood:     MoodHappy,
		Content:  "今天心情不错",
		IsPublic: true,
	}

	// 模拟今天未打卡
	mockRepo.On("GetTodayByUserID", userID).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.AnythingOfType("*checkin.DailyCheckin")).Return(nil)

	checkin, err := service.Create(userID, req)
	assert.NoError(t, err)
	assert.NotNil(t, checkin)
	assert.Equal(t, userID, checkin.UserID)
	assert.Equal(t, MoodHappy, checkin.Mood)
	mockRepo.AssertExpectations(t)
}

func TestCheckinService_Create_AlreadyCheckedIn(t *testing.T) {
	mockRepo := new(MockCheckinRepo)
	service := NewCheckinServiceWithDeps(mockRepo)

	userID := int64(1)
	req := &CreateRequest{
		Mood:     MoodHappy,
		Content:  "今天心情不错",
		IsPublic: true,
	}

	// 模拟今天已打卡
	existingCheckin := &DailyCheckin{
		ID:        1,
		UserID:    userID,
		Mood:      MoodCalm,
		CreatedAt: time.Now(),
	}
	mockRepo.On("GetTodayByUserID", userID).Return(existingCheckin, nil)

	checkin, err := service.Create(userID, req)
	assert.Error(t, err)
	assert.Nil(t, checkin)
	assert.Contains(t, err.Error(), "今天已经打过卡了")
	mockRepo.AssertExpectations(t)
}

func TestCheckinService_GetToday_CheckedIn(t *testing.T) {
	mockRepo := new(MockCheckinRepo)
	service := NewCheckinServiceWithDeps(mockRepo)

	userID := int64(1)
	existingCheckin := &DailyCheckin{
		ID:        1,
		UserID:    userID,
		Mood:      MoodHappy,
		CreatedAt: time.Now(),
	}

	mockRepo.On("GetTodayByUserID", userID).Return(existingCheckin, nil)

	checkin, err := service.GetToday(userID)
	assert.NoError(t, err)
	assert.NotNil(t, checkin)
	assert.Equal(t, existingCheckin, checkin)
	mockRepo.AssertExpectations(t)
}

func TestCheckinService_GetToday_NotCheckedIn(t *testing.T) {
	mockRepo := new(MockCheckinRepo)
	service := NewCheckinServiceWithDeps(mockRepo)

	userID := int64(1)

	mockRepo.On("GetTodayByUserID", userID).Return(nil, gorm.ErrRecordNotFound)

	checkin, err := service.GetToday(userID)
	assert.NoError(t, err)
	assert.Nil(t, checkin)
	mockRepo.AssertExpectations(t)
}

func TestCheckinService_GetHistory(t *testing.T) {
	mockRepo := new(MockCheckinRepo)
	service := NewCheckinServiceWithDeps(mockRepo)

	userID := int64(1)
	page := 1
	pageSize := 10

	expectedCheckins := []DailyCheckin{
		{ID: 2, UserID: userID, Mood: MoodCalm, CreatedAt: time.Now()},
		{ID: 1, UserID: userID, Mood: MoodHappy, CreatedAt: time.Now().AddDate(0, 0, -1)},
	}

	mockRepo.On("GetHistoryByUserID", userID, page, pageSize).Return(expectedCheckins, int64(2), nil)

	checkins, total, err := service.GetHistory(userID, page, pageSize)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, checkins, 2)
	mockRepo.AssertExpectations(t)
}

func TestCheckinService_GetHistory_DefaultPagination(t *testing.T) {
	mockRepo := new(MockCheckinRepo)
	service := NewCheckinServiceWithDeps(mockRepo)

	userID := int64(1)

	mockRepo.On("GetHistoryByUserID", userID, 1, 20).Return([]DailyCheckin{}, int64(0), nil)

	// 传入无效的分页参数
	checkins, total, err := service.GetHistory(userID, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, checkins, 0)
	mockRepo.AssertExpectations(t)
}

func TestCheckinService_GetMoodCurve(t *testing.T) {
	mockRepo := new(MockCheckinRepo)
	service := NewCheckinServiceWithDeps(mockRepo)

	userID := int64(1)

	// 模拟最近 2 天有打卡
	mockCheckins := []DailyCheckin{
		{ID: 1, UserID: userID, Mood: MoodHappy, CreatedAt: time.Now().AddDate(0, 0, -2)},
		{ID: 2, UserID: userID, Mood: MoodCalm, CreatedAt: time.Now().AddDate(0, 0, -1)},
	}

	mockRepo.On("GetMoodCurve", userID, 30).Return(mockCheckins, nil)

	items, err := service.GetMoodCurve(userID)
	assert.NoError(t, err)
	assert.Len(t, items, 30) // 应该返回 30 天的数据

	// 检查最近 2 天的数据
	assert.Equal(t, MoodHappy, items[27].Mood) // 28 天前
	assert.Equal(t, MoodCalm, items[28].Mood)  // 1 天前
	assert.Equal(t, "none", items[29].Mood)    // 今天未打卡
	mockRepo.AssertExpectations(t)
}

func TestCheckinService_GetPublicTimeline(t *testing.T) {
	mockRepo := new(MockCheckinRepo)
	service := NewCheckinServiceWithDeps(mockRepo)

	page := 1
	pageSize := 10

	expectedCheckins := []DailyCheckin{
		{ID: 1, UserID: 1, Mood: MoodHappy, IsPublic: true},
	}

	mockRepo.On("GetPublicTimeline", page, pageSize).Return(expectedCheckins, int64(1), nil)

	checkins, total, err := service.GetPublicTimeline(page, pageSize)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, checkins, 1)
	mockRepo.AssertExpectations(t)
}

func TestCheckinService_GetPublicTimeline_DefaultPagination(t *testing.T) {
	mockRepo := new(MockCheckinRepo)
	service := NewCheckinServiceWithDeps(mockRepo)

	mockRepo.On("GetPublicTimeline", 1, 20).Return([]DailyCheckin{}, int64(0), nil)

	// 传入无效的分页参数
	checkins, total, err := service.GetPublicTimeline(0, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, checkins, 0)
	mockRepo.AssertExpectations(t)
}

func TestCheckinService_Create_WithImages(t *testing.T) {
	mockRepo := new(MockCheckinRepo)
	service := NewCheckinServiceWithDeps(mockRepo)

	userID := int64(1)
	req := &CreateRequest{
		Mood:     MoodHappy,
		Content:  "今天心情不错",
		Images:   []string{"http://example.com/1.jpg", "http://example.com/2.jpg"},
		IsPublic: true,
	}

	// 模拟今天未打卡
	mockRepo.On("GetTodayByUserID", userID).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.AnythingOfType("*checkin.DailyCheckin")).Return(nil)

	checkin, err := service.Create(userID, req)
	assert.NoError(t, err)
	assert.NotNil(t, checkin)
	assert.Equal(t, userID, checkin.UserID)
	assert.Equal(t, MoodHappy, checkin.Mood)
	assert.NotEmpty(t, checkin.Images)
	mockRepo.AssertExpectations(t)
}

func TestCheckinService_Create_DBError(t *testing.T) {
	mockRepo := new(MockCheckinRepo)
	service := NewCheckinServiceWithDeps(mockRepo)

	userID := int64(1)
	req := &CreateRequest{
		Mood:     MoodHappy,
		Content:  "今天心情不错",
		IsPublic: true,
	}

	// 模拟今天未打卡，但创建时数据库错误
	mockRepo.On("GetTodayByUserID", userID).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.AnythingOfType("*checkin.DailyCheckin")).Return(errors.New("database error"))

	checkin, err := service.Create(userID, req)
	assert.Error(t, err)
	assert.Nil(t, checkin)
	mockRepo.AssertExpectations(t)
}

func TestCheckinService_GetMoodCurve_NoData(t *testing.T) {
	mockRepo := new(MockCheckinRepo)
	service := NewCheckinServiceWithDeps(mockRepo)

	userID := int64(1)

	mockRepo.On("GetMoodCurve", userID, 30).Return([]DailyCheckin{}, nil)

	items, err := service.GetMoodCurve(userID)
	assert.NoError(t, err)
	assert.Len(t, items, 30)
	// 所有日期的心情都应该是 none
	for _, item := range items {
		assert.Equal(t, "none", item.Mood)
	}
	mockRepo.AssertExpectations(t)
}
