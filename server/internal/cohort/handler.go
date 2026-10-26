package cohort

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// CohortHandler 同温层处理器
type CohortHandler struct {
	service *CohortService
}

// NewCohortHandler 创建同温层处理器
func NewCohortHandler() *CohortHandler {
	return &CohortHandler{
		service: NewCohortService(),
	}
}

// GetMyStats 获取当前用户的同温层统计
func (h *CohortHandler) GetMyStats(c *gin.Context) {
	// 获取当前用户 ID
	currentUserID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	stats, err := h.service.GetMyStats(currentUserID.(int64))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "获取同温层统计失败")
		return
	}

	response.Success(c, stats)
}

// GetStats 获取指定同温层统计
func (h *CohortHandler) GetStats(c *gin.Context) {
	city := c.Query("city")
	industry := c.Query("industry")
	timeRange := CohortTimeRange(c.DefaultQuery("time_range", "all"))

	stats, err := h.service.GetStats(city, industry, timeRange)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "获取同温层统计失败")
		return
	}

	response.Success(c, stats)
}

// ListPeers 获取同温层用户列表
func (h *CohortHandler) ListPeers(c *gin.Context) {
	city := c.Query("city")
	industry := c.Query("industry")
	timeRange := CohortTimeRange(c.DefaultQuery("time_range", "all"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filter := CohortFilter{
		City:      city,
		Industry:  industry,
		TimeRange: timeRange,
		Page:      page,
		PageSize:  pageSize,
	}

	peers, total, err := h.service.ListPeers(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "获取同温层用户列表失败")
		return
	}

	response.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     peers,
	})
}

// GetCities 获取所有城市列表
func (h *CohortHandler) GetCities(c *gin.Context) {
	cities, err := h.service.GetCities()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "获取城市列表失败")
		return
	}

	response.Success(c, cities)
}

// GetIndustries 获取所有行业列表
func (h *CohortHandler) GetIndustries(c *gin.Context) {
	industries, err := h.service.GetIndustries()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "获取行业列表失败")
		return
	}

	response.Success(c, industries)
}
