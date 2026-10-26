package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/admin"
	"github.com/shiyemo/server/internal/ambientsound"
	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/backup"
	"github.com/shiyemo/server/internal/checkin"
	"github.com/shiyemo/server/internal/comment"
	"github.com/shiyemo/server/internal/cohort"
	"github.com/shiyemo/server/internal/common/cache"
	"github.com/shiyemo/server/internal/common/config"
	"github.com/shiyemo/server/internal/common/database"
	"github.com/shiyemo/server/internal/common/middleware"
	"github.com/shiyemo/server/internal/common/storage"
	"github.com/shiyemo/server/internal/dailyimage"
	"github.com/shiyemo/server/internal/export"
	"github.com/shiyemo/server/internal/feedback"
	"github.com/shiyemo/server/internal/follow"
	"github.com/shiyemo/server/internal/journal"
	"github.com/shiyemo/server/internal/reading"
	"github.com/shiyemo/server/internal/health"
	"github.com/shiyemo/server/internal/interaction"
	logService "github.com/shiyemo/server/internal/log"
	"github.com/shiyemo/server/internal/moderation"
	"github.com/shiyemo/server/internal/notification"
	"github.com/shiyemo/server/internal/post"
	"github.com/shiyemo/server/internal/report"
	"github.com/shiyemo/server/internal/search"
	"github.com/shiyemo/server/internal/setting"
	"github.com/shiyemo/server/internal/stats"
	"github.com/shiyemo/server/internal/tag"
	"github.com/shiyemo/server/internal/tips"
	"github.com/shiyemo/server/internal/upload"
	"github.com/shiyemo/server/internal/user"
	"github.com/shiyemo/server/internal/verify"
	"github.com/shiyemo/server/internal/websocket"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	database.Init(&cfg.DB)

	// 初始化 Redis
	if err := cache.Init(cfg.Redis); err != nil {
		log.Printf("Warning: Redis init failed: %v (cache disabled)", err)
	} else {
		defer cache.Close()
		log.Println("Redis connected successfully")
	}

	// 初始化搜索服务
	searchService := search.NewSearchService()
	searchService.InitIndexes()

	// 初始化 WebSocket
	websocket.Init()
	log.Println("WebSocket hub initialized")

	// 初始化存储
	store, err := storage.NewStorage(storage.Config{
		Type: cfg.Storage.Type,
		Local: storage.LocalConfig{
			UploadDir: cfg.Storage.UploadDir,
			BaseURL:   cfg.Storage.BaseURL,
		},
		MinIO: storage.MinIOConfig{
			Endpoint:  cfg.Storage.MinIO.Endpoint,
			AccessKey: cfg.Storage.MinIO.AccessKey,
			SecretKey: cfg.Storage.MinIO.SecretKey,
			Bucket:    cfg.Storage.MinIO.Bucket,
			UseSSL:    cfg.Storage.MinIO.UseSSL,
			BaseURL:   cfg.Storage.MinIO.BaseURL,
		},
	})
	if err != nil {
		log.Fatalf("Failed to init storage: %v", err)
	}
	log.Printf("Storage type: %s", cfg.Storage.Type)

	// 初始化日志服务
	logLevel := logService.INFO
	switch strings.ToLower(cfg.Log.Level) {
	case "debug":
		logLevel = logService.DEBUG
	case "warn":
		logLevel = logService.WARN
	case "error":
		logLevel = logService.ERROR
	}
	appLogService, err := logService.NewLogService(cfg.Log.Dir, logLevel)
	if err != nil {
		log.Fatalf("Failed to init log service: %v", err)
	}
	defer appLogService.Close()
	log.Printf("Log service initialized, dir: %s, level: %s", cfg.Log.Dir, cfg.Log.Level)

	// 初始化备份服务
	backupService, err := backup.NewBackupService(
		cfg.Backup.Dir,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.DBName,
	)
	if err != nil {
		log.Printf("Warning: Backup service init failed: %v", err)
	} else {
		log.Printf("Backup service initialized, dir: %s", cfg.Backup.Dir)
	}

	// 创建处理器
	authHandler := auth.NewAuthHandler()
	userHandler := user.NewUserHandler()
	postHandler := post.NewPostHandler()
	postHandler.Service().SetIndexer(searchService)
	commentHandler := comment.NewCommentHandler()
	interactionHandler := interaction.NewInteractionHandler()
	followHandler := follow.NewFollowHandler()
	searchHandler := search.NewSearchHandler()
	notificationHandler := notification.NewHandler()
	adminHandler := admin.NewHandler()
	settingHandler := setting.NewHandler()
	moderationHandler := moderation.NewHandler()
	uploadService := upload.NewUploadServiceWithStorage(store)
	uploadHandler := upload.NewUploadHandlerWithService(uploadService)
	wsHandler := websocket.NewWebSocketHandler()
	exportHandler := export.NewExportHandler()
	reportHandler := report.NewReportHandler()
	healthHandler := health.NewHandler()
	verifyHandler := verify.NewVerifyHandler()
	tipsHandler := tips.NewTipsHandler()
	checkinHandler := checkin.NewCheckinHandler()
	cohortHandler := cohort.NewCohortHandler()
	statsHandler := stats.NewStatsHandler()
	feedbackHandler := feedback.NewHandler()
	dailyImageHandler := dailyimage.NewHandler()
	ambientSoundHandler := ambientsound.NewHandler()
	journalHandler := journal.NewJournalHandler()
	readingHandler := reading.NewReadingHandler()

	// 启动时执行一次统计聚合（后台）
	go func() {
		svc := stats.NewStatsService()
		if err := svc.RunDailyAggregation(); err != nil {
			log.Printf("Initial stats aggregation failed: %v", err)
		}
	}()

	var backupHandler *backup.Handler
	if backupService != nil {
		backupHandler = backup.NewHandler(backupService)
	}

	r := gin.Default()

	// 全局中间件
	r.Use(middleware.Cors())
	r.Use(logService.Middleware(appLogService))

	// 静态文件服务（头像）
	r.Static("/uploads", "./uploads")

	// 健康检查
	healthHandler.RegisterRoutes(r.Group(""))

	// 公开路由（无需认证）
	public := r.Group("/api/v1")
	{
		// 认证
		public.POST("/auth/register/email", authHandler.RegisterByEmail)
		public.POST("/auth/register/phone", authHandler.RegisterByPhone)
		public.POST("/auth/login/email", authHandler.LoginByEmail)
		public.POST("/auth/login/phone", authHandler.LoginByPhone)
		public.POST("/auth/sms/send", authHandler.SendSMS)

		// 邮箱验证
		public.GET("/auth/verify/:token", verifyHandler.VerifyEmail)

		// 密码重置
		public.POST("/auth/forgot-password", verifyHandler.SendResetEmail)
		public.POST("/auth/reset-password", verifyHandler.ResetPassword)

		// 生存小贴士
		public.GET("/tips/random", tipsHandler.GetRandom)
	public.GET("/daily", tipsHandler.GetDailyContent)
	public.GET("/daily/quote", tipsHandler.RandomQuote)
	public.GET("/daily/image", dailyImageHandler.RandomImage)

		// 静思空间环境音列表
		public.GET("/sounds", ambientSoundHandler.ListPublic)

		// 公开打卡时间线
		public.GET("/checkin/timeline", checkinHandler.GetPublicTimeline)

		// 公开的用户资料
		public.GET("/users/:id", userHandler.GetProfile)
		public.GET("/users/:id/posts", userHandler.GetUserPosts)
		public.GET("/users/:id/stats", userHandler.GetUserStats)

		// 公开的帖子详情和列表
		public.GET("/posts/:id", postHandler.GetDetail)
		public.GET("/posts", postHandler.List)

		// 公开的同温层统计
		public.GET("/cohort/public-stats", cohortHandler.GetStats)

		// 公开的趋势数据
		public.GET("/stats/trends", statsHandler.GetTrends)

		// 公开的评论列表
		public.GET("/posts/:id/comments", commentHandler.List)

		// 公开的关注/粉丝列表
		public.GET("/users/:id/following", followHandler.ListFollowing)
		public.GET("/users/:id/followers", followHandler.ListFollowers)

		// 搜索
		public.GET("/search/posts", searchHandler.SearchPosts)
		public.GET("/search/users", searchHandler.SearchUsers)
		public.GET("/search/suggestions", searchHandler.GetSuggestions)
	}

	// 注册标签路由
	tag.RegisterRoutes(r.Group("/api/v1"))

	// 需要认证的路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuth())
	{
		// 用户
		protected.GET("/auth/profile", authHandler.GetProfile)
		protected.PUT("/users/:id", userHandler.UpdateProfile)
		protected.POST("/users/avatar", userHandler.UploadAvatar)

		// 邮箱验证
		protected.POST("/auth/verify/send", verifyHandler.SendVerificationEmail)

		// 帖子（需要邮箱验证）
		protected.POST("/posts", middleware.EmailVerified(), postHandler.Create)
		protected.PUT("/posts/:id", postHandler.Update)
		protected.DELETE("/posts/:id", postHandler.Delete)

		// 信息流
		protected.GET("/feed", postHandler.Feed)
		protected.GET("/posts/following", postHandler.FollowingFeed)
		protected.GET("/posts/cohort", postHandler.CohortFeed)

		// 评论
		protected.POST("/posts/:id/comments", commentHandler.Create)
		protected.DELETE("/comments/:id", commentHandler.Delete)

		// 点赞/收藏
		protected.POST("/posts/:id/like", interactionHandler.Like)
		protected.POST("/posts/:id/bookmark", interactionHandler.Bookmark)
		protected.GET("/posts/:id/interaction", interactionHandler.State)
		protected.GET("/users/:id/bookmarks", interactionHandler.ListBookmarks)

		// 关注
		protected.POST("/users/:id/follow", followHandler.Follow)

		// 同温层
		protected.GET("/cohort/stats", cohortHandler.GetMyStats)
		protected.GET("/cohort/peers", cohortHandler.ListPeers)
		protected.GET("/cohort/cities", cohortHandler.GetCities)
		protected.GET("/cohort/industries", cohortHandler.GetIndustries)

		// 举报
		protected.POST("/reports", reportHandler.Create)

		// 索引同步（管理员）
		protected.POST("/search/sync", searchHandler.SyncIndexes)

		// 通知
		protected.GET("/notifications", notificationHandler.List)
		protected.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
		protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
		protected.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
		protected.POST("/notifications/batch/read", notificationHandler.BatchMarkAsRead)
		protected.POST("/notifications/batch/delete", notificationHandler.BatchDelete)

		// 反馈功能
		posted := protected.Group("/feedback")
		{
			posted.POST("", feedbackHandler.Create)
			posted.GET("/mine", feedbackHandler.ListMine)
		}

		// 管理员
		protected.GET("/admin/posts", middleware.AdminAuth(), adminHandler.ListPosts)
		protected.DELETE("/admin/posts/:id", middleware.AdminAuth(), adminHandler.DeletePost)
		protected.GET("/admin/users", middleware.AdminAuth(), adminHandler.ListUsers)
		protected.PUT("/admin/users/:id/status", middleware.AdminAuth(), adminHandler.UpdateUserStatus)
		protected.GET("/admin/stats", middleware.AdminAuth(), adminHandler.GetStats)

		// 内容审核管理
		protected.GET("/admin/moderation/words", middleware.AdminAuth(), moderationHandler.ListWords)
		protected.POST("/admin/moderation/words", middleware.AdminAuth(), moderationHandler.AddWord)
		protected.DELETE("/admin/moderation/words/:id", middleware.AdminAuth(), moderationHandler.DeleteWord)

		// 系统设置（SMTP）
		protected.GET("/admin/settings/smtp", middleware.AdminAuth(), settingHandler.GetSMTP)
		protected.PUT("/admin/settings/smtp", middleware.AdminAuth(), settingHandler.UpdateSMTP)
		protected.POST("/admin/settings/smtp/test", middleware.AdminAuth(), settingHandler.TestSMTP)

		// 打卡功能
		protected.POST("/checkin", checkinHandler.Create)
		protected.GET("/checkin/today", checkinHandler.GetToday)
		protected.GET("/checkin/history", checkinHandler.GetHistory)
		protected.GET("/checkin/mood-curve", checkinHandler.GetMoodCurve)

		// 贴士管理（管理员）
		protected.GET("/admin/tips", middleware.AdminAuth(), tipsHandler.List)
		protected.POST("/admin/tips", middleware.AdminAuth(), tipsHandler.Create)
		protected.PUT("/admin/tips/:id", middleware.AdminAuth(), tipsHandler.Update)
		protected.DELETE("/admin/tips/:id", middleware.AdminAuth(), tipsHandler.Delete)

		// 帖子审核队列
		protected.GET("/admin/posts/pending", middleware.AdminAuth(), adminHandler.ListPendingPosts)
		protected.POST("/admin/posts/:id/approve", middleware.AdminAuth(), adminHandler.ApprovePost)
		protected.POST("/admin/posts/:id/reject", middleware.AdminAuth(), adminHandler.RejectPost)

		// 评论管理
		protected.GET("/admin/comments", middleware.AdminAuth(), adminHandler.ListComments)
		protected.PUT("/admin/comments/:id/hide", middleware.AdminAuth(), adminHandler.HideComment)
		protected.PUT("/admin/comments/:id/restore", middleware.AdminAuth(), adminHandler.RestoreComment)

		// 举报管理
		protected.GET("/admin/reports", middleware.AdminAuth(), reportHandler.List)
		protected.PUT("/admin/reports/:id", middleware.AdminAuth(), reportHandler.Handle)

		// 帖子置顶/加精
		protected.PUT("/admin/posts/:id/top", middleware.AdminAuth(), adminHandler.SetPostTop)
		protected.PUT("/admin/posts/:id/featured", middleware.AdminAuth(), adminHandler.SetPostFeatured)

		// 风景图管理（管理员）
		protected.GET("/admin/images", middleware.AdminAuth(), dailyImageHandler.List)
		protected.POST("/admin/images", middleware.AdminAuth(), dailyImageHandler.Create)
		protected.PUT("/admin/images/:id", middleware.AdminAuth(), dailyImageHandler.Update)
		protected.DELETE("/admin/images/:id", middleware.AdminAuth(), dailyImageHandler.Delete)
		protected.PUT("/admin/images/reorder", middleware.AdminAuth(), dailyImageHandler.Reorder)

		// 环境音管理（管理员）
		protected.GET("/admin/sounds", middleware.AdminAuth(), ambientSoundHandler.List)
		protected.POST("/admin/sounds", middleware.AdminAuth(), ambientSoundHandler.Create)
		protected.PUT("/admin/sounds/:id", middleware.AdminAuth(), ambientSoundHandler.Update)
		protected.DELETE("/admin/sounds/:id", middleware.AdminAuth(), ambientSoundHandler.Delete)
		protected.PUT("/admin/sounds/reorder", middleware.AdminAuth(), ambientSoundHandler.Reorder)

		// 反馈管理（管理员）
		protected.GET("/admin/feedback", middleware.AdminAuth(), feedbackHandler.List)
		protected.PUT("/admin/feedback/:id/status", middleware.AdminAuth(), feedbackHandler.UpdateStatus)

		// 备份管理（管理员）
		if backupHandler != nil {
			backupHandler.RegisterRoutes(protected.Group(""))
		}

		// 文件上传
		protected.POST("/upload/image", uploadHandler.UploadImage)
		protected.POST("/upload/images", uploadHandler.UploadMultiple)

		// 静思空间 - 私密日记
		protected.POST("/journal", journalHandler.Create)
		protected.GET("/journal", journalHandler.List)
		protected.DELETE("/journal/:id", journalHandler.Delete)

		// 静思空间 - 今日阅读
		public.GET("/readings/today", readingHandler.GetToday)
		public.GET("/readings/random", readingHandler.GetRandom)

		// WebSocket
		protected.GET("/ws", wsHandler.HandleWebSocket)

		// 数据导出（管理员）
		protected.GET("/admin/export/users", middleware.AdminAuth(), exportHandler.ExportUsers)
		protected.GET("/admin/export/posts", middleware.AdminAuth(), exportHandler.ExportPosts)
		protected.GET("/admin/export/stats", middleware.AdminAuth(), exportHandler.GetStats)
	}

	log.Printf("Server starting on :%s\n", cfg.App.Port)
	if err := r.Run(fmt.Sprintf(":%s", cfg.App.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
