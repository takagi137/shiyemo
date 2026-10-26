package tag

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

func TestTagRepo_GetByName(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewTagRepoWithDB(db)

		mock.ExpectQuery(`SELECT .+ FROM tags`).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Go"))

		tag, err := repo.GetByName("Go")
		require.NoError(t, err)
		assert.Equal(t, "Go", tag.Name)
	})

	t.Run("not_found", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewTagRepoWithDB(db)

		mock.ExpectQuery(`SELECT .+ FROM tags`).WillReturnError(gorm.ErrRecordNotFound)
		_, err := repo.GetByName("NotExist")
		require.Error(t, err)
	})
}

func TestTagRepo_GetOrCreate(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewTagRepoWithDB(db)

		mock.ExpectQuery(`SELECT .+ FROM tags`).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Go"))

		tag, err := repo.GetOrCreate("Go")
		require.NoError(t, err)
		assert.Equal(t, "Go", tag.Name)
	})

	t.Run("new_tag", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewTagRepoWithDB(db)

		mock.ExpectQuery(`SELECT .+ FROM tags`).WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectQuery(`INSERT INTO tags`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		tag, err := repo.GetOrCreate("Python")
		require.NoError(t, err)
		assert.NotNil(t, tag)
	})
}

func TestTagRepo_AddUserSkill(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTagRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM user_skill_tags`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec(`INSERT INTO user_skill_tags`).WillReturnResult(sqlmock.NewResult(1, 1))

	require.NoError(t, repo.AddUserSkill(1, 1))
}

func TestTagRepo_RemoveUserSkill(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTagRepoWithDB(db)

	mock.ExpectExec(`DELETE FROM user_skill_tags WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, repo.RemoveUserSkill(1, 1))
}

func TestTagRepo_GetUserSkills(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTagRepoWithDB(db)

	mock.ExpectQuery(`SELECT .+ FROM tags .+ JOIN user_skill_tags`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Go").AddRow(2, "Python"))

	skills, err := repo.GetUserSkills(1)
	require.NoError(t, err)
	assert.Len(t, skills, 2)
}

func TestTagRepo_GetUserIDsByTag(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTagRepoWithDB(db)

	mock.ExpectQuery(`SELECT user_id FROM user_skill_tags WHERE .+`).WillReturnRows(
		sqlmock.NewRows([]string{"user_id"}).AddRow(1).AddRow(2))

	ids, err := repo.GetUserIDsByTag(1)
	require.NoError(t, err)
	assert.Equal(t, []int64{1, 2}, ids)
}

func TestTagRepo_GetPostTags(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTagRepoWithDB(db)

	mock.ExpectQuery(`SELECT .+ FROM tags .+ JOIN post_tags`).WillReturnRows(
		sqlmock.NewRows([]string{"name"}).AddRow("Go").AddRow("Python"))

	names, err := repo.GetPostTags(1)
	require.NoError(t, err)
	assert.Equal(t, []string{"Go", "Python"}, names)
}

func TestTagRepo_UpdatePostTags(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTagRepoWithDB(db)

	mock.ExpectExec(`DELETE FROM post_tags WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO post_tags`).WillReturnResult(sqlmock.NewResult(1, 2))

	err := repo.UpdatePostTags(1, []int64{1, 2})
	require.NoError(t, err)
}

func TestTagRepo_UpdatePostTags_Empty(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTagRepoWithDB(db)

	mock.ExpectExec(`DELETE FROM post_tags WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdatePostTags(1, []int64{})
	require.NoError(t, err)
}

func TestTagRepo_GetHotTags(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTagRepoWithDB(db)

	mock.ExpectQuery(`SELECT .+ FROM tags .+ LEFT JOIN post_tags`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "count"}).
			AddRow(1, "Go", 10).
			AddRow(2, "Python", 5))

	tags, err := repo.GetHotTags(10)
	require.NoError(t, err)
	assert.Len(t, tags, 2)
}

func TestTagRepo_AddUserSkill_AlreadyExists(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTagRepoWithDB(db)

	mock.ExpectQuery(`SELECT count\(\*\) FROM user_skill_tags`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	err := repo.AddUserSkill(1, 1)
	require.NoError(t, err)
}
