package feedback

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/common/response"
)

// Handler 反馈处理器
type Handler struct {
	service *Service
}

// NewHandler 创建处理器
func NewHandler() *Handler {
	return &Handler{service: NewService()}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	// 用户端接口
	feedback := r.Group("/feedback")
	{
		feedback.POST("", h.Create)
		feedback.GET("/mine", h.ListMine)
	}

	// 管理端接口
	admin := r.Group("/admin/feedback")
	{
		admin.GET("", h.List)
		admin.PUT("/:id/status", h.UpdateStatus)
	}
}

// Create 创建反馈
func (h *Handler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}

	var req struct {
		Type    string `json:"type"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Contact string `json:"contact"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}

	if err := h.service.Create(userID, req.Type, req.Title, req.Content, req.Contact); err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "提交反馈失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// ListMine 获取我的反馈列表
func (h *Handler) ListMine(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	feedbacks, total, err := h.service.ListByUser(userID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取反馈列表失败")
		return
	}

	response.Success(c, gin.H{
		"items":     feedbacks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// List 获取反馈列表（管理员）
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	feedbackType := c.Query("type")

	var status *int
	if s := c.Query("status"); s != "" {
		v, err := strconv.Atoi(s)
		if err == nil {
			status = &v
		}
	}

	feedbacks, total, err := h.service.List(page, pageSize, status, feedbackType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取反馈列表失败")
		return
	}

	response.Success(c, gin.H{
		"items":     feedbacks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateStatus 更新反馈状态（管理员）
func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的反馈ID")
		return
	}

	var req struct {
		Status int8   `json:"status" binding:"required"`
		Reply  string `json:"reply"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}

	if err := h.service.UpdateStatus(id, req.Status, req.Reply); err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "更新状态失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}
