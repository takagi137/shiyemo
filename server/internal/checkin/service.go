package checkin

import (
	"encoding/json"
	"errors"
	"time"

	bizerr "github.com/shiyemo/server/internal/common/errors"
	"gorm.io/gorm"
)

// CheckinService 打卡服务
type CheckinService struct {
	repo CheckinRepoInterface
}

// NewCheckinService 创建打卡服务
func NewCheckinService() *CheckinService {
	return &CheckinService{
		repo: NewCheckinRepo(),
	}
}

// NewCheckinServiceWithDeps 创建打卡服务（注入依赖，用于测试）
func NewCheckinServiceWithDeps(repo CheckinRepoInterface) *CheckinService {
	return &CheckinService{
		repo: repo,
	}
}

// CreateRequest 创建打卡请求
type CreateRequest struct {
	Mood     string   `json:"mood" binding:"required,oneof=happy calm sad angry anxious"`
	Content  string   `json:"content"`
	Images   []string `json:"images"`
	IsPublic bool     `json:"is_public"`
}

// Create 创建打卡
func (s *CheckinService) Create(userID int64, req *CreateRequest) (*DailyCheckin, error) {
	// 检查今天是否已打卡
	_, err := s.repo.GetTodayByUserID(userID)
	if err == nil {
		return nil, bizerr.New(40201, "今天已经打过卡了")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, bizerr.ErrDB
	}

	// 将图片列表转为 JSON 字符串
	imagesJSON := ""
	if len(req.Images) > 0 {
		imagesBytes, err := json.Marshal(req.Images)
		if err != nil {
			return nil, bizerr.ErrInternal
		}
		imagesJSON = string(imagesBytes)
	}

	checkin := &DailyCheckin{
		UserID:   userID,
		Mood:     req.Mood,
		Content:  req.Content,
		Images:   imagesJSON,
		IsPublic: req.IsPublic,
	}

	if err := s.repo.Create(checkin); err != nil {
		return nil, bizerr.ErrDB
	}

	return checkin, nil
}

// GetToday 获取今天的打卡记录
func (s *CheckinService) GetToday(userID int64) (*DailyCheckin, error) {
	checkin, err := s.repo.GetTodayByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 今天未打卡，返回 nil
		}
		return nil, bizerr.ErrDB
	}
	return checkin, nil
}

// GetHistory 获取打卡历史
func (s *CheckinService) GetHistory(userID int64, page, pageSize int) ([]DailyCheckin, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	return s.repo.GetHistoryByUserID(userID, page, pageSize)
}

// MoodCurveItem 心情曲线数据项
type MoodCurveItem struct {
	Date  string `json:"date"`
	Mood  string `json:"mood"`
	Emoji string `json:"emoji"`
}

// GetMoodCurve 获取近 30 天心情曲线
func (s *CheckinService) GetMoodCurve(userID int64) ([]MoodCurveItem, error) {
	checkins, err := s.repo.GetMoodCurve(userID, 30)
	if err != nil {
		return nil, bizerr.ErrDB
	}

	// 构建日期到心情的映射
	moodMap := make(map[string]string)
	for _, c := range checkins {
		date := c.CreatedAt.Format("2006-01-02")
		moodMap[date] = c.Mood
	}

	// 生成近 30 天的数据（包括未打卡的日期）
	var items []MoodCurveItem
	now := time.Now()
	for i := 29; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		mood, exists := moodMap[date]
		if !exists {
			mood = "none"
		}
		items = append(items, MoodCurveItem{
			Date:  date,
			Mood:  mood,
			Emoji: MoodEmoji[mood],
		})
	}

	return items, nil
}

// GetPublicTimeline 获取公开打卡时间线
func (s *CheckinService) GetPublicTimeline(page, pageSize int) ([]DailyCheckin, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	return s.repo.GetPublicTimeline(page, pageSize)
}
