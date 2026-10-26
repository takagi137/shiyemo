package dailyimage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shiyemo/server/internal/common/cache"
)

// Service 风景图服务
type Service struct {
	repo RepoInterface
}

// NewService 创建服务
func NewService() *Service {
	return &Service{repo: NewRepo()}
}

// NewServiceWithDeps 创建服务（注入仓储，用于测试）
func NewServiceWithDeps(repo RepoInterface) *Service {
	return &Service{repo: repo}
}

// RandomImage 随机获取一张图片（带缓存）
func (s *Service) RandomImage() (*DailyImage, error) {
	ctx := context.Background()
	cacheKey := "dailyimage:random:" + time.Now().Format("2006-01-02-15")

	// 尝试缓存
	var cached DailyImage
	if err := cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	img, err := s.repo.RandomActive()
	if err != nil {
		return nil, err
	}

	// 缓存 1 小时
	if err := cache.Set(ctx, cacheKey, img, 1*time.Hour); err != nil {
		log.Printf("Failed to cache daily image: %v", err)
	}

	return img, nil
}

// RandomImageFresh 每次真正随机取一张，不走缓存
// 用于「换一张」刷新：走缓存会导致一小时内始终返回同一张。
func (s *Service) RandomImageFresh() (*DailyImage, error) {
	return s.repo.RandomActive()
}

// List 获取图片列表
func (s *Service) List(page, pageSize int) ([]DailyImage, int64, error) {
	return s.repo.List(page, pageSize, false)
}

// Create 创建图片
func (s *Service) Create(url, description, credit string, sortOrder int) (*DailyImage, error) {
	if url == "" {
		return nil, fmt.Errorf("图片URL不能为空")
	}
	img := &DailyImage{
		URL:         url,
		Description: description,
		Credit:      credit,
		SortOrder:   sortOrder,
		IsActive:    true,
	}
	if err := s.repo.Create(img); err != nil {
		return nil, err
	}
	return img, nil
}

// Update 更新图片
func (s *Service) Update(id int64, url, description, credit string, sortOrder int, isActive bool) error {
	img, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	img.URL = url
	img.Description = description
	img.Credit = credit
	img.SortOrder = sortOrder
	img.IsActive = isActive
	return s.repo.Update(img)
}

// Delete 删除图片
func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Reorder 批量调整排序
func (s *Service) Reorder(ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("图片ID列表不能为空")
	}
	return s.repo.UpdateSortOrder(ids)
}
