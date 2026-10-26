package auth

import (
	"testing"

	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ============================================================
// Mock 实现
// ============================================================

type mockAuthRepo struct {
	users    map[int64]*User
	nextID   int64
	createDB bool // simulate DB error
}

func newMockAuthRepo() *mockAuthRepo {
	return &mockAuthRepo{users: make(map[int64]*User), nextID: 1}
}

type mockVerifyService struct {
	sentEmails []int64
}

func newMockVerifyService() *mockVerifyService {
	return &mockVerifyService{sentEmails: make([]int64, 0)}
}

func (m *mockVerifyService) SendVerificationEmail(userID int64) error {
	m.sentEmails = append(m.sentEmails, userID)
	return nil
}

func (m *mockAuthRepo) CreateUser(user *User) error {
	if m.createDB {
		return gorm.ErrInvalidDB
	}
	user.ID = m.nextID
	m.nextID++
	if user.Status == 0 {
		user.Status = 1
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockAuthRepo) GetByID(id int64) (*User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) GetByEmail(email string) (*User, error) {
	for _, u := range m.users {
		if u.Email != nil && *u.Email == email {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) GetByPhone(phone string) (*User, error) {
	for _, u := range m.users {
		if u.Phone != nil && *u.Phone == phone {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) EmailExists(email string) bool {
	for _, u := range m.users {
		if u.Email != nil && *u.Email == email {
			return true
		}
	}
	return false
}

func (m *mockAuthRepo) PhoneExists(phone string) bool {
	for _, u := range m.users {
		if u.Phone != nil && *u.Phone == phone {
			return true
		}
	}
	return false
}

func (m *mockAuthRepo) DisableUser(phone string) {
	for _, u := range m.users {
		if u.Phone != nil && *u.Phone == phone {
			u.Status = 0
			return
		}
	}
}

func (m *mockAuthRepo) UpdateEmailVerified(userID int64, verified bool) error {
	if u, ok := m.users[userID]; ok {
		u.EmailVerified = verified
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *mockAuthRepo) UpdatePassword(userID int64, hashedPassword string) error {
	if u, ok := m.users[userID]; ok {
		u.Password = hashedPassword
		return nil
	}
	return gorm.ErrRecordNotFound
}

// ============================================================
// RegisterByEmail 测试
// ============================================================

func TestAuthService_RegisterByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockAuthRepo()
		verifyService := newMockVerifyService()
		service := NewAuthServiceWithDeps(repo, verifyService)

		token, user, err := service.RegisterByEmail(RegisterParams{Email: "test@example.com", Password: "a123456", Nickname: "测试用户"})

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", *user.Email)
		assert.Equal(t, "测试用户", user.Nickname)
		assert.NotEmpty(t, user.Password) // 密码已加密
	})

	t.Run("success_with_optional_fields", func(t *testing.T) {
		repo := newMockAuthRepo()
		verifyService := newMockVerifyService()
		service := NewAuthServiceWithDeps(repo, verifyService)

		token, user, err := service.RegisterByEmail(RegisterParams{
			Email:     "with_fields@example.com",
			Password:  "a123456",
			Nickname:  "有字段用户",
			City:      "北京",
			Industry:  "互联网",
			LeftJobAt: "2024-01-15",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotNil(t, user)
		assert.Equal(t, "北京", user.City)
		assert.Equal(t, "互联网", user.Industry)
		assert.NotNil(t, user.LeftJobAt)
	})

	t.Run("duplicate_email", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		service.RegisterByEmail(RegisterParams{Email: "dup@example.com", Password: "a123456", Nickname: "用户1"})
		_, _, err := service.RegisterByEmail(RegisterParams{Email: "dup@example.com", Password: "a123456", Nickname: "用户2"})

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrEmailExists, err)
	})

	t.Run("db_error", func(t *testing.T) {
		repo := newMockAuthRepo()
		repo.createDB = true
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		_, _, err := service.RegisterByEmail(RegisterParams{Email: "new@example.com", Password: "a123456", Nickname: "用户"})

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrDB, err)
	})
}

// ============================================================
// LoginByEmail 测试
// ============================================================

func TestAuthService_LoginByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		// 先注册
		service.RegisterByEmail(RegisterParams{Email: "login@example.com", Password: "a123456", Nickname: "登录用户"})

		// 登录
		token, user, err := service.LoginByEmail("login@example.com", "a123456")

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, "login@example.com", *user.Email)
	})

	t.Run("wrong_password", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		service.RegisterByEmail(RegisterParams{Email: "wp@example.com", Password: "a123456", Nickname: "用户"})
		_, _, err := service.LoginByEmail("wp@example.com", "wrongpass")

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrPasswordWrong, err)
	})

	t.Run("user_not_found", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		_, _, err := service.LoginByEmail("noexist@example.com", "a123456")

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrUserNotFound, err)
	})

	t.Run("user_disabled", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		// 注册后手动禁用
		service.RegisterByEmail(RegisterParams{Email: "disabled@example.com", Password: "a123456", Nickname: "用户"})
		repo.users[1].Status = 0

		_, _, err := service.LoginByEmail("disabled@example.com", "a123456")

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrUserDisabled, err)
	})
}

