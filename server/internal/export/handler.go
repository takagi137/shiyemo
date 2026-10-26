package export

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// ExportHandler 导出处理器
type ExportHandler struct {
	service *ExportService
}

// NewExportHandler 创建导出处理器
func NewExportHandler() *ExportHandler {
	return &ExportHandler{
		service: NewExportService(),
	}
}

// ExportUsers 导出用户数据
func (h *ExportHandler) ExportUsers(c *gin.Context) {
	data, err := h.service.ExportUsers()
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	// 设置响应头
	filename := fmt.Sprintf("users_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}

// ExportPosts 导出帖子数据
func (h *ExportHandler) ExportPosts(c *gin.Context) {
	postType := c.Query("type")

	data, err := h.service.ExportPosts(postType)
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	// 设置响应头
	filename := fmt.Sprintf("posts_%s_%s.csv", postType, time.Now().Format("20060102_150405"))
	if postType == "" {
		filename = fmt.Sprintf("posts_all_%s.csv", time.Now().Format("20060102_150405"))
	}
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}

// GetStats 获取统计数据
func (h *ExportHandler) GetStats(c *gin.Context) {
	stats, err := h.service.ExportStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.Success(c, stats)
}
