package ambientsound

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/common/response"
)

// Handler 环境音处理器
type Handler struct {
	svc *Service
}

// NewHandler 创建处理器
func NewHandler() *Handler {
	return &Handler{svc: NewService()}
}

// ListPublic 获取启用音效列表（公开）
// GET /api/v1/sounds
func (h *Handler) ListPublic(c *gin.Context) {
	sounds, err := h.svc.ListActive()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "获取音效列表失败")
		return
	}
	// 保证是数组而非 null
	if sounds == nil {
		sounds = []AmbientSound{}
	}
	response.Success(c, gin.H{"items": sounds})
}

// List 获取音效列表（管理员）
// GET /api/v1/admin/sounds
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	sounds, total, err := h.svc.List(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "获取列表失败")
		return
	}

	response.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     sounds,
	})
}

// CreateRequest 创建请求
type CreateRequest struct {
	Key       string `json:"key" binding:"required"`
	Name      string `json:"name"`
	Emoji     string `json:"emoji"`
	URL       string `json:"url" binding:"required"`
	Credit    string `json:"credit"`
	SortOrder int    `json:"sort_order"`
}

// Create 创建音效（管理员）
// POST /api/v1/admin/sounds
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	sound, err := h.svc.Create(req.Key, req.Name, req.Emoji, req.URL, req.Credit, req.SortOrder)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, sound)
}

// UpdateRequest 更新请求
type UpdateRequest struct {
	Key       string `json:"key" binding:"required"`
	Name      string `json:"name"`
	Emoji     string `json:"emoji"`
	URL       string `json:"url" binding:"required"`
	Credit    string `json:"credit"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
}

// Update 更新音效（管理员）
// PUT /api/v1/admin/sounds/:id
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的ID")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.Update(id, req.Key, req.Name, req.Emoji, req.URL, req.Credit, req.SortOrder, req.IsActive); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除音效（管理员）
// DELETE /api/v1/admin/sounds/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效的ID")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "删除失败")
		return
	}

	response.Success(c, nil)
}

// ReorderRequest 排序请求
type ReorderRequest struct {
	IDs []int64 `json:"ids" binding:"required"`
}

// Reorder 批量调整排序（管理员）
// PUT /api/v1/admin/sounds/reorder
func (h *Handler) Reorder(c *gin.Context) {
	var req ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.Reorder(req.IDs); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}
