package reading

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// ReadingHandler 阅读处理器
type ReadingHandler struct {
	service *ReadingService
}

// NewReadingHandler 创建阅读处理器
func NewReadingHandler() *ReadingHandler {
	return &ReadingHandler{
		service: NewReadingService(),
	}
}

// GetToday 获取今日阅读
// GET /api/v1/readings/today
func (h *ReadingHandler) GetToday(c *gin.Context) {
	reading, err := h.service.GetToday()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "获取今日阅读失败")
		return
	}

	response.Success(c, reading)
}

// GetRandom 获取随机阅读
// GET /api/v1/readings/random?exclude=1,2,3
func (h *ReadingHandler) GetRandom(c *gin.Context) {
	var excludeIDs []int64
	excludeStr := c.Query("exclude")
	if excludeStr != "" {
		for _, idStr := range strings.Split(excludeStr, ",") {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err == nil {
				excludeIDs = append(excludeIDs, id)
			}
		}
	}

	reading, err := h.service.GetRandom(excludeIDs)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "获取随机阅读失败")
		return
	}

	response.Success(c, reading)
}
