package interaction

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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
// Like / Unlike / IsLiked
// ============================================================

func TestInteractionRepo_Like(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewInteractionRepoWithDB(db)

	mock.ExpectQuery(`INSERT INTO likes`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	require.NoError(t, repo.Like(1, 1))
}

func TestInteractionRepo_Unlike(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewInteractionRepoWithDB(db)

	mock.ExpectExec(`DELETE FROM likes WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.Unlike(1, 1))
}

func TestInteractionRepo_IsLiked(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewInteractionRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM likes WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	assert.True(t, repo.IsLiked(1, 1))
}

func TestInteractionRepo_IsLiked_False(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewInteractionRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM likes WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	assert.False(t, repo.IsLiked(1, 1))
}

// ============================================================
// Bookmark / Unbookmark / IsBookmarked
// ============================================================

func TestInteractionRepo_Bookmark(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewInteractionRepoWithDB(db)

	mock.ExpectQuery(`INSERT INTO bookmarks`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	require.NoError(t, repo.Bookmark(1, 1))
}

func TestInteractionRepo_Unbookmark(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewInteractionRepoWithDB(db)

	mock.ExpectExec(`DELETE FROM bookmarks WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.Unbookmark(1, 1))
}

func TestInteractionRepo_IsBookmarked(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewInteractionRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM bookmarks WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	assert.True(t, repo.IsBookmarked(1, 1))
}

// ============================================================
// CountLikes
// ============================================================

func TestInteractionRepo_CountLikes(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewInteractionRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM likes WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	count, err := repo.CountLikes(1)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestInteractionRepo_CountLikes_Error(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewInteractionRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM likes WHERE .+`).WillReturnError(sql.ErrConnDone)

	_, err := repo.CountLikes(1)
	require.Error(t, err)
}

// ============================================================
// ListBookmarks
// ============================================================

func TestInteractionRepo_ListBookmarks(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewInteractionRepoWithDB(db)

		mock.ExpectQuery(`SELECT count\(\*\) FROM bookmarks WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		mock.ExpectQuery(`SELECT .+ FROM bookmarks WHERE .+ ORDER BY`).WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "post_id", "created_at"}).
				AddRow(1, 1, 10, nil).
				AddRow(2, 1, 20, nil))

		ids, total, err := repo.ListBookmarks(1, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, []int64{10, 20}, ids)
	})

	t.Run("empty", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewInteractionRepoWithDB(db)

		mock.ExpectQuery(`SELECT count\(\*\) FROM bookmarks WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(`SELECT .+ FROM bookmarks WHERE .+ ORDER BY`).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "created_at"}))

		ids, total, err := repo.ListBookmarks(1, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, ids)
	})
}
