package integration

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/admin"
	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/comment"
	"github.com/shiyemo/server/internal/common/database"
	"github.com/shiyemo/server/internal/common/middleware"
	"github.com/shiyemo/server/internal/common/response"
	"github.com/shiyemo/server/internal/follow"
	"github.com/shiyemo/server/internal/interaction"
	"github.com/shiyemo/server/internal/migrationx"
	"github.com/shiyemo/server/internal/notification"
	"github.com/shiyemo/server/internal/post"
	"github.com/shiyemo/server/internal/report"
	"github.com/shiyemo/server/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	testDB        *gorm.DB
	migrationsDir string
	dbHost        string
	dbPort        string
	dbUser        string
	dbPass        string
	testDBName    string
)

// TestMain 集成测试入口
// 1. 创建/重建测试数据库（通过 gorm 建立与 postgres 的连接）
// 2. 读取并执行迁移 SQL
// 3. 运行测试
// 4. 删除测试数据库
func TestMain(m *testing.M) {
	// ❶ 环境变量 — 在 config.Load() 之前设置
	os.Setenv("JWT_SECRET", "integration-test-secret")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "shiyemo")
	os.Setenv("DB_PASSWORD", "shiyemo123")
	os.Setenv("DB_NAME", "shiyemo")
	os.Setenv("APP_MODE", "debug")

	// ❷ 解析自定义参数（可通过环境变量覆盖）
	flag.StringVar(&migrationsDir, "migrations.dir", "", "迁移文件目录（相对于 server/）")
	flag.StringVar(&testDBName, "test.dbname", "shiyemo_test", "测试数据库名")
	flag.StringVar(&dbHost, "db.host", os.Getenv("TEST_DB_HOST"), "测试数据库主机")
	flag.StringVar(&dbPort, "db.port", os.Getenv("TEST_DB_PORT"), "测试数据库端口")
	flag.StringVar(&dbUser, "db.user", "shiyemo", "测试数据库用户")
	flag.StringVar(&dbPass, "db.pass", "shiyemo123", "测试数据库密码")
	flag.Parse()

	// ❸ 搜索迁移目录 — 向上查找 migrations/ 文件夹
	if migrationsDir == "" {
		migrationsDir = findMigrationsDir()
		if migrationsDir == "" {
			migrationsDir = "../../migrations" // fallback
		}
	}
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	log.Printf("[integration] 迁移目录: %s, 数据库: %s@%s:%s/%s", migrationsDir, dbUser, dbHost, dbPort, testDBName)

	// ❹ 创建/重建测试数据库（连接 postgres 系统库，DROP+CREATE）
	if err := recreateTestDB(); err != nil {
		log.Fatalf("[integration] 创建测试数据库失败: %v", err)
	}

	// ❺ 连接测试数据库
	if err := connectTestDB(); err != nil {
		log.Fatalf("[integration] 连接测试数据库失败: %v", err)
	}

	// ❻ 运行迁移
	if err := runMigrations(); err != nil {
		log.Fatalf("[integration] 运行迁移失败: %v", err)
	}
	log.Println("[integration] ✅ 迁移完成")

	// ❼ 运行测试
	code := m.Run()

	// ❽ 清理
	teardownTestDB()
	os.Exit(code)
}

// findMigrationsDir 向上搜索包含 migrations/ 的目录
// 处理两种场景：
//   1. 本地开发：tests/integration/ → 向上找 server/migrations/ 或直接 ./migrations/
//   2. Docker 挂载：/app/ 本身就是 server/ → ./migrations/
func findMigrationsDir() string {
	dir, _ := os.Getwd()
	for i := 0; i < 10; i++ {
		// 当前目录直接有 migrations/
		if info, err := os.Stat(filepath.Join(dir, "migrations")); err == nil && info.IsDir() {
			return filepath.Join(dir, "migrations")
		}
		// server/ 子目录有 migrations/
		if info, err := os.Stat(filepath.Join(dir, "server", "migrations")); err == nil && info.IsDir() {
			return filepath.Join(dir, "server", "migrations")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// recreateTestDB 连接到 postgres 系统库，DROP + CREATE 目标测试库
func recreateTestDB() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		dbHost, dbPort, dbUser, dbPass)

	// 通过 gorm 连接（不走 database 包，避免 driver 注册问题）
	adminDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("连接 postgres 失败（确认 PostgreSQL 正在运行）: %w", err)
	}

	rawDB, _ := adminDB.DB()
	defer rawDB.Close()

	// 终止 testDBName 的现有连接
	rawDB.Exec(fmt.Sprintf(
		"SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s' AND pid <> pg_backend_pid()",
		testDBName,
	))

	// DROP + CREATE
	rawDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	if _, err := rawDB.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName)); err != nil {
		return fmt.Errorf("CREATE DATABASE %s 失败: %w", testDBName, err)
	}

	log.Printf("[integration] 数据库 %s 已创建", testDBName)
	return nil
}

