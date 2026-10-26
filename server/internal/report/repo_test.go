package report

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

func TestReportRepo_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewReportRepoWithDB(db)

	mock.ExpectQuery(`INSERT INTO reports`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	r := &Report{UserID: 1, TargetType: "post", TargetID: 10, Reason: "spam"}
	require.NoError(t, repo.Create(r))
	assert.Equal(t, int64(1), r.ID)
}

func TestReportRepo_ExistsByUserAndTarget(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewReportRepoWithDB(db)

	t.Run("exists", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM reports`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		exists, err := repo.ExistsByUserAndTarget(1, "post", 10)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("not_exists", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM reports`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		exists, err := repo.ExistsByUserAndTarget(1, "post", 10)
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestReportRepo_List(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewReportRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM reports`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery(`SELECT .+ FROM reports`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "user_id", "target_type", "target_id", "reason", "status"}).
			AddRow(1, 1, "post", 10, "spam", 1).
			AddRow(2, 2, "comment", 5, "abuse", 1))

	reports, total, err := repo.List(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, reports, 2)
}

func TestReportRepo_UpdateStatus(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewReportRepoWithDB(db)

	mock.ExpectExec(`UPDATE reports SET .+ WHERE id = .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.UpdateStatus(1, 2))
}
