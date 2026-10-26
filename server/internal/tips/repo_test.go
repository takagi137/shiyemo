package tips

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

// setupMockDB 创建 gorm.DB + sqlmock
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

// tipColumns 贴士表列名
var tipColumns = []string{
	"id", "content", "category", "is_active", "created_at",
}

// ============================================================
// GetRandom 测试
// ============================================================

func TestTipsRepo_GetRandom(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewTipsRepoWithDB(gormDB)

		now := time.Now()
		rows := sqlmock.NewRows(tipColumns).AddRow(
			1, "见过加班猝死的，没见过不上班饿死的。", "mindset", true, now,
		)
		mock.ExpectQuery(`SELECT \* FROM survival_tips WHERE is_active = .+ ORDER BY RANDOM\(\)`).
			WillReturnRows(rows)

		tip, err := repo.GetRandom()
		require.NoError(t, err)
		assert.Equal(t, int64(1), tip.ID)
		assert.Equal(t, "见过加班猝死的，没见过不上班饿死的。", tip.Content)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewTipsRepoWithDB(gormDB)

		mock.ExpectQuery(`SELECT \* FROM survival_tips WHERE is_active = .+ ORDER BY RANDOM\(\)`).
			WillReturnError(gorm.ErrRecordNotFound)

		tip, err := repo.GetRandom()
		require.Error(t, err)
		assert.Nil(t, tip)
	})
}

// ============================================================
// GetList 测试
// ============================================================

func TestTipsRepo_GetList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewTipsRepoWithDB(gormDB)

		now := time.Now()

		// 统计总数
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(`SELECT count\(\*\) FROM survival_tips WHERE is_active =`).
			WillReturnRows(countRows)

		// 分页查询
		rows := sqlmock.NewRows(tipColumns).
			AddRow(2, "一个人生活自己做饭，每月 800 就能吃得很好。", "finance", true, now).
			AddRow(1, "见过加班猝死的，没见过不上班饿死的。", "mindset", true, now)
		mock.ExpectQuery(`SELECT \* FROM survival_tips WHERE is_active = .+ ORDER BY id DESC`).
			WillReturnRows(rows)

		tips, total, err := repo.GetList(1, 20)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, tips, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewTipsRepoWithDB(gormDB)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(`SELECT count\(\*\) FROM survival_tips WHERE is_active =`).
			WillReturnRows(countRows)

		rows := sqlmock.NewRows(tipColumns)
		mock.ExpectQuery(`SELECT \* FROM survival_tips WHERE is_active = .+ ORDER BY id DESC`).
			WillReturnRows(rows)

		tips, total, err := repo.GetList(1, 20)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, tips)
	})
}

// ============================================================
// Create 测试
// ============================================================

func TestTipsRepo_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewTipsRepoWithDB(gormDB)

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(`INSERT INTO survival_tips`).
			WillReturnRows(rows)

		tip := &SurvivalTip{
			Content:  "测试贴士",
			Category: "general",
			IsActive: true,
		}

		err := repo.Create(tip)
		require.NoError(t, err)
		assert.Equal(t, int64(1), tip.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewTipsRepoWithDB(gormDB)

		mock.ExpectQuery(`INSERT INTO survival_tips`).
			WillReturnError(sql.ErrConnDone)

		tip := &SurvivalTip{
			Content:  "测试贴士",
			Category: "general",
			IsActive: true,
		}

		err := repo.Create(tip)
		require.Error(t, err)
	})
}

// ============================================================
// Update 测试
// ============================================================

func TestTipsRepo_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewTipsRepoWithDB(gormDB)

		mock.ExpectExec(`UPDATE survival_tips`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		tip := &SurvivalTip{
			ID:       1,
			Content:  "更新后的贴士",
			Category: "finance",
		}

		err := repo.Update(tip)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewTipsRepoWithDB(gormDB)

		mock.ExpectExec(`UPDATE survival_tips`).
			WillReturnError(sql.ErrConnDone)

		tip := &SurvivalTip{
			ID:       1,
			Content:  "更新后的贴士",
			Category: "finance",
		}

		err := repo.Update(tip)
		require.Error(t, err)
	})
}

// ============================================================
// Delete 测试
// ============================================================

func TestTipsRepo_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewTipsRepoWithDB(gormDB)

		mock.ExpectExec(`UPDATE survival_tips SET .+ WHERE id =`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(1)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		gormDB, mock := setupMockDB(t)
		repo := NewTipsRepoWithDB(gormDB)

		mock.ExpectExec(`UPDATE survival_tips SET .+ WHERE id =`).
			WillReturnError(sql.ErrConnDone)

		err := repo.Delete(1)
		require.Error(t, err)
	})
}