// connectTestDB 使用 gorm 连接到目标测试库
func connectTestDB() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, dbPort, dbUser, dbPass, testDBName,
	)

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return err
	}

	log.Printf("[integration] 已连接测试库 %s", testDBName)
	return nil
}

// runMigrations 读取并执行所有迁移 SQL（CREATE TABLE / ALTER TABLE / CREATE INDEX 等）
func runMigrations() error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("读取迁移目录 %s 失败: %w", migrationsDir, err)
	}

	var sqlFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			sqlFiles = append(sqlFiles, e.Name())
		}
	}
	migrationx.SortInPlace(sqlFiles)

	rawDB, err := testDB.DB()
	if err != nil {
		return fmt.Errorf("获取 raw DB 失败: %w", err)
	}

	for _, name := range sqlFiles {
		path := filepath.Join(migrationsDir, name)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取 %s 失败: %w", name, err)
		}

		for _, stmt := range splitSQL(string(content)) {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" || strings.HasPrefix(stmt, "--") {
				continue
			}
			// 迁移错误必须直接失败 —— 不能吞掉。
			// 历史教训：此处曾把错误降级为警告（"多为幂等语句"）后继续，
			// 导致迁移 012 的 read 列名错、缺失的建表迁移等全新部署必崩的 bug
			// 在集成测试里被静默放过，潜伏两个月。迁移失败 = 全新部署失败，
			// 必须让测试红掉。
			if _, err := rawDB.Exec(stmt); err != nil {
				return fmt.Errorf("迁移 %s 执行失败: %w\n语句: %s", name, err, strings.TrimSpace(stmt))
			}
		}
		log.Printf("[integration]   %s", name)
	}

	return nil
}

// 老的 sortMigrations 已迁移到 internal/migrationx.SortInPlace（共享唯一实现）。
// 排序规则由 migrationx 单测覆盖；行为必须与 scripts/_lib_migrations.sh 一致。

// splitSQL 按分号分割 SQL（正确处理注释行和空行）
func splitSQL(content string) []string {
	var result []string
	var current strings.Builder
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}
		current.WriteString(line)
		current.WriteString("\n")
		if strings.HasSuffix(trimmed, ";") {
			result = append(result, current.String())
			current.Reset()
		}
	}
	if current.Len() > 0 {
		result = append(result, current.String())
	}
	return result
}

// teardownTestDB 删除测试数据库
func teardownTestDB() {
	if testDB != nil {
		if rawDB, err := testDB.DB(); err == nil {
			rawDB.Close()
		}
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		dbHost, dbPort, dbUser, dbPass)
	adminDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err == nil {
		if rawDB, err := adminDB.DB(); err == nil {
			defer rawDB.Close()
			rawDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
		}
	}
	log.Printf("[integration] 数据库 %s 已清理", testDBName)
}

// ============================================================
// 路由搭建
// ============================================================

// TestServer 集成测试服务器包装
type TestServer struct {
	Router *gin.Engine
	DB     *gorm.DB
}

// NewTestServer 创建集成测试服务器
func NewTestServer() *TestServer {
	gin.SetMode(gin.TestMode)

	// 注入测试 DB 到全局 database 包（服务层通过 database.GetDB() 使用）
	database.SetDB(testDB)

	r := gin.New()
	r.Use(gin.Recovery())
	setupRoutes(r)

	return &TestServer{Router: r, DB: testDB}
}

