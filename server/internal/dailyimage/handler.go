package dailyimage

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/common/response"
)

// Handler 风景图处理器
type Handler struct {
	svc *Service
}

// NewHandler 创建处理器
func NewHandler() *Handler {
	return &Handler{svc: NewService()}
}

// RandomImage 随机获取一张图片（公开）
// GET /api/v1/daily/image
func (h *Handler) RandomImage(c *gin.Context) {
	img, err := h.svc.RandomImageFresh()
	if err != nil {
		// 如果数据库为空，返回默认图片
		response.Success(c, gin.H{
			"url":         "https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=1920&q=80",
			"description": "山谷湖泊",
			"credit":      "Photo by Baechtold on Unsplash",
		})
		return
	}
	response.Success(c, img)
}

// List 获取图片列表（管理员）
// GET /api/v1/admin/images
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	images, total, err := h.svc.List(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "获取列表失败")
		return
	}

	response.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     images,
	})
}

// CreateRequest 创建请求
type CreateRequest struct {
	URL         string `json:"url" binding:"required"`
	Description string `json:"description"`
	Credit      string `json:"credit"`
	SortOrder   int    `json:"sort_order"`
}

// Create 创建图片（管理员）
// POST /api/v1/admin/images
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	img, err := h.svc.Create(req.URL, req.Description, req.Credit, req.SortOrder)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, img)
}

// UpdateRequest 更新请求
type UpdateRequest struct {
	URL         string `json:"url" binding:"required"`
	Description string `json:"description"`
	Credit      string `json:"credit"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

// Update 更新图片（管理员）
// PUT /api/v1/admin/images/:id
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

	if err := h.svc.Update(id, req.URL, req.Description, req.Credit, req.SortOrder, req.IsActive); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除图片（管理员）
// DELETE /api/v1/admin/images/:id
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
// PUT /api/v1/admin/images/reorder
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
