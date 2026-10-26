package verify

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/cache"
	"github.com/shiyemo/server/internal/common/config"
	"github.com/shiyemo/server/internal/common/email"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/setting"
	"gorm.io/gorm"
	"strings"
)

// VerifyService 验证服务
type VerifyService struct {
	repo     VerifyRepoInterface
	authRepo auth.AuthRepository
	email    email.EmailSender
}

// NewVerifyService 创建验证服务。
// 邮件发送器使用配置提供者，每次发信时从 system_settings 表重新读取 SMTP 配置，
// 因此管理后台保存设置后立即生效，无需重启服务。
func NewVerifyService() *VerifyService {
	return &VerifyService{
		repo:     NewVerifyRepo(),
		authRepo: auth.NewAuthRepo(),
		email:    email.NewSMTPEmailSenderWithProvider(setting.NewSettingService().SMTPConfigProvider()),
	}
}

// NewVerifyServiceWithDeps 创建验证服务（注入依赖，用于测试）
func NewVerifyServiceWithDeps(repo VerifyRepoInterface, authRepo auth.AuthRepository, emailSender email.EmailSender) *VerifyService {
	return &VerifyService{
		repo:     repo,
		authRepo: authRepo,
		email:    emailSender,
	}
}

// SendVerificationEmail 发送验证邮件
func (s *VerifyService) SendVerificationEmail(userID int64) error {
	// 获取用户信息
	user, err := s.authRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerr.ErrUserNotFound
		}
		return bizerr.ErrDB
	}

	// 检查用户是否有邮箱
	if !user.HasEmail() {
		return bizerr.ErrParamInvalid
	}

	// 检查邮箱是否已验证
	if user.IsEmailVerified() {
		return bizerr.ErrParamInvalid
	}

	// 生成验证 token
	token, err := generateToken()
	if err != nil {
		return bizerr.ErrInternal
	}

	// 创建验证记录
	verification := &EmailVerification{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24小时过期
	}

	if err := s.repo.Create(verification); err != nil {
		return bizerr.ErrDB
	}

	// 发送验证邮件
	baseURL := strings.TrimRight(config.Load().App.PublicURL, "/")
	if err := s.email.SendVerificationEmail(*user.Email, token, baseURL); err != nil {
		// 邮件发送失败不影响注册，只记录日志
		// log.Printf("发送验证邮件失败: %v", err)
	}

	return nil
}

// VerifyEmail 验证邮箱
func (s *VerifyService) VerifyEmail(token string) error {
	// 查询验证记录
	verification, err := s.repo.GetByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerr.ErrParamInvalid
		}
		return bizerr.ErrDB
	}

	// 检查是否已验证
	if verification.IsVerified() {
		return bizerr.ErrParamInvalid
	}

	// 检查是否过期
	if verification.IsExpired() {
		return bizerr.ErrParamInvalid
	}

	// 标记为已验证
	verification.Verify()
	if err := s.repo.Update(verification); err != nil {
		return bizerr.ErrDB
	}

	// 更新用户邮箱验证状态
	user, err := s.authRepo.GetByID(verification.UserID)
	if err != nil {
		return bizerr.ErrDB
	}

	// 使用 repo 更新用户字段
	if err := s.updateUserEmailVerified(user.ID, true); err != nil {
		return err
	}

	return nil
}

// GetVerificationByUserID 获取用户的验证记录
func (s *VerifyService) GetVerificationByUserID(userID int64) (*EmailVerification, error) {
	return s.repo.GetByUserID(userID)
}

// updateUserEmailVerified 更新用户邮箱验证状态
func (s *VerifyService) updateUserEmailVerified(userID int64, verified bool) error {
	return s.authRepo.UpdateEmailVerified(userID, verified)
}

// generateToken 生成随机 token
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SendResetEmail 发送密码重置邮件
func (s *VerifyService) SendResetEmail(emailAddr string) error {
	// 查找用户
	user, err := s.authRepo.GetByEmail(emailAddr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在也返回成功，防止邮箱枚举
			return nil
		}
		return bizerr.ErrDB
	}

	// 生成 token
	token, err := generateToken()
	if err != nil {
		return bizerr.ErrInternal
	}

	// 存入 Redis，15 分钟过期
	ctx := context.Background()
	cacheKey := "reset_password:" + token
	if err := cache.Set(ctx, cacheKey, user.ID, 15*time.Minute); err != nil {
		return bizerr.ErrInternal
	}

	// 发送邮件
	baseURL := strings.TrimRight(config.Load().App.PublicURL, "/")
	if err := s.email.SendResetPasswordEmail(*user.Email, token, baseURL); err != nil {
		// 邮件发送失败只记录日志
		return nil
	}

	return nil
}

// ResetPassword 重置密码
func (s *VerifyService) ResetPassword(token, newPassword string) error {
	// 从 Redis 获取用户 ID
	ctx := context.Background()
	cacheKey := "reset_password:" + token
	var userID int64
	if err := cache.Get(ctx, cacheKey, &userID); err != nil {
		return bizerr.ErrParamInvalid
	}

	// 加密密码
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return bizerr.ErrInternal
	}

	// 更新密码
	if err := s.authRepo.UpdatePassword(userID, hashedPassword); err != nil {
		return bizerr.ErrDB
	}

	// 删除 token
	cache.Delete(ctx, cacheKey)

	return nil
}
