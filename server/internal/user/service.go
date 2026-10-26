package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/cache"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	repo UserRepoInterface
	sf   *cache.Singleflight
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{
		repo: NewUserRepo(),
		sf:   cache.NewSingleflight(),
	}
}

// NewUserServiceWithRepo 使用注入的 repo 创建用户服务（测试用）
func NewUserServiceWithRepo(repo UserRepoInterface) *UserService {
	return &UserService{repo: repo, sf: cache.NewSingleflight()}
}

// GetProfile 获取用户完整资料
func (s *UserService) GetProfile(id int64) (*auth.User, error) {
	ctx := context.Background()
	cacheKey := cache.FormatUserDetailKey(id)

	// 尝试从缓存获取
	var cachedUser auth.User
	if err := cache.Get(ctx, cacheKey, &cachedUser); err == nil {
		return &cachedUser, nil
	}

	// 缓存未命中，使用 singleflight 防止缓存雪崩
	v, err := s.sf.Do(cacheKey, func() (interface{}, error) {
		user, err := s.repo.GetByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, bizerr.ErrUserNotFound
			}
			return nil, bizerr.ErrDB
		}
		return user, nil
	})
	if err != nil {
		return nil, err
	}
	user := v.(*auth.User)

	// 写入缓存
	if err := cache.Set(ctx, cacheKey, user, cache.CacheExpireUserDetail); err != nil {
		log.Printf("Failed to cache user profile %d: %v", id, err)
	}

	return user, nil
}

// GetPublicProfile 获取用户公开资料
func (s *UserService) GetPublicProfile(id int64) (map[string]interface{}, error) {
	ctx := context.Background()
	cacheKey := cache.FormatUserDetailKey(id) + ":public"

	// 尝试从缓存获取
	var cachedProfile map[string]interface{}
	if err := cache.Get(ctx, cacheKey, &cachedProfile); err == nil {
		log.Printf("Cache hit for user public profile: %d", id)
		return cachedProfile, nil
	}

	// 缓存未命中，从数据库获取
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerr.ErrUserNotFound
		}
		return nil, bizerr.ErrDB
	}

	// 只返回公开信息
	result := map[string]interface{}{
		"id":         user.ID,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"bio":        user.Bio,
		"company":    user.Company,
		"occupation": user.Occupation,
		"city":       user.City,
		"industry":   user.Industry,
		"created_at": user.CreatedAt,
	}

	// 离职时间
	if user.LeftJobAt != nil {
		result["left_job_at"] = user.LeftJobAt
	}

	// 计算自由天数
	freedomDays := user.FreedomDays()
	if freedomDays != nil {
		result["freedom_days"] = *freedomDays
	}

	// 写入缓存
	if err := cache.Set(ctx, cacheKey, result, cache.CacheExpireUserDetail); err != nil {
		log.Printf("Failed to cache user public profile %d: %v", id, err)
	}

	return result, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(id int64, req UpdateProfileRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerr.ErrUserNotFound
		}
		return bizerr.ErrDB
	}

	// 更新字段
	fields := make(map[string]interface{})
	if req.Nickname != "" {
		fields["nickname"] = req.Nickname
	}
	if req.Bio != "" {
		fields["bio"] = req.Bio
	}
	if req.Company != "" {
		fields["company"] = req.Company
	}
	if req.Occupation != "" {
		fields["occupation"] = req.Occupation
	}
	if req.City != "" {
		fields["city"] = req.City
	}
	if req.Industry != "" {
		fields["industry"] = req.Industry
	}
	if req.LeftJobAt != nil {
		if *req.LeftJobAt == "" {
			fields["left_job_at"] = nil
		} else {
			parsedTime, err := time.Parse("2006-01-02", *req.LeftJobAt)
			if err != nil {
				return bizerr.ErrParamInvalid
			}
			fields["left_job_at"] = parsedTime
		}
	}
	if req.LastWorkDate != nil {
		if *req.LastWorkDate == "" {
			fields["last_work_date"] = nil
		} else {
			parsedTime, err := time.Parse("2006-01-02", *req.LastWorkDate)
			if err != nil {
				return bizerr.ErrParamInvalid
			}
			fields["last_work_date"] = parsedTime
		}
	}

	if len(fields) > 0 {
		if err := s.repo.UpdateFields(user.ID, fields); err != nil {
			return err
		}
		// 清除用户缓存
		s.invalidateUserCache(id)
	}
	return nil
}

// UploadAvatar 上传头像
func (s *UserService) UploadAvatar(userID int64, file io.Reader, filename string) (string, error) {
	// 检查文件格式
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return "", bizerr.ErrParamInvalid
	}

	// 创建上传目录
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", bizerr.ErrInternal
	}

	// 生成文件名
	newFilename := fmt.Sprintf("%d%s", userID, ext)
	filePath := filepath.Join(uploadDir, newFilename)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", bizerr.ErrInternal
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", bizerr.ErrInternal
	}

	// 更新用户头像字段
	avatarURL := "/" + filePath
	if err := s.repo.UpdateFields(userID, map[string]interface{}{"avatar": avatarURL}); err != nil {
		return "", bizerr.ErrDB
	}

	return avatarURL, nil
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Nickname     string  `json:"nickname"`
	Bio          string  `json:"bio"`
	Company      string  `json:"company"`
	Occupation   string  `json:"occupation"`
	City         string  `json:"city"`
	Industry     string  `json:"industry"`
	LeftJobAt    *string `json:"left_job_at"`    // 格式：2006-01-02
	LastWorkDate *string `json:"last_work_date"` // 格式：2006-01-02
}

// GetUserPosts 获取用户帖子列表
func (s *UserService) GetUserPosts(userID int64, page, pageSize int) ([]map[string]interface{}, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	return s.repo.GetUserPosts(userID, page, pageSize)
}

// ============================================================
// 缓存辅助方法
// ============================================================

// invalidateUserCache 清除用户缓存
func (s *UserService) invalidateUserCache(userID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 清除用户详情缓存
	cacheKey := cache.FormatUserDetailKey(userID)
	cache.Delete(ctx, cacheKey)

	// 清除用户公开资料缓存
	cache.Delete(ctx, cacheKey+":public")

	// 清除用户统计缓存
	cache.Delete(ctx, cacheKey+":stats")
}

// GetUserStats 获取用户统计数据
func (s *UserService) GetUserStats(userID int64) (map[string]interface{}, error) {
	ctx := context.Background()
	cacheKey := cache.FormatUserDetailKey(userID) + ":stats"

	// 尝试从缓存获取
	var cachedStats map[string]interface{}
	if err := cache.Get(ctx, cacheKey, &cachedStats); err == nil {
		log.Printf("Cache hit for user stats: %d", userID)
		return cachedStats, nil
	}

	// 获取用户信息
	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerr.ErrUserNotFound
		}
		return nil, bizerr.ErrDB
	}

	// 获取帖子数量
	postCount, err := s.repo.CountUserPosts(userID)
	if err != nil {
		return nil, bizerr.ErrDB
	}

	// 获取获赞数量
	likeReceived, err := s.repo.CountUserLikesReceived(userID)
	if err != nil {
		return nil, bizerr.ErrDB
	}

	result := map[string]interface{}{
		"post_count":      postCount,
		"like_received":   likeReceived,
		"follower_count":  user.FollowerCount,
		"following_count": user.FollowingCount,
	}

	// 写入缓存
	if err := cache.Set(ctx, cacheKey, result, cache.CacheExpireUserStats); err != nil {
		log.Printf("Failed to cache user stats %d: %v", userID, err)
	}

	return result, nil
}
