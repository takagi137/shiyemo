package follow

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// FollowHandler 关注处理器
type FollowHandler struct {
	service *FollowService
}

// NewFollowHandler 创建关注处理器
func NewFollowHandler() *FollowHandler {
	return &FollowHandler{
		service: NewFollowService(),
	}
}

// Follow 关注/取消关注
func (h *FollowHandler) Follow(c *gin.Context) {
	followingIDStr := c.Param("id")
	if followingIDStr == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "用户ID不能为空")
		return
	}

	followingID, err := strconv.ParseInt(followingIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "无效的用户ID")
		return
	}

	// 从 context 获取 user_id
	followerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	following, err := h.service.Follow(followerID.(int64), followingID)
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, gin.H{"following": following})
}

// ListFollowing 获取关注列表
func (h *FollowHandler) ListFollowing(c *gin.Context) {
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

	following, total, err := h.service.ListFollowing(userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.SuccessPage(c, following, page, size, total)
}

// ListFollowers 获取粉丝列表
func (h *FollowHandler) ListFollowers(c *gin.Context) {
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

	followers, total, err := h.service.ListFollowers(userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.SuccessPage(c, followers, page, size, total)
}
