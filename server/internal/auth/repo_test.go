package auth

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

// userColumns 用户表列名
var userColumns = []string{
	"id", "email", "phone", "password", "nickname",
	"bio", "company", "occupation", "avatar",
	"following_count", "follower_count", "is_admin", "status",
	"created_at", "updated_at",
}

// ============================================================
// CreateUser 测试
// ============================================================

func TestAuthRepo_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(`INSERT INTO users`).
			WillReturnRows(rows)

		email := "test@example.com"
		user := &User{
			Email:    &email,
			Password: "hashed_password",
			Nickname: "testuser",
		}

		err := repo.CreateUser(user)
		require.NoError(t, err)
		assert.Equal(t, int64(1), user.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		mock.ExpectQuery(`INSERT INTO users`).
			WillReturnError(sql.ErrConnDone)

		email := "test@example.com"
		user := &User{
			Email:    &email,
			Password: "hashed_password",
			Nickname: "testuser",
		}

		err := repo.CreateUser(user)
		require.Error(t, err)
	})
}

// ============================================================
// GetByID 测试
// ============================================================

func TestAuthRepo_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		now := time.Now()
		rows := sqlmock.NewRows(userColumns).AddRow(
			1, "test@example.com", nil, "hashed_pwd", "testuser",
			"", "", "", "",
			0, 0, false, int8(1),
			now, now,
		)
		mock.ExpectQuery(`SELECT \* FROM users WHERE id = .+ ORDER BY`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		user, err := repo.GetByID(1)
		require.NoError(t, err)
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "testuser", user.Nickname)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		mock.ExpectQuery(`SELECT \* FROM users WHERE id = .+ ORDER BY`).
			WithArgs(int64(999)).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.GetByID(999)
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

// ============================================================
// GetByEmail 测试
// ============================================================

func TestAuthRepo_GetByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		now := time.Now()
		rows := sqlmock.NewRows(userColumns).AddRow(
			1, "test@example.com", nil, "hashed_pwd", "testuser",
			"", "", "", "",
			0, 0, false, int8(1),
			now, now,
		)
		mock.ExpectQuery(`SELECT \* FROM users WHERE email = .+ ORDER BY`).
			WithArgs("test@example.com").
			WillReturnRows(rows)

		user, err := repo.GetByEmail("test@example.com")
		require.NoError(t, err)
		assert.Equal(t, "test@example.com", *user.Email)
	})

	t.Run("not_found", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		mock.ExpectQuery(`SELECT \* FROM users WHERE email = .+ ORDER BY`).
			WithArgs("noone@example.com").
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.GetByEmail("noone@example.com")
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

// ============================================================
// GetByPhone 测试
// ============================================================

func TestAuthRepo_GetByPhone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		now := time.Now()
		phone := "13800000000"
		rows := sqlmock.NewRows(userColumns).AddRow(
			1, nil, &phone, "hashed_pwd", "phoneuser",
			"", "", "", "",
			0, 0, false, int8(1),
			now, now,
		)
		mock.ExpectQuery(`SELECT \* FROM users WHERE phone = .+ ORDER BY`).
			WithArgs("13800000000").
			WillReturnRows(rows)

		user, err := repo.GetByPhone("13800000000")
		require.NoError(t, err)
		assert.Equal(t, "13800000000", *user.Phone)
	})

	t.Run("not_found", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		mock.ExpectQuery(`SELECT \* FROM users WHERE phone = .+ ORDER BY`).
			WithArgs("13900000000").
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.GetByPhone("13900000000")
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

// ============================================================
// EmailExists 测试
// ============================================================

func TestAuthRepo_EmailExists(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(`SELECT count\(\*\) FROM users WHERE email =`).
			WithArgs("test@example.com").
			WillReturnRows(rows)

		exists := repo.EmailExists("test@example.com")
		assert.True(t, exists)
	})

	t.Run("not_exists", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(`SELECT count\(\*\) FROM users WHERE email =`).
			WithArgs("no@example.com").
			WillReturnRows(rows)

		exists := repo.EmailExists("no@example.com")
		assert.False(t, exists)
	})
}

// ============================================================
// PhoneExists 测试
// ============================================================

func TestAuthRepo_PhoneExists(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(`SELECT count\(\*\) FROM users WHERE phone =`).
			WithArgs("13800000000").
			WillReturnRows(rows)

		exists := repo.PhoneExists("13800000000")
		assert.True(t, exists)
	})

	t.Run("not_exists", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(`SELECT count\(\*\) FROM users WHERE phone =`).
			WithArgs("13900000000").
			WillReturnRows(rows)

		exists := repo.PhoneExists("13900000000")
		assert.False(t, exists)
	})
}

// ============================================================
// NicknameExists 测试
// ============================================================

func TestAuthRepo_NicknameExists(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(`SELECT count\(\*\) FROM users WHERE nickname =`).
			WithArgs("testuser").
			WillReturnRows(rows)

		exists := repo.NicknameExists("testuser")
		assert.True(t, exists)
	})

	t.Run("not_exists", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(`SELECT count\(\*\) FROM users WHERE nickname =`).
			WithArgs("newuser").
			WillReturnRows(rows)

		exists := repo.NicknameExists("newuser")
		assert.False(t, exists)
	})
}

// ============================================================
// GetAll 测试
// ============================================================

func TestAuthRepo_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		now := time.Now()
		rows := sqlmock.NewRows(userColumns).
			AddRow(1, "a@test.com", nil, "pwd", "user1", "", "", "", "", 0, 0, false, int8(1), now, now).
			AddRow(2, "b@test.com", nil, "pwd", "user2", "", "", "", "", 0, 0, false, int8(1), now, now)
		mock.ExpectQuery(`SELECT \* FROM users ORDER BY`).
			WillReturnRows(rows)

		users, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "user1", users[0].Nickname)
		assert.Equal(t, "user2", users[1].Nickname)
	})

	t.Run("empty", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		rows := sqlmock.NewRows(userColumns)
		mock.ExpectQuery(`SELECT \* FROM users ORDER BY`).
			WillReturnRows(rows)

		users, err := repo.GetAll()
		require.NoError(t, err)
		assert.Empty(t, users)
	})
}

// ============================================================
// Count 测试
// ============================================================

func TestAuthRepo_Count(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
		mock.ExpectQuery(`SELECT count\(\*\) FROM users`).
			WillReturnRows(rows)

		count, err := repo.Count()
		require.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("empty", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewAuthRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(`SELECT count\(\*\) FROM users`).
			WillReturnRows(rows)

		count, err := repo.Count()
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}
