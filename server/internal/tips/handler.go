package tips

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// TipsHandler 贴士处理器
type TipsHandler struct {
	service *TipsService
}

// NewTipsHandler 创建贴士处理器
func NewTipsHandler() *TipsHandler {
	return &TipsHandler{
		service: NewTipsService(),
	}
}

// GetRandom 获取随机贴士
func (h *TipsHandler) GetRandom(c *gin.Context) {
	tip, err := h.service.GetRandom()
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, tip)
}

// List 获取贴士列表（管理员）
func (h *TipsHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	tips, total, err := h.service.GetList(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrDB.Code, bizerr.ErrDB.Message)
		return
	}

	response.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     tips,
	})
}

// CreateRequest 创建贴士请求
type CreateRequest struct {
	Content  string `json:"content" binding:"required"`
	Category string `json:"category"`
}

// Create 创建贴士（管理员）
func (h *TipsHandler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	tip, err := h.service.Create(req.Content, req.Category)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrDB.Code, bizerr.ErrDB.Message)
		return
	}

	response.Success(c, tip)
}

// UpdateRequest 更新贴士请求
type UpdateRequest struct {
	Content  string `json:"content" binding:"required"`
	Category string `json:"category"`
}

// Update 更新贴士（管理员）
func (h *TipsHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "无效的贴士ID")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	if err := h.service.Update(id, req.Content, req.Category); err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrDB.Code, bizerr.ErrDB.Message)
		return
	}

	response.Success(c, nil)
}

// Delete 删除贴士（管理员）
func (h *TipsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "无效的贴士ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrDB.Code, bizerr.ErrDB.Message)
		return
	}

	response.Success(c, nil)
}

// GetDailyContent 获取每日内容（风景图 + 文案）
// GET /api/v1/daily
func (h *TipsHandler) GetDailyContent(c *gin.Context) {
	data, err := h.service.GetDailyContent()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "获取每日内容失败")
		return
	}

	response.Success(c, data)
}

// RandomQuote 随机获取文案
// GET /api/v1/daily/quote
func (h *TipsHandler) RandomQuote(c *gin.Context) {
	quote, err := h.service.RandomQuote()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "获取文案失败")
		return
	}

	response.Success(c, quote)
}
