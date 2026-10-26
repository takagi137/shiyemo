package verify

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupMockDB 创建 gorm.DB + sqlmock
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 db,
		WithoutQuotingCheck:  true,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	return gormDB, mock
}

// verificationColumns 验证表列名
var verificationColumns = []string{
	"id", "user_id", "token", "expires_at", "verified_at", "created_at",
}

// ============================================================
// Create 测试
// ============================================================

func TestVerifyRepo_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewVerifyRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(`INSERT INTO email_verifications`).
			WillReturnRows(rows)

		verification := &EmailVerification{
			UserID:    1,
			Token:     "test-token-123",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		err := repo.Create(verification)
		require.NoError(t, err)
		assert.Equal(t, int64(1), verification.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewVerifyRepoWithDB(gormDB)

		mock.ExpectQuery(`INSERT INTO email_verifications`).
			WillReturnError(sql.ErrConnDone)

		verification := &EmailVerification{
			UserID:    1,
			Token:     "test-token-123",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		err := repo.Create(verification)
		require.Error(t, err)
	})
}

// ============================================================
// GetByToken 测试
// ============================================================

func TestVerifyRepo_GetByToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewVerifyRepoWithDB(gormDB)

		now := time.Now()
		expiresAt := now.Add(24 * time.Hour)
		rows := sqlmock.NewRows(verificationColumns).AddRow(
			1, 1, "test-token-123", expiresAt, nil, now,
		)
		mock.ExpectQuery(`SELECT \* FROM email_verifications WHERE token = .+ ORDER BY`).
			WithArgs("test-token-123").
			WillReturnRows(rows)

		verification, err := repo.GetByToken("test-token-123")
		require.NoError(t, err)
		assert.Equal(t, int64(1), verification.ID)
		assert.Equal(t, "test-token-123", verification.Token)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewVerifyRepoWithDB(gormDB)

		mock.ExpectQuery(`SELECT \* FROM email_verifications WHERE token = .+ ORDER BY`).
			WithArgs("nonexistent-token").
			WillReturnError(gorm.ErrRecordNotFound)

		verification, err := repo.GetByToken("nonexistent-token")
		require.Error(t, err)
		assert.Nil(t, verification)
	})
}

// ============================================================
// GetByUserID 测试
// ============================================================

func TestVerifyRepo_GetByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewVerifyRepoWithDB(gormDB)

		now := time.Now()
		expiresAt := now.Add(24 * time.Hour)
		rows := sqlmock.NewRows(verificationColumns).AddRow(
			1, 1, "test-token-123", expiresAt, nil, now,
		)
		mock.ExpectQuery(`SELECT \* FROM email_verifications WHERE user_id = .+ ORDER BY`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		verification, err := repo.GetByUserID(1)
		require.NoError(t, err)
		assert.Equal(t, int64(1), verification.UserID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewVerifyRepoWithDB(gormDB)

		mock.ExpectQuery(`SELECT \* FROM email_verifications WHERE user_id = .+ ORDER BY`).
			WithArgs(int64(999)).
			WillReturnError(gorm.ErrRecordNotFound)

		verification, err := repo.GetByUserID(999)
		require.Error(t, err)
		assert.Nil(t, verification)
	})
}

// ============================================================
// Update 测试
// ============================================================

func TestVerifyRepo_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewVerifyRepoWithDB(gormDB)

		mock.ExpectExec(`UPDATE email_verifications`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		verification := &EmailVerification{
			ID:        1,
			UserID:    1,
			Token:     "test-token-123",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		err := repo.Update(verification)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewVerifyRepoWithDB(gormDB)

		mock.ExpectExec(`UPDATE email_verifications`).
			WillReturnError(sql.ErrConnDone)

		verification := &EmailVerification{
			ID:        1,
			UserID:    1,
			Token:     "test-token-123",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		err := repo.Update(verification)
		require.Error(t, err)
	})
}

// ============================================================
// Delete 测试
// ============================================================

func TestVerifyRepo_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewVerifyRepoWithDB(gormDB)

		mock.ExpectExec(`DELETE FROM email_verifications`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(1)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewVerifyRepoWithDB(gormDB)

		mock.ExpectExec(`DELETE FROM email_verifications`).
			WillReturnError(sql.ErrConnDone)

		err := repo.Delete(1)
		require.Error(t, err)
	})
}
