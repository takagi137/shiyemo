package admin

import (
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
		Conn: db, WithoutQuotingCheck: true, PreferSimpleProtocol: true,
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent), SkipDefaultTransaction: true})
	require.NoError(t, err)
	return gormDB, mock
}

func TestAdminRepo_ListPosts_All(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM posts WHERE 1=1`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery(`SELECT .+ FROM posts WHERE 1=1 ORDER BY`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "user_id", "type", "title", "content", "status", "created_at", "deleted_at"}).
			AddRow(1, 1, "emotion", "p1", "c1", 1, nil, nil).
			AddRow(2, 2, "insight", "p2", "c2", 1, nil, nil))

	posts, total, err := repo.ListPosts(1, 10, "all")
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, posts, 2)
}

func TestAdminRepo_ListUsers(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM users`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
	mock.ExpectQuery(`SELECT .+ FROM users ORDER BY`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "nickname", "email", "status"}).AddRow(1, "u1", "a@b.com", 1))

	users, total, err := repo.ListUsers(1, 10, "", nil)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, users, 1)
}

func TestAdminRepo_UpdateUserStatus(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectExec(`UPDATE users SET .+ WHERE id = .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.UpdateUserStatus(1, 0))
}

func TestAdminRepo_GetStats(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM users`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
	mock.ExpectQuery(`SELECT count\(\*\) FROM posts WHERE deleted_at IS NULL`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(20))
	mock.ExpectQuery(`SELECT count\(\*\) FROM comments`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(50))

	stats, err := repo.GetStats()
	require.NoError(t, err)
	assert.Equal(t, int64(10), stats["total_users"])
	assert.Equal(t, int64(20), stats["total_posts"])
	assert.Equal(t, int64(50), stats["total_comments"])
}

func TestAdminRepo_ApprovePost(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectExec(`UPDATE posts SET .+ WHERE id = .+$`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.ApprovePost(1, 1))
}

func TestAdminRepo_ApprovePost_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectExec(`UPDATE posts SET .+ WHERE id = .+$`).WillReturnResult(sqlmock.NewResult(0, 0))
	err := repo.ApprovePost(999, 1)
	require.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestAdminRepo_RejectPost(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectExec(`UPDATE posts SET .+ WHERE id = .+$`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.RejectPost(1, 1, "violation"))
}

func TestAdminRepo_RejectPost_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectExec(`UPDATE posts SET .+ WHERE id = .+$`).WillReturnResult(sqlmock.NewResult(0, 0))
	err := repo.RejectPost(999, 1, "reason")
	require.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
