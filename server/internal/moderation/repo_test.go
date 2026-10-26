package moderation

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

func TestModerationRepo_CreateWord(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewModerationRepoWithDB(db)

	mock.ExpectQuery(`INSERT INTO sensitive_words`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	w := &SensitiveWord{Word: "badword"}
	require.NoError(t, repo.CreateWord(w))
}

func TestModerationRepo_GetWordByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewModerationRepoWithDB(db)

	mock.ExpectQuery(`SELECT .+ FROM sensitive_words`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "word"}).AddRow(1, "badword"))

	w, err := repo.GetWordByID(1)
	require.NoError(t, err)
	assert.Equal(t, "badword", w.Word)
}

func TestModerationRepo_GetWordByName(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewModerationRepoWithDB(db)

	t.Run("found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM sensitive_words`).WillReturnRows(
			sqlmock.NewRows([]string{"id", "word"}).AddRow(1, "badword"))
		w, err := repo.GetWordByName("badword")
		require.NoError(t, err)
		assert.NotNil(t, w)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT .+ FROM sensitive_words`).WillReturnError(gorm.ErrRecordNotFound)
		w, err := repo.GetWordByName("notexist")
		require.NoError(t, err)
		assert.Nil(t, w)
	})
}

func TestModerationRepo_ListWords(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewModerationRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM sensitive_words`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery(`SELECT .+ FROM sensitive_words`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "word"}).AddRow(1, "w1").AddRow(2, "w2"))

	words, total, err := repo.ListWords(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, words, 2)
}

func TestModerationRepo_DeleteWord(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewModerationRepoWithDB(db)

	mock.ExpectExec(`DELETE FROM sensitive_words`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.DeleteWord(1))
}

func TestModerationRepo_GetAllWords(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewModerationRepoWithDB(db)

	mock.ExpectQuery(`SELECT word FROM sensitive_words`).WillReturnRows(
		sqlmock.NewRows([]string{"word"}).AddRow("w1").AddRow("w2"))

	words, err := repo.GetAllWords()
	require.NoError(t, err)
	assert.Equal(t, []string{"w1", "w2"}, words)
}
