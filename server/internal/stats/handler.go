package stats

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// StatsHandler 统计处理器
type StatsHandler struct {
	service StatsServiceInterface
}

// NewStatsHandler 创建统计处理器
func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		service: NewStatsService(),
	}
}

// NewStatsHandlerWithService 使用注入的服务创建处理器（测试用）
func NewStatsHandlerWithService(service StatsServiceInterface) *StatsHandler {
	return &StatsHandler{service: service}
}

// GetTrends 获取趋势数据
func (h *StatsHandler) GetTrends(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 {
		days = 30
	}

	trends, err := h.service.GetTrends(days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "获取趋势数据失败")
		return
	}

	response.Success(c, trends)
}
