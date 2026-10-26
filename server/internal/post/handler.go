package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// PostHandler 帖子处理器
type PostHandler struct {
	service *PostService
}

// NewPostHandler 创建帖子处理器
func NewPostHandler() *PostHandler {
	return &PostHandler{
		service: NewPostService(),
	}
}

// Service 获取内部 PostService（用于外部注入依赖）
func (h *PostHandler) Service() *PostService {
	return h.service
}

// Create 创建帖子
func (h *PostHandler) Create(c *gin.Context) {
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

	post, err := h.service.Create(userID.(int64), req)
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, post)
}

// GetDetail 获取帖子详情
func (h *PostHandler) GetDetail(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "帖子ID不能为空")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "无效的帖子ID")
		return
	}

	detail, err := h.service.GetDetail(id)
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, detail)
}

// List 帖子列表
func (h *PostHandler) List(c *gin.Context) {
	postType := c.Query("type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	posts, total, err := h.service.List(postType, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.SuccessPage(c, posts, page, size, total)
}

// Feed 信息流
func (h *PostHandler) Feed(c *gin.Context) {
	postType := c.Query("type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	posts, total, err := h.service.List(postType, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.SuccessPage(c, posts, page, size, total)
}

// FollowingFeed 关注动态
func (h *PostHandler) FollowingFeed(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	// 从 context 获取 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	posts, total, err := h.service.FollowingFeed(userID.(int64), page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.SuccessPage(c, posts, page, size, total)
}

// CohortFeed 同温层帖子列表
func (h *PostHandler) CohortFeed(c *gin.Context) {
	postType := c.Query("type")
	city := c.Query("city")
	industry := c.Query("industry")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	posts, total, err := h.service.CohortFeed(postType, city, industry, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.SuccessPage(c, posts, page, size, total)
}

// Update 更新帖子
func (h *PostHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "帖子ID不能为空")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "无效的帖子ID")
		return
	}

	// 从 context 获取 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	if err := h.service.Update(userID.(int64), id, req); err != nil {
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

// Delete 删除帖子
func (h *PostHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "帖子ID不能为空")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "无效的帖子ID")
		return
	}

	// 从 context 获取 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	if err := h.service.Delete(userID.(int64), id); err != nil {
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
