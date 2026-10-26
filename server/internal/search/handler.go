package search

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// SearchHandler 搜索处理器
type SearchHandler struct {
	service *SearchService
}

// NewSearchHandler 创建搜索处理器
func NewSearchHandler() *SearchHandler {
	return &SearchHandler{
		service: NewSearchService(),
	}
}

// SearchPosts 搜索帖子
func (h *SearchHandler) SearchPosts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "搜索关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	posts, total, err := h.service.SearchPosts(query, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.SuccessPage(c, posts, page, size, total)
}

// SearchUsers 搜索用户
func (h *SearchHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "搜索关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	users, total, err := h.service.SearchUsers(query, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.SuccessPage(c, users, page, size, total)
}

// SyncIndexes 同步索引
func (h *SearchHandler) SyncIndexes(c *gin.Context) {
	if err := h.service.SyncAllPosts(); err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "同步帖子失败")
		return
	}

	if err := h.service.SyncAllUsers(); err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "同步用户失败")
		return
	}

	response.Success(c, gin.H{"message": "索引同步完成"})
}

// GetSuggestions 获取搜索建议
func (h *SearchHandler) GetSuggestions(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "搜索关键词不能为空")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	suggestions, err := h.service.GetSuggestions(query, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "获取搜索建议失败")
		return
	}

	response.Success(c, gin.H{"suggestions": suggestions})
}
