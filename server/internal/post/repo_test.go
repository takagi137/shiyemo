package post

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

var postColumns = []string{
	"id", "user_id", "type", "title", "content", "images", "tags",
	"like_count", "comment_count", "is_top", "is_featured",
	"status", "review_status", "created_at", "updated_at", "deleted_at",
}

// ============================================================
// Create
// ============================================================

func TestPostRepo_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewPostRepoWithDB(db)

		mock.ExpectQuery(`INSERT INTO posts`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		post := &Post{UserID: 1, Type: "emotion", Title: "test", Content: "hello"}
		require.NoError(t, repo.Create(post))
		assert.Equal(t, int64(1), post.ID)
	})

	t.Run("db_error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewPostRepoWithDB(db)

		mock.ExpectQuery(`INSERT INTO posts`).WillReturnError(sql.ErrConnDone)
		require.Error(t, repo.Create(&Post{UserID: 1}))
	})
}

// ============================================================
// GetByID
// ============================================================

func TestPostRepo_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewPostRepoWithDB(db)

		now := time.Now()
		mock.ExpectQuery(`SELECT .+ FROM posts WHERE .+`).WillReturnRows(
			sqlmock.NewRows(postColumns).AddRow(1, 1, "emotion", "p", "c", "", "",
				5, 2, false, false, 1, "approved", now, now, nil))

		post, err := repo.GetByID(1)
		require.NoError(t, err)
		assert.Equal(t, int64(1), post.ID)
	})

	t.Run("not_found", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewPostRepoWithDB(db)

		mock.ExpectQuery(`SELECT .+ FROM posts WHERE .+`).WillReturnError(gorm.ErrRecordNotFound)
		_, err := repo.GetByID(999)
		require.Error(t, err)
	})
}

// ============================================================
// List
// ============================================================

func TestPostRepo_List(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM posts WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	now := time.Now()
	mock.ExpectQuery(`SELECT .+ FROM posts WHERE .+ ORDER BY`).WillReturnRows(
		sqlmock.NewRows(postColumns).
			AddRow(1, 1, "emotion", "p1", "c", "", "", 5, 2, false, false, 1, "approved", now, now, nil).
			AddRow(2, 1, "insight", "p2", "c", "", "", 3, 1, false, false, 1, "approved", now, now, nil))

	posts, total, err := repo.List("", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, posts, 2)
}

// ============================================================
// Update
// ============================================================

func TestPostRepo_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostRepoWithDB(db)

	mock.ExpectExec(`UPDATE posts SET`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.Update(&Post{ID: 1, Title: "updated"}))
}

// ============================================================
// Delete
// ============================================================

func TestPostRepo_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostRepoWithDB(db)

	mock.ExpectExec(`UPDATE posts SET status`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.Delete(1))
}

// ============================================================
// Count
// ============================================================

func TestPostRepo_Count(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPostRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM posts WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	count, err := repo.Count()
	require.NoError(t, err)
	assert.Equal(t, int64(10), count)
}

// ============================================================
// GetByUserIDs
// ============================================================

func TestPostRepo_GetByUserIDs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewPostRepoWithDB(db)

		mock.ExpectQuery(`SELECT count\(\*\) FROM posts WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		now := time.Now()
		mock.ExpectQuery(`SELECT .+ FROM posts WHERE .+ ORDER BY`).WillReturnRows(
			sqlmock.NewRows(postColumns).
				AddRow(1, 1, "emotion", "p1", "", "", "", 0, 0, false, false, 1, "approved", now, now, nil).
				AddRow(2, 2, "insight", "p2", "", "", "", 0, 0, false, false, 1, "approved", now, now, nil))

		posts, total, err := repo.GetByUserIDs([]int64{1, 2}, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, posts, 2)
	})

	t.Run("empty", func(t *testing.T) {
		db, _ := setupMockDB(t)
		repo := NewPostRepoWithDB(db)

		posts, total, err := repo.GetByUserIDs([]int64{}, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, posts)
	})
}
