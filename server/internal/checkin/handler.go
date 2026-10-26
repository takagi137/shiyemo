package checkin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// CheckinHandler 打卡处理器
type CheckinHandler struct {
	service *CheckinService
}

// NewCheckinHandler 创建打卡处理器
func NewCheckinHandler() *CheckinHandler {
	return &CheckinHandler{
		service: NewCheckinService(),
	}
}

// Create 创建打卡
func (h *CheckinHandler) Create(c *gin.Context) {
	// 从 JWT 中获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, "未登录")
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	checkin, err := h.service.Create(userID.(int64), &req)
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, checkin)
}

// GetToday 获取今天的打卡记录
func (h *CheckinHandler) GetToday(c *gin.Context) {
	// 从 JWT 中获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, "未登录")
		return
	}

	checkin, err := h.service.GetToday(userID.(int64))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrDB.Code, bizerr.ErrDB.Message)
		return
	}

	// 返回结果（可能为 nil，表示今天未打卡）
	response.Success(c, checkin)
}

// GetHistory 获取打卡历史
func (h *CheckinHandler) GetHistory(c *gin.Context) {
	// 从 JWT 中获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	checkins, total, err := h.service.GetHistory(userID.(int64), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrDB.Code, bizerr.ErrDB.Message)
		return
	}

	response.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     checkins,
	})
}

// GetMoodCurve 获取心情曲线
func (h *CheckinHandler) GetMoodCurve(c *gin.Context) {
	// 从 JWT 中获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, "未登录")
		return
	}

	items, err := h.service.GetMoodCurve(userID.(int64))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrDB.Code, bizerr.ErrDB.Message)
		return
	}

	response.Success(c, items)
}

// GetPublicTimeline 获取公开打卡时间线
func (h *CheckinHandler) GetPublicTimeline(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	checkins, total, err := h.service.GetPublicTimeline(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrDB.Code, bizerr.ErrDB.Message)
		return
	}

	response.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     checkins,
	})
}
