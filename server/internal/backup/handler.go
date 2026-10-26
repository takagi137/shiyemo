package backup

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/common/response"
)

// Handler 备份处理器
type Handler struct {
	service *BackupService
}

// NewHandler 创建备份处理器
func NewHandler(service *BackupService) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	backup := r.Group("/backup")
	{
		backup.POST("/create", h.CreateBackup)
		backup.GET("/list", h.ListBackups)
		backup.POST("/restore", h.RestoreBackup)
		backup.DELETE("/:name", h.DeleteBackup)
		backup.POST("/cleanup", h.CleanupOldBackups)
	}
}

// CreateBackup 创建备份
func (h *Handler) CreateBackup(c *gin.Context) {
	backupFile, err := h.service.Backup()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "备份失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"file": backupFile,
		"time": time.Now().Format(time.RFC3339),
	})
}

// ListBackups 列出备份
func (h *Handler) ListBackups(c *gin.Context) {
	backups, err := h.service.ListBackups()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取备份列表失败: "+err.Error())
		return
	}

	response.Success(c, backups)
}

// RestoreBackup 恢复备份
func (h *Handler) RestoreBackup(c *gin.Context) {
	var req struct {
		File string `json:"file" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}

	if err := h.service.Restore(req.File); err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "恢复失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "恢复成功",
		"file":    req.File,
	})
}

// DeleteBackup 删除备份
func (h *Handler) DeleteBackup(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.Error(c, http.StatusBadRequest, 40001, "备份文件名不能为空")
		return
	}

	// 构建完整路径
	backupFile := h.service.backupDir + "/" + name

	if err := h.service.DeleteBackup(backupFile); err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "删除失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "删除成功",
		"file":    name,
	})
}

// CleanupOldBackups 清理旧备份
func (h *Handler) CleanupOldBackups(c *gin.Context) {
	// 获取保留天数参数
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		days = 7
	}

	maxAge := time.Duration(days) * 24 * time.Hour
	deleted, err := h.service.CleanupOldBackups(maxAge)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "清理失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"deleted": deleted,
		"days":    days,
	})
}