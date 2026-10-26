package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// UserHandler 用户处理器
type UserHandler struct {
	service *UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		service: NewUserService(),
	}
}

// GetProfile 获取指定用户公开资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "用户ID不能为空")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "无效的用户ID")
		return
	}

	publicProfile, err := h.service.GetPublicProfile(id)
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, publicProfile)
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "用户ID不能为空")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "无效的用户ID")
		return
	}

	// 检查权限：只能更新自己的资料
	currentUserID, exists := c.Get("user_id")
	if !exists || currentUserID.(int64) != id {
		response.Error(c, http.StatusForbidden, 40302, "无权操作")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	if err := h.service.UpdateProfile(id, req); err != nil {
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

// UploadAvatar 上传头像
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	// 获取当前用户 ID
	currentUserID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	// 获取上传文件
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "请选择头像文件")
		return
	}
	defer file.Close()

	// 检查文件大小（2MB）
	if header.Size > 2*1024*1024 {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "文件大小超过限制")
		return
	}

	// 上传头像
	avatarURL, err := h.service.UploadAvatar(currentUserID.(int64), file, header.Filename)
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, gin.H{
		"avatar_url": avatarURL,
	})
}

// GetUserPosts 获取用户帖子列表
func (h *UserHandler) GetUserPosts(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, 40001, "用户ID不能为空")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的用户ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	posts, total, err := h.service.GetUserPosts(id, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取用户帖子失败")
		return
	}

	response.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     posts,
	})
}

// GetUserStats 获取用户统计数据
func (h *UserHandler) GetUserStats(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, 40001, "用户ID不能为空")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的用户ID")
		return
	}

	stats, err := h.service.GetUserStats(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}
