package comment

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

var commentColumns = []string{
	"id", "post_id", "user_id", "parent_id", "content",
	"like_count", "status", "created_at", "updated_at", "deleted_at",
}

// ============================================================
// Create
// ============================================================

func TestCommentRepo_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewCommentRepoWithDB(db)

		mock.ExpectQuery(`INSERT INTO comments`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		c := &Comment{PostID: 1, UserID: 1, Content: "great!"}
		require.NoError(t, repo.Create(c))
		assert.Equal(t, int64(1), c.ID)
	})

	t.Run("db_error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewCommentRepoWithDB(db)

		mock.ExpectQuery(`INSERT INTO comments`).WillReturnError(sql.ErrConnDone)
		require.Error(t, repo.Create(&Comment{PostID: 1}))
	})
}

// ============================================================
// GetByID
// ============================================================

func TestCommentRepo_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewCommentRepoWithDB(db)

		now := time.Now()
		mock.ExpectQuery(`SELECT .+ FROM comments WHERE .+`).WillReturnRows(
			sqlmock.NewRows(commentColumns).AddRow(1, 1, 1, nil, "hello", 0, 1, now, now, nil))

		c, err := repo.GetByID(1)
		require.NoError(t, err)
		assert.Equal(t, int64(1), c.ID)
	})

	t.Run("not_found", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewCommentRepoWithDB(db)

		mock.ExpectQuery(`SELECT .+ FROM comments WHERE .+`).WillReturnError(gorm.ErrRecordNotFound)
		_, err := repo.GetByID(999)
		require.Error(t, err)
	})
}

// ============================================================
// GetByPostID
// ============================================================

func TestCommentRepo_GetByPostID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewCommentRepoWithDB(db)

		mock.ExpectQuery(`SELECT count\(\*\) FROM comments WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		now := time.Now()
		mock.ExpectQuery(`SELECT .+ FROM comments WHERE .+ ORDER BY`).WillReturnRows(
			sqlmock.NewRows(commentColumns).
				AddRow(1, 1, 1, nil, "c1", 0, 1, now, now, nil).
				AddRow(2, 1, 2, nil, "c2", 0, 1, now, now, nil))

		comments, total, err := repo.GetByPostID(1, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, comments, 2)
	})

	t.Run("empty", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewCommentRepoWithDB(db)

		mock.ExpectQuery(`SELECT count\(\*\) FROM comments WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		mock.ExpectQuery(`SELECT .+ FROM comments WHERE .+ ORDER BY`).WillReturnRows(sqlmock.NewRows(commentColumns))

		comments, total, err := repo.GetByPostID(999, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, comments)
	})
}

// ============================================================
// GetReplies
// ============================================================

func TestCommentRepo_GetReplies(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepoWithDB(db)

	now := time.Now()
	pid := int64(1)
	mock.ExpectQuery(`SELECT .+ FROM comments WHERE .+`).WillReturnRows(
		sqlmock.NewRows(commentColumns).
			AddRow(2, 1, 2, &pid, "reply", 0, 1, now, now, nil))

	replies, err := repo.GetReplies(1)
	require.NoError(t, err)
	assert.Len(t, replies, 1)
}

// ============================================================
// Delete
// ============================================================

func TestCommentRepo_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepoWithDB(db)

	mock.ExpectExec(`UPDATE comments SET status`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.Delete(1))
}

// ============================================================
// CountByPostID
// ============================================================

func TestCommentRepo_CountByPostID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM comments WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	count, err := repo.CountByPostID(1)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}
