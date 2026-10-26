package health

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/common/cache"
	"github.com/shiyemo/server/internal/common/database"
)

// Handler 健康检查处理器
type Handler struct {
	startTime time.Time
}

// NewHandler 创建健康检查处理器
func NewHandler() *Handler {
	return &Handler{
		startTime: time.Now(),
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/health", h.HealthCheck)
	r.GET("/health/detailed", h.DetailedHealthCheck)
}

// HealthCheck 基础健康检查
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
	})
}

// DetailedHealthCheck 详细健康检查
func (h *Handler) DetailedHealthCheck(c *gin.Context) {
	checks := make(map[string]interface{})
	overallStatus := "ok"

	// 检查数据库连接
	dbStatus := "ok"
	dbErr := h.checkDatabase()
	if dbErr != nil {
		dbStatus = "error"
		overallStatus = "degraded"
	}
	checks["database"] = map[string]interface{}{
		"status": dbStatus,
		"error":  dbErr,
	}

	// 检查 Redis 连接
	redisStatus := "ok"
	redisErr := h.checkRedis()
	if redisErr != nil {
		redisStatus = "error"
		overallStatus = "degraded"
	}
	checks["redis"] = map[string]interface{}{
		"status": redisStatus,
		"error":  redisErr,
	}

	// 系统信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	checks["system"] = map[string]interface{}{
		"goroutines":    runtime.NumGoroutine(),
		"heap_alloc_mb": float64(m.HeapAlloc) / 1024 / 1024,
		"sys_mb":        float64(m.Sys) / 1024 / 1024,
		"gc_cycles":     m.NumGC,
	}

	// 运行时间
	uptime := time.Since(h.startTime)
	checks["uptime"] = map[string]interface{}{
		"seconds": int(uptime.Seconds()),
		"human":   uptime.String(),
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    overallStatus,
		"timestamp": time.Now().Unix(),
		"checks":    checks,
	})
}

// checkDatabase 检查数据库连接
func (h *Handler) checkDatabase() error {
	db := database.GetDB()
	if db == nil {
		return ErrDatabaseNotInitialized
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// checkRedis 检查 Redis 连接
func (h *Handler) checkRedis() error {
	client := cache.GetClient()
	if client == nil {
		return ErrRedisNotInitialized
	}

	return client.Ping(nil).Err()
}

// 错误定义
var (
	ErrDatabaseNotInitialized = &HealthError{Message: "database not initialized"}
	ErrRedisNotInitialized    = &HealthError{Message: "redis not initialized"}
)

// HealthError 健康检查错误
type HealthError struct {
	Message string
}

func (e *HealthError) Error() string {
	return e.Message
}