// setupRoutes 与 cmd/main.go 一致的路由结构（去除 Redis/WebSocket/搜索等外部依赖）
func setupRoutes(r *gin.Engine) {
	authHandler := auth.NewAuthHandler()
	userHandler := user.NewUserHandler()
	postHandler := post.NewPostHandler()
	commentHandler := comment.NewCommentHandler()
	interactionHandler := interaction.NewInteractionHandler()
	followHandler := follow.NewFollowHandler()
	reportHandler := report.NewReportHandler()
	notificationHandler := notification.NewHandler()
	adminHandler := admin.NewHandler()

	// 公开路由
	public := r.Group("/api/v1")
	{
		public.POST("/auth/register/email", authHandler.RegisterByEmail)
		public.POST("/auth/login/email", authHandler.LoginByEmail)
		public.GET("/users/:id", userHandler.GetProfile)
		public.GET("/users/:id/posts", userHandler.GetUserPosts)
		public.GET("/users/:id/stats", userHandler.GetUserStats)
		public.GET("/posts/:id", postHandler.GetDetail)
		public.GET("/posts", postHandler.List)
		public.GET("/posts/:id/comments", commentHandler.List)
		public.GET("/users/:id/following", followHandler.ListFollowing)
		public.GET("/users/:id/followers", followHandler.ListFollowers)
	}

	// 需要认证的路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/auth/profile", authHandler.GetProfile)
		protected.PUT("/users/:id", userHandler.UpdateProfile)
		protected.POST("/users/avatar", userHandler.UploadAvatar)
		protected.POST("/posts", postHandler.Create)
		protected.PUT("/posts/:id", postHandler.Update)
		protected.DELETE("/posts/:id", postHandler.Delete)
		protected.GET("/feed", postHandler.Feed)
		protected.GET("/posts/following", postHandler.FollowingFeed)
		protected.POST("/posts/:id/comments", commentHandler.Create)
		protected.DELETE("/comments/:id", commentHandler.Delete)
		protected.POST("/posts/:id/like", interactionHandler.Like)
		protected.POST("/posts/:id/bookmark", interactionHandler.Bookmark)
		protected.GET("/users/:id/bookmarks", interactionHandler.ListBookmarks)
		protected.POST("/users/:id/follow", followHandler.Follow)
		protected.POST("/reports", reportHandler.Create)
		protected.GET("/notifications", notificationHandler.List)
		protected.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
		protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
		protected.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
	}

	// 管理员路由
	adminGroup := r.Group("/api/v1")
	adminGroup.Use(middleware.JWTAuth(), middleware.AdminAuth())
	{
		adminGroup.GET("/admin/posts", adminHandler.ListPosts)
		adminGroup.DELETE("/admin/posts/:id", adminHandler.DeletePost)
		adminGroup.GET("/admin/users", adminHandler.ListUsers)
		adminGroup.PUT("/admin/users/:id/status", adminHandler.UpdateUserStatus)
		adminGroup.GET("/admin/stats", adminHandler.GetStats)
	}

	// 健康检查（集成测试冒烟端点）
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok", "mode": "integration-test"})
	})

	// 调试端点：直接测试 FollowingFeed 逻辑
	r.GET("/debug/following/:id", func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID := int64(0)
		fmt.Sscanf(userIDStr, "%d", &userID)

		// 1. 直接查 follows 表
		var followIDs []int64
		database.GetDB().Table("follows").Where("follower_id = ?", userID).Pluck("following_id", &followIDs)

		// 2. 用原生 SQL 查 posts
		var rows []map[string]interface{}
		database.GetDB().Raw("SELECT id, user_id, status, title, deleted_at FROM posts WHERE user_id IN ?", followIDs).Scan(&rows)

		response.Success(c, gin.H{
			"followIDs":  followIDs,
			"posts":      rows,
			"dbAddr":     fmt.Sprintf("%p", database.GetDB()),
		})
	})
}

// CleanupTestData 清空所有测试数据（在需要隔离的测试中手动调用）
func CleanupTestData() {
	if testDB == nil {
		return
	}
	tables := []string{
		"notifications", "post_tags", "reports", "comments",
		"bookmarks", "likes", "follows", "posts", "tags", "users",
	}
	for _, t := range tables {
		testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", t))
	}
}
