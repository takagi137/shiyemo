package user

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shiyemo/server/internal/auth"
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

// ============================================================
// GetByID 测试
// ============================================================

func TestUserRepo_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewUserRepoWithDB(gormDB)

		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"id", "email", "phone", "password", "nickname",
			"bio", "company", "occupation", "avatar",
			"following_count", "follower_count", "is_admin", "status",
			"created_at", "updated_at",
		}).AddRow(
			1, "test@example.com", nil, "hashed_pwd", "testuser",
			"bio text", "company", "engineer", "avatar.png",
			10, 5, false, int8(1),
			now, now,
		)
		mock.ExpectQuery(`SELECT \* FROM users WHERE id = .+ ORDER BY`).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		user, err := repo.GetByID(1)
		require.NoError(t, err)
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "testuser", user.Nickname)
		assert.Equal(t, "bio text", user.Bio)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewUserRepoWithDB(gormDB)

		mock.ExpectQuery(`SELECT \* FROM users WHERE id = .+ ORDER BY`).
			WithArgs(int64(999)).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.GetByID(999)
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

// ============================================================
// Update 测试
// ============================================================

func TestUserRepo_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewUserRepoWithDB(gormDB)

		mock.ExpectExec(`UPDATE users SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		user := &auth.User{
			ID:       1,
			Nickname: "updated_name",
			Bio:      "new bio",
		}

		err := repo.Update(user)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewUserRepoWithDB(gormDB)

		mock.ExpectExec(`UPDATE users SET`).
			WillReturnError(sql.ErrConnDone)

		user := &auth.User{
			ID:       1,
			Nickname: "updated_name",
		}

		err := repo.Update(user)
		require.Error(t, err)
	})
}

// ============================================================
// UpdateFields 测试
// ============================================================

func TestUserRepo_UpdateFields(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewUserRepoWithDB(gormDB)

		mock.ExpectExec(`UPDATE users SET`).
			WithArgs("new avatar", sqlmock.AnyArg(), int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateFields(1, map[string]interface{}{
			"avatar": "new avatar",
		})
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewUserRepoWithDB(gormDB)

		mock.ExpectExec(`UPDATE users SET`).
			WithArgs("new avatar", sqlmock.AnyArg(), int64(1)).
			WillReturnError(sql.ErrConnDone)

		err := repo.UpdateFields(1, map[string]interface{}{
			"avatar": "new avatar",
		})
		require.Error(t, err)
	})
}

// ============================================================
// GetUserPosts 测试
// ============================================================

func TestUserRepo_GetUserPosts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewUserRepoWithDB(gormDB)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(`SELECT count\(\*\) FROM posts WHERE`).
			WillReturnRows(countRows)

		postRows := sqlmock.NewRows([]string{
			"id", "user_id", "type", "title", "content",
			"images", "tags", "like_count", "comment_count", "created_at",
		}).
			AddRow(1, 1, "emotion", "post1", "content1", "", "", 5, 2, time.Now()).
			AddRow(2, 1, "insight", "post2", "content2", "", "", 3, 1, time.Now())
		mock.ExpectQuery(`SELECT .+ FROM posts WHERE`).
			WillReturnRows(postRows)

		posts, total, err := repo.GetUserPosts(1, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, posts, 2)
	})

	t.Run("empty", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewUserRepoWithDB(gormDB)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(`SELECT count\(\*\) FROM posts WHERE`).
			WillReturnRows(countRows)

		postRows := sqlmock.NewRows([]string{
			"id", "user_id", "type", "title", "content",
			"images", "tags", "like_count", "comment_count", "created_at",
		})
		mock.ExpectQuery(`SELECT .+ FROM posts WHERE`).
			WillReturnRows(postRows)

		posts, total, err := repo.GetUserPosts(999, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, posts)
	})
}

// ============================================================
// CountUserPosts 测试
// ============================================================

func TestUserRepo_CountUserPosts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewUserRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
		mock.ExpectQuery(`SELECT count\(\*\) FROM posts WHERE`).
			WillReturnRows(rows)

		count, err := repo.CountUserPosts(1)
		require.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}

// ============================================================
// CountUserLikesReceived 测试
// ============================================================

func TestUserRepo_CountUserLikesReceived(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewUserRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"coalesce"}).AddRow(42)
		mock.ExpectQuery(`SELECT COALESCE\(SUM\(like_count\), 0\) FROM posts WHERE`).
			WillReturnRows(rows)

		count, err := repo.CountUserLikesReceived(1)
		require.NoError(t, err)
		assert.Equal(t, int64(42), count)
	})

	t.Run("zero", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewUserRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"coalesce"}).AddRow(0)
		mock.ExpectQuery(`SELECT COALESCE\(SUM\(like_count\), 0\) FROM posts WHERE`).
			WillReturnRows(rows)

		count, err := repo.CountUserLikesReceived(999)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}
