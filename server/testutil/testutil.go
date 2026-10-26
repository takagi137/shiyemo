package testutil

import (
	"fmt"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

// User 测试用用户模型
type User struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	Phone          *string   `gorm:"uniqueIndex;size:20"`
	Email          *string   `gorm:"uniqueIndex;size:255"`
	Password       string    `gorm:"not null;size:255"`
	Nickname       string    `gorm:"not null;size:50;index"`
	Bio            string    `gorm:"size:500;default:''"`
	Company        string    `gorm:"size:100;default:''"`
	Occupation     string    `gorm:"size:100;default:''"`
	Avatar         string    `gorm:"size:500;default:''"`
	FollowingCount int       `gorm:"default:0"`
	FollowerCount  int       `gorm:"default:0"`
	Status         int8      `gorm:"default:1;index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Post 测试用帖子模型
type Post struct {
	ID           int64          `gorm:"primaryKey;autoIncrement"`
	UserID       int64          `gorm:"index"`
	Type         string         `gorm:"size:20;not null"`
	Title        string         `gorm:"size:200;not null"`
	Content      string         `gorm:"type:text;not null"`
	Images       string         `gorm:"type:jsonb;default:'[]'"`
	Tags         string         `gorm:"type:jsonb;default:'[]'"`
	LikeCount    int            `gorm:"default:0"`
	CommentCount int            `gorm:"default:0"`
	Status       int            `gorm:"default:1;index"`
	ReviewStatus string         `gorm:"size:20;default:'pending';index"`
	RejectReason string         `gorm:"size:500;default:''"`
	ReviewerID   *int64         `gorm:"index"`
	ReviewedAt   *time.Time
	IsTop        bool           `gorm:"default:false;index"`
	IsFeatured   bool           `gorm:"default:false;index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// Comment 测试用评论模型
type Comment struct {
	ID        int64      `gorm:"primaryKey;autoIncrement"`
	PostID    int64      `gorm:"index"`
	UserID    int64      `gorm:"index"`
	ParentID  *int64     `gorm:"index"`
	Content   string     `gorm:"type:text"`
	LikeCount int        `gorm:"default:0"`
	Status    int        `gorm:"default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Like 测试用点赞模型
type Like struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64 `gorm:"index"`
	PostID    int64 `gorm:"index"`
	CreatedAt time.Time
}

// Bookmark 测试用收藏模型
type Bookmark struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64 `gorm:"index"`
	PostID    int64 `gorm:"index"`
	CreatedAt time.Time
}

// Follow 测试用关注模型
type Follow struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	FollowerID  int64 `gorm:"index"`
	FollowingID int64 `gorm:"index"`
	CreatedAt   time.Time
}

// SetupTestDB 初始化测试数据库连接和 Redis
func SetupTestDB() {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		// 默认测试数据库 DSN
		dsn = "host=localhost user=shiyemo password=shiyemo123 dbname=shiyemo_test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	}

	var err error
	TestDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("连接测试数据库失败: %v", err))
	}

	// 自动迁移表结构
	err = TestDB.AutoMigrate(
		&User{},
		&Post{},
		&Comment{},
		&Like{},
		&Bookmark{},
		&Follow{},
	)
	if err != nil {
		panic(fmt.Sprintf("自动迁移失败: %v", err))
	}

	// 初始化 Redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}

	// 这里不初始化 Redis，因为 cache 包有自己的初始化逻辑
	// 测试中会单独调用 cache.Init()
}

// CleanupTestData 清理测试数据
func CleanupTestData(t *testing.T) {
	t.Helper()

	// 使用 TRUNCATE CASCADE 一次性清空所有表，避免外键约束问题
	TestDB.Exec("TRUNCATE TABLE notifications, comments, likes, bookmarks, follows, posts, users RESTART IDENTITY CASCADE")
}

// CreateUser 直接创建用户（不依赖业务包）
func CreateUser(t *testing.T, email, nickname, password string) int64 {
	t.Helper()

	user := &User{
		Email:    &email,
		Nickname: nickname,
		Password: password,
		Status:   1,
	}

	if err := TestDB.Create(user).Error; err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	return user.ID
}

// CreatePost 直接创建帖子（不依赖业务包）
func CreatePost(t *testing.T, userID int64, postType, title, content string) int64 {
	t.Helper()

	post := &Post{
		UserID:  userID,
		Type:    postType,
		Title:   title,
		Content: content,
		Status:  1,
	}

	if err := TestDB.Create(post).Error; err != nil {
		t.Fatalf("创建测试帖子失败: %v", err)
	}

	return post.ID
}

// CreatePostWithStatus 创建指定状态的帖子
func CreatePostWithStatus(t *testing.T, userID int64, postType, title, content, reviewStatus string) int64 {
	t.Helper()

	post := &Post{
		UserID:       userID,
		Type:         postType,
		Title:        title,
		Content:      content,
		Status:       1,
		ReviewStatus: reviewStatus,
	}

	if err := TestDB.Create(post).Error; err != nil {
		t.Fatalf("创建测试帖子失败: %v", err)
	}

	return post.ID
}

// GetPost 获取帖子信息
func GetPost(t *testing.T, postID int64) map[string]interface{} {
	t.Helper()

	var post Post
	if err := TestDB.Where("id = ?", postID).First(&post).Error; err != nil {
		t.Fatalf("获取帖子失败: %v", err)
	}

	return map[string]interface{}{
		"id":            post.ID,
		"user_id":       post.UserID,
		"type":          post.Type,
		"title":         post.Title,
		"content":       post.Content,
		"review_status": post.ReviewStatus,
		"reject_reason": post.RejectReason,
		"is_top":        post.IsTop,
		"is_featured":   post.IsFeatured,
	}
}
