package moderation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/common/response"
)

// Handler 审核处理器
type Handler struct {
	service *Service
}

// NewHandler 创建审核处理器
func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin/moderation")
	{
		admin.GET("/words", h.ListWords)
		admin.POST("/words", h.AddWord)
		admin.DELETE("/words/:id", h.DeleteWord)
	}
}

// ListWordsRequest 敏感词列表请求
type ListWordsRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

// ListWords 获取敏感词列表
func (h *Handler) ListWords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	words, total, err := h.service.ListWords(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取敏感词列表失败")
		return
	}

	response.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     words,
	})
}

// AddWordRequest 添加敏感词请求
type AddWordRequest struct {
	Word string `json:"word" binding:"required,max=50"`
}

// AddWord 添加敏感词
func (h *Handler) AddWord(c *gin.Context) {
	var req AddWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}

	word, err := h.service.AddWord(req.Word)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "添加敏感词失败")
		return
	}

	response.Success(c, word)
}

// DeleteWord 删除敏感词
func (h *Handler) DeleteWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的敏感词ID")
		return
	}

	if err := h.service.DeleteWord(id); err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "删除敏感词失败")
		return
	}

	response.Success(c, nil)
}
