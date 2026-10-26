package interaction

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// InteractionHandler 互动处理器
type InteractionHandler struct {
	service *InteractionService
}

// NewInteractionHandler 创建互动处理器
func NewInteractionHandler() *InteractionHandler {
	return &InteractionHandler{
		service: NewInteractionService(),
	}
}

// Like 点赞/取消点赞
func (h *InteractionHandler) Like(c *gin.Context) {
	postIDStr := c.Param("id")
	if postIDStr == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "帖子ID不能为空")
		return
	}

	postID, err := strconv.ParseInt(postIDStr, 10, 64)
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

	liked, likeCount, err := h.service.Like(userID.(int64), postID)
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, gin.H{"liked": liked, "like_count": likeCount})
}

// State 获取当前用户对某帖子的点赞/收藏状态
// GET /api/v1/posts/:id/interaction
func (h *InteractionHandler) State(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "无效的帖子ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	liked, bookmarked := h.service.GetUserPostState(userID.(int64), postID)
	response.Success(c, gin.H{"liked": liked, "bookmarked": bookmarked})
}

// Bookmark 收藏/取消收藏
func (h *InteractionHandler) Bookmark(c *gin.Context) {
	postIDStr := c.Param("id")
	if postIDStr == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "帖子ID不能为空")
		return
	}

	postID, err := strconv.ParseInt(postIDStr, 10, 64)
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

	bookmarked, err := h.service.Bookmark(userID.(int64), postID)
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, gin.H{"bookmarked": bookmarked})
}

// ListBookmarks 获取用户收藏列表
func (h *InteractionHandler) ListBookmarks(c *gin.Context) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "用户ID不能为空")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "无效的用户ID")
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

	bookmarks, total, err := h.service.ListBookmarks(userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.SuccessPage(c, bookmarks, page, size, total)
}
