package follow

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
// Follow / Unfollow / IsFollowing
// ============================================================

func TestFollowRepo_Follow(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewFollowRepoWithDB(db)

	mock.ExpectQuery(`INSERT INTO follows`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	require.NoError(t, repo.Follow(1, 2))
}

func TestFollowRepo_Unfollow(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewFollowRepoWithDB(db)

	mock.ExpectExec(`DELETE FROM follows WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.Unfollow(1, 2))
}

func TestFollowRepo_IsFollowing(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewFollowRepoWithDB(db)

	t.Run("true", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM follows WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		assert.True(t, repo.IsFollowing(1, 2))
	})

	t.Run("false", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM follows WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		assert.False(t, repo.IsFollowing(1, 2))
	})
}

// ============================================================
// ListFollowing / ListFollowers
// ============================================================

func TestFollowRepo_ListFollowing(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewFollowRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM follows WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery(`SELECT .+ FROM follows WHERE .+ ORDER BY`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "follower_id", "following_id", "created_at"}).
			AddRow(1, 1, 10, nil).
			AddRow(2, 1, 20, nil))

	ids, total, err := repo.ListFollowing(1, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, []int64{10, 20}, ids)
}

func TestFollowRepo_ListFollowers(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewFollowRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM follows WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT .+ FROM follows WHERE .+ ORDER BY`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "follower_id", "following_id", "created_at"}).
			AddRow(1, 10, 1, nil))

	ids, total, err := repo.ListFollowers(1, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, []int64{10}, ids)
}

// ============================================================
// CountFollowing / CountFollowers
// ============================================================

func TestFollowRepo_CountFollowing(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewFollowRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM follows WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	count, err := repo.CountFollowing(1)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestFollowRepo_CountFollowers(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewFollowRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM follows WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	count, err := repo.CountFollowers(1)
	require.NoError(t, err)
	assert.Equal(t, int64(10), count)
}
