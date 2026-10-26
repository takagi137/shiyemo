package tips

import (
	"context"
	"log"
	"time"

	"github.com/shiyemo/server/internal/common/cache"
	"github.com/shiyemo/server/internal/dailyimage"
)

// DailyQuote 温暖文案
type DailyQuote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

// DailyResponse 每日内容响应
type DailyResponse struct {
	Image dailyimage.DailyImage `json:"image"`
	Quote DailyQuote            `json:"quote"`
}

// GetDailyContent 获取每日内容（带缓存）
func (s *TipsService) GetDailyContent() (*DailyResponse, error) {
	return s.GetDailyContentWithImageService(dailyimage.NewService())
}

// GetDailyContentWithImageService 获取每日内容（可注入图片服务，用于测试）
func (s *TipsService) GetDailyContentWithImageService(imageSvc *dailyimage.Service) (*DailyResponse, error) {
	ctx := context.Background()
	cacheKey := "daily:content:" + time.Now().Format("2006-01-02")

	// 尝试从缓存获取
	var cached DailyResponse
	if err := cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	// 获取随机文案
	tip, err := s.repo.GetRandom()
	if err != nil {
		// 如果文案表为空，使用默认文案
		return &DailyResponse{
			Image: dailyimage.DailyImage{
				URL:         "https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=1920&q=80",
				Description: "山谷湖泊",
				Credit:      "Photo by Baechtold on Unsplash",
			},
			Quote: DailyQuote{
				Text:   "今天的你，已经比昨天勇敢了一点。",
				Author: "",
			},
		}, nil
	}

	// 从数据库获取随机图片
	image, err := imageSvc.RandomImage()
	if err != nil {
		// 数据库为空时使用默认图片
		image = &dailyimage.DailyImage{
			URL:         "https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=1920&q=80",
			Description: "山谷湖泊",
			Credit:      "Photo by Baechtold on Unsplash",
		}
	}

	result := &DailyResponse{
		Image: *image,
		Quote: DailyQuote{
			Text: tip.Content,
			// category 是内部分类（mindset/finance/...），不作为作者展示
			Author: "",
		},
	}

	// 写入缓存，TTL = 24小时
	if err := cache.Set(ctx, cacheKey, result, 24*time.Hour); err != nil {
		log.Printf("Failed to cache daily content: %v", err)
	}

	return result, nil
}

// RandomQuote 随机获取一条文案（换一条功能）
func (s *TipsService) RandomQuote() (*DailyQuote, error) {
	tip, err := s.repo.GetRandom()
	if err != nil {
		return &DailyQuote{
			Text:   "今天的你，已经比昨天勇敢了一点。",
			Author: "",
		}, nil
	}
	return &DailyQuote{
		Text:   tip.Content,
		Author: "",
	}, nil
}
