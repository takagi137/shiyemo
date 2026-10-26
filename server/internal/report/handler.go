package report

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// ReportHandler 举报处理器
type ReportHandler struct {
	service *ReportService
}

// NewReportHandler 创建举报处理器
func NewReportHandler() *ReportHandler {
	return &ReportHandler{
		service: NewReportService(),
	}
}

// Create 创建举报
func (h *ReportHandler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	// 从 context 获取 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	if err := h.service.Create(userID.(int64), req); err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, nil)
}

// List 获取举报列表（管理员）
func (h *ReportHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	reports, total, err := h.service.List(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.SuccessPage(c, reports, page, size, total)
}

// Handle 处理举报（管理员）
func (h *ReportHandler) Handle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "无效的举报ID")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=2 3"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	if err := h.service.Handle(id, req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.Success(c, nil)
}
