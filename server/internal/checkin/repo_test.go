package checkin

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestCheckinRepo_Create(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewCheckinRepoWithDB(gormDB)

	checkin := &DailyCheckin{
		UserID:   1,
		Mood:     MoodHappy,
		Content:  "今天心情不错",
		IsPublic: true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "daily_checkins"`).
		WithArgs(checkin.UserID, checkin.Mood, checkin.Content, checkin.Images, checkin.IsPublic, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(checkin)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), checkin.ID)
}

func TestCheckinRepo_GetTodayByUserID(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewCheckinRepoWithDB(gormDB)

	userID := int64(1)
	today := time.Now().Format("2006-01-02")

	rows := sqlmock.NewRows([]string{"id", "user_id", "mood", "content", "images", "is_public", "created_at", "updated_at"}).
		AddRow(1, userID, MoodHappy, "今天心情不错", "[]", true, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "daily_checkins" WHERE user_id = \$1 AND DATE\(created_at\) = \$2`).
		WithArgs(userID, today).
		WillReturnRows(rows)

	checkin, err := repo.GetTodayByUserID(userID)
	assert.NoError(t, err)
	assert.NotNil(t, checkin)
	assert.Equal(t, userID, checkin.UserID)
	assert.Equal(t, MoodHappy, checkin.Mood)
}

func TestCheckinRepo_GetTodayByUserID_NotFound(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewCheckinRepoWithDB(gormDB)

	userID := int64(1)
	today := time.Now().Format("2006-01-02")

	mock.ExpectQuery(`SELECT \* FROM "daily_checkins" WHERE user_id = \$1 AND DATE\(created_at\) = \$2`).
		WithArgs(userID, today).
		WillReturnError(gorm.ErrRecordNotFound)

	checkin, err := repo.GetTodayByUserID(userID)
	assert.Error(t, err)
	assert.Nil(t, checkin)
}

func TestCheckinRepo_GetHistoryByUserID(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewCheckinRepoWithDB(gormDB)

	userID := int64(1)
	page := 1
	pageSize := 10

	// Mock count query
	mock.ExpectQuery(`SELECT count\(\*\) FROM "daily_checkins" WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	// Mock select query
	rows := sqlmock.NewRows([]string{"id", "user_id", "mood", "content", "images", "is_public", "created_at", "updated_at"}).
		AddRow(2, userID, MoodCalm, "今天很平静", "[]", true, time.Now(), time.Now()).
		AddRow(1, userID, MoodHappy, "今天心情不错", "[]", true, time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, -1))

	mock.ExpectQuery(`SELECT \* FROM "daily_checkins" WHERE user_id = \$1 ORDER BY created_at DESC LIMIT 10`).
		WithArgs(userID).
		WillReturnRows(rows)

	checkins, total, err := repo.GetHistoryByUserID(userID, page, pageSize)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, checkins, 2)
}

func TestCheckinRepo_GetMoodCurve(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewCheckinRepoWithDB(gormDB)

	userID := int64(1)
	days := 30
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	rows := sqlmock.NewRows([]string{"id", "user_id", "mood", "content", "images", "is_public", "created_at", "updated_at"}).
		AddRow(1, userID, MoodHappy, "第1天", "[]", true, time.Now().AddDate(0, 0, -2), time.Now().AddDate(0, 0, -2)).
		AddRow(2, userID, MoodCalm, "第2天", "[]", true, time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, -1))

	mock.ExpectQuery(`SELECT \* FROM "daily_checkins" WHERE user_id = \$1 AND DATE\(created_at\) >= \$2 ORDER BY created_at ASC`).
		WithArgs(userID, startDate).
		WillReturnRows(rows)

	checkins, err := repo.GetMoodCurve(userID, days)
	assert.NoError(t, err)
	assert.Len(t, checkins, 2)
}

func TestCheckinRepo_GetPublicTimeline(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	repo := NewCheckinRepoWithDB(gormDB)

	page := 1
	pageSize := 10

	// Mock count query
	mock.ExpectQuery(`SELECT count\(\*\) FROM "daily_checkins" WHERE is_public = \$1`).
		WithArgs(true).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Mock select query
	rows := sqlmock.NewRows([]string{"id", "user_id", "mood", "content", "images", "is_public", "created_at", "updated_at"}).
		AddRow(1, 1, MoodHappy, "今天心情不错", "[]", true, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "daily_checkins" WHERE is_public = \$1 ORDER BY created_at DESC LIMIT 10`).
		WithArgs(true).
		WillReturnRows(rows)

	checkins, total, err := repo.GetPublicTimeline(page, pageSize)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, checkins, 1)
}
