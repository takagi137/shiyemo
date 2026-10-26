package notification

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

func TestNotificationRepo_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectQuery(`INSERT INTO notifications`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	n := &Notification{UserID: 1, Type: "like", Content: "liked your post"}
	require.NoError(t, repo.Create(n))
}

func TestNotificationRepo_ListByUserID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM notifications WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT .+ FROM notifications WHERE .+ ORDER BY`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "user_id", "type", "content", "target_id", "target_type", "is_read", "created_at"}).
			AddRow(1, 1, "like", "liked", nil, nil, false, nil))

	notifications, total, err := repo.ListByUserID(1, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, notifications, 1)
}

func TestNotificationRepo_ListByType(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM notifications WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT .+ FROM notifications WHERE .+ ORDER BY`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "user_id", "type", "content", "target_id", "target_type", "is_read", "created_at"}).
			AddRow(1, 1, "like", "liked", nil, nil, false, nil))

	notifications, total, err := repo.ListByType(1, "like", false, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, notifications, 1)
}

func TestNotificationRepo_BatchMarkAsRead(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectExec(`UPDATE notifications SET .+ WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 2))
	require.NoError(t, repo.BatchMarkAsRead([]int64{1, 2}, 1))
}

func TestNotificationRepo_BatchDelete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectExec(`DELETE FROM notifications WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 2))
	require.NoError(t, repo.BatchDelete([]int64{1, 2}, 1))
}

func TestNotificationRepo_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectQuery(`SELECT .+ FROM notifications WHERE .+`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "user_id", "type", "content", "target_id", "target_type", "is_read", "created_at"}).
			AddRow(1, 1, "like", "liked", nil, nil, false, nil))

	n, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, int64(1), n.ID)
}

func TestNotificationRepo_MarkAsRead(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectExec(`UPDATE notifications SET .+ WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.MarkAsRead(1, 1))
}

func TestNotificationRepo_MarkAllAsRead(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectExec(`UPDATE notifications SET .+ WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 5))
	require.NoError(t, repo.MarkAllAsRead(1))
}

func TestNotificationRepo_CountUnread(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM notifications WHERE .+`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	count, err := repo.CountUnread(1)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}
