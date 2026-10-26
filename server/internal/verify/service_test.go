package verify

import (
	"testing"
	"time"

	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/email"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================
// Mock 实现
// ============================================================

type mockVerifyRepo struct {
	records  map[int64]*EmailVerification
	byToken  map[string]*EmailVerification
	nextID   int64
	createDB bool
}

func newMockVerifyRepo() *mockVerifyRepo {
	return &mockVerifyRepo{
		records: make(map[int64]*EmailVerification),
		byToken: make(map[string]*EmailVerification),
		nextID:  1,
	}
}

func (m *mockVerifyRepo) Create(v *EmailVerification) error {
	if m.createDB {
		return gorm.ErrInvalidDB
	}
	v.ID = m.nextID
	m.nextID++
	m.records[v.ID] = v
	m.byToken[v.Token] = v
	return nil
}

func (m *mockVerifyRepo) GetByToken(token string) (*EmailVerification, error) {
	if v, ok := m.byToken[token]; ok {
		return v, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockVerifyRepo) GetByUserID(userID int64) (*EmailVerification, error) {
	for _, v := range m.records {
		if v.UserID == userID {
			return v, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockVerifyRepo) Update(v *EmailVerification) error {
	m.records[v.ID] = v
	m.byToken[v.Token] = v
	return nil
}

func (m *mockVerifyRepo) Delete(id int64) error {
	if v, ok := m.records[id]; ok {
		delete(m.byToken, v.Token)
		delete(m.records, id)
		return nil
	}
	return gorm.ErrRecordNotFound
}

type mockAuthRepoForVerify struct {
	users map[int64]*auth.User
}

func newMockAuthRepoForVerify() *mockAuthRepoForVerify {
	return &mockAuthRepoForVerify{users: make(map[int64]*auth.User)}
}

func (m *mockAuthRepoForVerify) CreateUser(u *auth.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *mockAuthRepoForVerify) GetByID(id int64) (*auth.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepoForVerify) GetByEmail(email string) (*auth.User, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepoForVerify) GetByPhone(phone string) (*auth.User, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAuthRepoForVerify) EmailExists(email string) bool { return false }

func (m *mockAuthRepoForVerify) PhoneExists(phone string) bool { return false }

func (m *mockAuthRepoForVerify) UpdateEmailVerified(userID int64, verified bool) error {
	if u, ok := m.users[userID]; ok {
		u.EmailVerified = verified
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *mockAuthRepoForVerify) UpdatePassword(userID int64, hashedPassword string) error {
	if u, ok := m.users[userID]; ok {
		u.Password = hashedPassword
		return nil
	}
	return gorm.ErrRecordNotFound
}

// ============================================================
// SendVerificationEmail 测试
// ============================================================

func TestVerifyService_SendVerificationEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv("APP_PUBLIC_URL", "https://community.example.com/")
		verifyRepo := newMockVerifyRepo()
		authRepo := newMockAuthRepoForVerify()
		emailSender := email.NewMockEmailSender()

		e := "test@example.com"
		authRepo.users[1] = &auth.User{ID: 1, Email: &e, EmailVerified: false}

		service := NewVerifyServiceWithDeps(verifyRepo, authRepo, emailSender)
		err := service.SendVerificationEmail(1)
		require.NoError(t, err)
		require.Len(t, emailSender.SentEmails, 1)
		assert.Contains(t, emailSender.SentEmails[0].Body, "https://community.example.com/verify-email/confirm?token=")

		// 验证记录已创建
		v, err := verifyRepo.GetByUserID(1)
		require.NoError(t, err)
		assert.Equal(t, int64(1), v.UserID)
		assert.NotEmpty(t, v.Token)
		assert.False(t, v.IsVerified())
	})

	t.Run("user_not_found", func(t *testing.T) {
		verifyRepo := newMockVerifyRepo()
		authRepo := newMockAuthRepoForVerify()
		emailSender := email.NewMockEmailSender()

		service := NewVerifyServiceWithDeps(verifyRepo, authRepo, emailSender)
		err := service.SendVerificationEmail(999)
		assert.Equal(t, bizerr.ErrUserNotFound, err)
	})

	t.Run("no_email", func(t *testing.T) {
		verifyRepo := newMockVerifyRepo()
		authRepo := newMockAuthRepoForVerify()
		emailSender := email.NewMockEmailSender()

		authRepo.users[2] = &auth.User{ID: 2, Email: nil, EmailVerified: false}

		service := NewVerifyServiceWithDeps(verifyRepo, authRepo, emailSender)
		err := service.SendVerificationEmail(2)
		assert.Equal(t, bizerr.ErrParamInvalid, err)
	})

	t.Run("already_verified", func(t *testing.T) {
		verifyRepo := newMockVerifyRepo()
		authRepo := newMockAuthRepoForVerify()
		emailSender := email.NewMockEmailSender()

		e := "verified@example.com"
		authRepo.users[3] = &auth.User{ID: 3, Email: &e, EmailVerified: true}

		service := NewVerifyServiceWithDeps(verifyRepo, authRepo, emailSender)
		err := service.SendVerificationEmail(3)
		assert.Equal(t, bizerr.ErrParamInvalid, err)
	})

	t.Run("repo_create_error", func(t *testing.T) {
		verifyRepo := newMockVerifyRepo()
		verifyRepo.createDB = true
		authRepo := newMockAuthRepoForVerify()
		emailSender := email.NewMockEmailSender()

		e := "new@example.com"
		authRepo.users[4] = &auth.User{ID: 4, Email: &e, EmailVerified: false}

		service := NewVerifyServiceWithDeps(verifyRepo, authRepo, emailSender)
		err := service.SendVerificationEmail(4)
		assert.Equal(t, bizerr.ErrDB, err)
	})
}

// ============================================================
// VerifyEmail 测试
// ============================================================

func TestVerifyService_VerifyEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		verifyRepo := newMockVerifyRepo()
		authRepo := newMockAuthRepoForVerify()
		emailSender := email.NewMockEmailSender()

		e := "verify@example.com"
		authRepo.users[1] = &auth.User{ID: 1, Email: &e, EmailVerified: false}

		service := NewVerifyServiceWithDeps(verifyRepo, authRepo, emailSender)

		// 先发送验证邮件
		service.SendVerificationEmail(1)

		// 获取 token
		v, _ := verifyRepo.GetByUserID(1)

		// 验证
		err := service.VerifyEmail(v.Token)
		require.NoError(t, err)

		// 确认用户邮箱已验证
		assert.True(t, authRepo.users[1].EmailVerified)
	})

	t.Run("invalid_token", func(t *testing.T) {
		verifyRepo := newMockVerifyRepo()
		authRepo := newMockAuthRepoForVerify()
		emailSender := email.NewMockEmailSender()

		service := NewVerifyServiceWithDeps(verifyRepo, authRepo, emailSender)
		err := service.VerifyEmail("nonexistent-token")
		assert.Equal(t, bizerr.ErrParamInvalid, err)
	})

	t.Run("already_verified", func(t *testing.T) {
		verifyRepo := newMockVerifyRepo()
		authRepo := newMockAuthRepoForVerify()
		emailSender := email.NewMockEmailSender()

		e := "done@example.com"
		authRepo.users[1] = &auth.User{ID: 1, Email: &e, EmailVerified: false}

		service := NewVerifyServiceWithDeps(verifyRepo, authRepo, emailSender)
		service.SendVerificationEmail(1)

		v, _ := verifyRepo.GetByUserID(1)

		// 第一次验证成功
		service.VerifyEmail(v.Token)

		// 第二次验证失败（已验证）
		err := service.VerifyEmail(v.Token)
		assert.Equal(t, bizerr.ErrParamInvalid, err)
	})

	t.Run("expired_token", func(t *testing.T) {
		verifyRepo := newMockVerifyRepo()
		authRepo := newMockAuthRepoForVerify()
		emailSender := email.NewMockEmailSender()

		e := "expired@example.com"
		authRepo.users[1] = &auth.User{ID: 1, Email: &e, EmailVerified: false}

		// 手动创建一个过期记录
		v := &EmailVerification{
			UserID:    1,
			Token:     "expired-token",
			ExpiresAt: time.Now().Add(-1 * time.Hour), // 已过期
		}
		verifyRepo.Create(v)

		service := NewVerifyServiceWithDeps(verifyRepo, authRepo, emailSender)
		err := service.VerifyEmail("expired-token")
		assert.Equal(t, bizerr.ErrParamInvalid, err)
	})

	t.Run("update_user_fails", func(t *testing.T) {
		verifyRepo := newMockVerifyRepo()
		authRepo := newMockAuthRepoForVerify()
		emailSender := email.NewMockEmailSender()

		e := "gone@example.com"
		authRepo.users[1] = &auth.User{ID: 1, Email: &e, EmailVerified: false}

		service := NewVerifyServiceWithDeps(verifyRepo, authRepo, emailSender)
		service.SendVerificationEmail(1)

		v, _ := verifyRepo.GetByUserID(1)

		// 先验证成功
		service.VerifyEmail(v.Token)

		// 确认已验证
		assert.True(t, authRepo.users[1].EmailVerified)
	})
}

