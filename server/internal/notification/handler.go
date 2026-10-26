package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/common/middleware"
	"github.com/shiyemo/server/internal/common/response"
)

// Handler 通知处理器
type Handler struct {
	service *Service
}

// NewHandler 创建通知处理器
func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	notifications := r.Group("/notifications")
	notifications.Use(middleware.JWTAuth())
	{
		notifications.GET("", h.List)
		notifications.GET("/unread-count", h.GetUnreadCount)
		notifications.PUT("/:id/read", h.MarkAsRead)
		notifications.PUT("/read-all", h.MarkAllAsRead)
		notifications.POST("/batch/read", h.BatchMarkAsRead)
		notifications.POST("/batch/delete", h.BatchDelete)
	}
}

// List 获取通知列表
func (h *Handler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 解析筛选参数
	nType := c.Query("type")
	unreadOnly := c.Query("unread") == "true"

	var notifications []Notification
	var total int64
	var err error

	if nType != "" || unreadOnly {
		// 按类型或未读筛选
		notifications, total, err = h.service.ListByType(userID, nType, unreadOnly, page, pageSize)
	} else {
		// 获取全部
		notifications, total, err = h.service.List(userID, page, pageSize)
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取通知列表失败")
		return
	}

	response.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     notifications,
	})
}

// GetUnreadCount 获取未读数量
func (h *Handler) GetUnreadCount(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}

	count, err := h.service.GetUnreadCount(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取未读数量失败")
		return
	}

	response.Success(c, gin.H{"count": count})
}

// MarkAsRead 标记单条已读
func (h *Handler) MarkAsRead(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的通知ID")
		return
	}

	err = h.service.MarkAsRead(id, userID)
	if err != nil {
		if err.Error() == "通知不存在" {
			response.Error(c, http.StatusNotFound, 40401, "通知不存在")
			return
		}
		if err.Error() == "无权操作该通知" {
			response.Error(c, http.StatusForbidden, 40301, "无权操作该通知")
			return
		}
		response.Error(c, http.StatusInternalServerError, 50001, "标记已读失败")
		return
	}

	response.Success(c, nil)
}

// MarkAllAsRead 标记全部已读
func (h *Handler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}

	err := h.service.MarkAllAsRead(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "标记全部已读失败")
		return
	}

	response.Success(c, nil)
}

// BatchMarkAsRead 批量标记已读
func (h *Handler) BatchMarkAsRead(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}

	var req struct {
		IDs []int64 `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}

	err := h.service.BatchMarkAsRead(req.IDs, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "批量标记已读失败")
		return
	}

	response.Success(c, nil)
}

// BatchDelete 批量删除
func (h *Handler) BatchDelete(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}

	var req struct {
		IDs []int64 `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}

	err := h.service.BatchDelete(req.IDs, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "批量删除失败")
		return
	}

	response.Success(c, nil)
}