// ============================================================
// RegisterByPhone 测试
// ============================================================

func TestAuthService_RegisterByPhone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		token, user, err := service.RegisterByPhone("13800138000", "123456", "手机用户")

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, "13800138000", *user.Phone)
		assert.Equal(t, "手机用户", user.Nickname)
	})

	t.Run("duplicate_phone", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		service.RegisterByPhone("13800138001", "123456", "用户1")
		_, _, err := service.RegisterByPhone("13800138001", "123456", "用户2")

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrPhoneExists, err)
	})

	t.Run("invalid_code", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		_, _, err := service.RegisterByPhone("13800138002", "000000", "用户")

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrCodeWrong, err)
	})
}

// ============================================================
// LoginByPhone 测试
// ============================================================

func TestAuthService_LoginByPhone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		service.RegisterByPhone("13800138003", "123456", "手机登录")
		token, user, err := service.LoginByPhone("13800138003", "123456")

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, "13800138003", *user.Phone)
	})

	t.Run("wrong_code", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		service.RegisterByPhone("13800138004", "123456", "用户")
		_, _, err := service.LoginByPhone("13800138004", "000000")

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrCodeWrong, err)
	})

	t.Run("user_not_found", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		_, _, err := service.LoginByPhone("13900139000", "123456")

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrUserNotFound, err)
	})

	t.Run("user_disabled", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		// 注册用户
		service.RegisterByPhone("13800138005", "123456", "禁用用户")

		// 禁用用户
		repo.DisableUser("13800138005")

		_, _, err := service.LoginByPhone("13800138005", "123456")

		require.Error(t, err)
		assert.Equal(t, bizerr.ErrUserDisabled, err)
	})
}

// ============================================================
// SendSMS 测试
// ============================================================

func TestAuthService_SendSMS(t *testing.T) {
	repo := newMockAuthRepo()
	service := NewAuthServiceWithDeps(repo, newMockVerifyService())
	err := service.SendSMS("13800138000")
	require.NoError(t, err)
}

// ============================================================
// GetProfile 测试
// =============================================================

func TestAuthService_GetProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		_, user, _ := service.RegisterByEmail(RegisterParams{Email: "profile@example.com", Password: "a123456", Nickname: "资料用户"})

		profile, err := service.GetProfile(user.ID)
		require.NoError(t, err)
		assert.Equal(t, "profile@example.com", *profile.Email)
	})

	t.Run("not_found", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		_, err := service.GetProfile(999)
		require.Error(t, err)
	})
}

// ============================================================
// UpdateEmailVerified 测试（mock 仓储层）
// =============================================================

func TestAuthRepo_UpdateEmailVerified(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := newMockAuthRepo()
		service := NewAuthServiceWithDeps(repo, newMockVerifyService())

		_, user, _ := service.RegisterByEmail(RegisterParams{Email: "verify@example.com", Password: "a123456", Nickname: "验证用户"})

		err := repo.UpdateEmailVerified(user.ID, true)
		require.NoError(t, err)
		assert.True(t, repo.users[user.ID].EmailVerified)
	})

	t.Run("user_not_found", func(t *testing.T) {
		repo := newMockAuthRepo()

		err := repo.UpdateEmailVerified(999, true)
		require.Error(t, err)
	})
}

// ============================================================
// 辅助：验证密码哈希
// ============================================================

func TestPasswordHashing(t *testing.T) {
	t.Run("hash_and_verify", func(t *testing.T) {
		password := "a123456"
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		// 验证正确密码
		err = bcrypt.CompareHashAndPassword(hash, []byte(password))
		assert.NoError(t, err)

		// 验证错误密码
		err = bcrypt.CompareHashAndPassword(hash, []byte("wrong"))
		assert.Error(t, err)
	})
}
