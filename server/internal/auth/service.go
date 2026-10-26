package auth

import (
	"errors"
	"time"

	"github.com/shiyemo/server/internal/common/config"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	jwtutil "github.com/shiyemo/server/internal/common/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyServiceInterface 验证服务接口
type VerifyServiceInterface interface {
	SendVerificationEmail(userID int64) error
}

// AuthService 认证服务
type AuthService struct {
	repo          AuthRepository
	verifyService VerifyServiceInterface
}

// NewAuthService 创建认证服务（使用数据库）
func NewAuthService() *AuthService {
	return &AuthService{
		repo: NewAuthRepo(),
	}
}

// NewAuthServiceWithDeps 创建认证服务（注入依赖，用于测试）
func NewAuthServiceWithDeps(repo AuthRepository, verifyService VerifyServiceInterface) *AuthService {
	return &AuthService{
		repo:          repo,
		verifyService: verifyService,
	}
}

// RegisterByEmail 邮箱注册
func (s *AuthService) RegisterByEmail(params RegisterParams) (string, *User, error) {
	// 检查邮箱是否已注册
	if s.repo.EmailExists(params.Email) {
		return "", nil, bizerr.ErrEmailExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, bizerr.ErrInternal
	}

	// 创建用户
	user := &User{
		Email:     &params.Email,
		Password:  string(hashedPassword),
		Nickname:  params.Nickname,
		City:      params.City,
		Industry:  params.Industry,
	}

	// 解析离职日期
	if params.LeftJobAt != "" {
		if t, err := time.Parse("2006-01-02", params.LeftJobAt); err == nil {
			user.LeftJobAt = &t
		}
	}

	if err := s.repo.CreateUser(user); err != nil {
		return "", nil, bizerr.ErrDB
	}

	// 生成 token
	token, err := jwtutil.GenerateToken(user.ID)
	if err != nil {
		return "", nil, bizerr.ErrInternal
	}

	// 注册成功后触发发送验证邮件
	if s.verifyService != nil {
		s.verifyService.SendVerificationEmail(user.ID)
	}

	return token, user, nil
}

// RegisterByPhone 手机号注册
func (s *AuthService) RegisterByPhone(phone, code, nickname string) (string, *User, error) {
	// 验证验证码
	if err := s.verifyCode(phone, code); err != nil {
		return "", nil, err
	}

	// 检查手机号是否已注册
	if s.repo.PhoneExists(phone) {
		return "", nil, bizerr.ErrPhoneExists
	}

	// 创建用户（手机号注册不设密码）
	user := &User{
		Phone:    &phone,
		Password: "",
		Nickname: nickname,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return "", nil, bizerr.ErrDB
	}

	// 生成 token
	token, err := jwtutil.GenerateToken(user.ID)
	if err != nil {
		return "", nil, bizerr.ErrInternal
	}

	return token, user, nil
}

// LoginByEmail 邮箱登录
func (s *AuthService) LoginByEmail(email, password string) (string, *User, error) {
	// 查询用户
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, bizerr.ErrUserNotFound
		}
		return "", nil, bizerr.ErrDB
	}

	// 检查账号状态
	if !user.IsActive() {
		return "", nil, bizerr.ErrUserDisabled
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, bizerr.ErrPasswordWrong
	}

	// 生成 token
	token, err := jwtutil.GenerateToken(user.ID)
	if err != nil {
		return "", nil, bizerr.ErrInternal
	}

	return token, user, nil
}

// LoginByPhone 手机号登录
func (s *AuthService) LoginByPhone(phone, code string) (string, *User, error) {
	// 验证验证码
	if err := s.verifyCode(phone, code); err != nil {
		return "", nil, err
	}

	// 查询用户
	user, err := s.repo.GetByPhone(phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, bizerr.ErrUserNotFound
		}
		return "", nil, bizerr.ErrDB
	}

	// 检查账号状态
	if !user.IsActive() {
		return "", nil, bizerr.ErrUserDisabled
	}

	// 生成 token
	token, err := jwtutil.GenerateToken(user.ID)
	if err != nil {
		return "", nil, bizerr.ErrInternal
	}

	return token, user, nil
}

// SendSMS 发送验证码（MVP 阶段 mock）
func (s *AuthService) SendSMS(phone string) error {
	return nil
}

// GetProfile 获取当前用户信息
func (s *AuthService) GetProfile(userID int64) (*User, error) {
	return s.repo.GetByID(userID)
}

// verifyCode 验证验证码（MVP 阶段 mock 固定 123456）
func (s *AuthService) verifyCode(phone, code string) error {
	cfg := config.Load()

	if cfg.App.Mode == "debug" && code == "123456" {
		return nil
	}

	return bizerr.ErrCodeWrong
}