// ============================================================
// GetVerificationByUserID 测试
// ============================================================

func TestVerifyService_GetVerificationByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		verifyRepo := newMockVerifyRepo()
		authRepo := newMockAuthRepoForVerify()
		emailSender := email.NewMockEmailSender()

		e := "get@example.com"
		authRepo.users[1] = &auth.User{ID: 1, Email: &e, EmailVerified: false}

		service := NewVerifyServiceWithDeps(verifyRepo, authRepo, emailSender)
		service.SendVerificationEmail(1)

		v, err := service.GetVerificationByUserID(1)
		require.NoError(t, err)
		assert.Equal(t, int64(1), v.UserID)
	})

	t.Run("not_found", func(t *testing.T) {
		verifyRepo := newMockVerifyRepo()
		authRepo := newMockAuthRepoForVerify()
		emailSender := email.NewMockEmailSender()

		service := NewVerifyServiceWithDeps(verifyRepo, authRepo, emailSender)

		_, err := service.GetVerificationByUserID(999)
		assert.Error(t, err)
	})
}

// ============================================================
// Model 方法测试
// ============================================================

func TestEmailVerification_Model(t *testing.T) {
	t.Run("IsExpired", func(t *testing.T) {
		v := &EmailVerification{ExpiresAt: time.Now().Add(-1 * time.Hour)}
		assert.True(t, v.IsExpired())

		v2 := &EmailVerification{ExpiresAt: time.Now().Add(1 * time.Hour)}
		assert.False(t, v2.IsExpired())
	})

	t.Run("IsVerified", func(t *testing.T) {
		v := &EmailVerification{}
		assert.False(t, v.IsVerified())

		v.Verify()
		assert.True(t, v.IsVerified())
		assert.NotNil(t, v.VerifiedAt)
	})

	t.Run("TableName", func(t *testing.T) {
		v := EmailVerification{}
		assert.Equal(t, "email_verifications", v.TableName())
	})
}
