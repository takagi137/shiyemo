package setting

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupMockDB 创建 gorm.DB + sqlmock（与其他模块保持一致，测试不依赖真实 PostgreSQL）
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

var settingColumns = []string{"key", "value", "updated_at"}

func TestSettingRepo_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := setupMockDB(t)
		rows := sqlmock.NewRows(settingColumns).
			AddRow("smtp.host", "smtp.qq.com", time.Now()).
			AddRow("smtp.enabled", "true", time.Now())
		mock.ExpectQuery(`SELECT \* FROM system_settings`).WillReturnRows(rows)

		got, err := NewSettingRepoWithDB(db).GetAll()

		require.NoError(t, err)
		assert.Equal(t, "smtp.qq.com", got["smtp.host"])
		assert.Equal(t, "true", got["smtp.enabled"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("空表返回空 map 而非 nil", func(t *testing.T) {
		db, mock := setupMockDB(t)
		mock.ExpectQuery(`SELECT \* FROM system_settings`).
			WillReturnRows(sqlmock.NewRows(settingColumns))

		got, err := NewSettingRepoWithDB(db).GetAll()

		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Len(t, got, 0)
	})
}

func TestSettingRepo_Get(t *testing.T) {
	t.Run("命中", func(t *testing.T) {
		db, mock := setupMockDB(t)
		mock.ExpectQuery(`SELECT \* FROM system_settings WHERE key = \$1`).
			WithArgs("smtp.host").
			WillReturnRows(sqlmock.NewRows(settingColumns).AddRow("smtp.host", "smtp.163.com", time.Now()))

		got, err := NewSettingRepoWithDB(db).Get("smtp.host")

		require.NoError(t, err)
		assert.Equal(t, "smtp.163.com", got)
	})

	// key 不存在属于正常情况（尚未配置），不应作为错误抛出
	t.Run("key 不存在时返回空字符串且不报错", func(t *testing.T) {
		db, mock := setupMockDB(t)
		mock.ExpectQuery(`SELECT \* FROM system_settings WHERE key = \$1`).
			WithArgs("smtp.nonexistent").
			WillReturnRows(sqlmock.NewRows(settingColumns))

		got, err := NewSettingRepoWithDB(db).Get("smtp.nonexistent")

		require.NoError(t, err)
		assert.Equal(t, "", got)
	})
}

func TestSettingRepo_Set(t *testing.T) {
	t.Run("upsert 写入", func(t *testing.T) {
		db, mock := setupMockDB(t)
		mock.ExpectExec(`INSERT INTO system_settings`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := NewSettingRepoWithDB(db).Set("smtp.host", "smtp.qq.com")

		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSettingRepo_SetMany(t *testing.T) {
	// 批量写入必须在单事务内，避免部分成功导致配置处于中间状态
	t.Run("在单个事务内写入", func(t *testing.T) {
		db, mock := setupMockDB(t)
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO system_settings`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := NewSettingRepoWithDB(db).SetMany(map[string]string{"smtp.host": "smtp.qq.com"})

		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
