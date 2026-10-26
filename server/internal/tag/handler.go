package tag

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/middleware"
	"github.com/shiyemo/server/internal/common/response"
)

// RegisterRoutes 注册标签路由
func RegisterRoutes(r *gin.RouterGroup) {
	svc := NewTagService()
	skillSvc := NewSkillTagService()
	h := &TagHandler{svc: svc, skillSvc: skillSvc}

	tagGroup := r.Group("/tags")
	{
		tagGroup.GET("", h.GetHotTags)
		tagGroup.GET("/skill/search", h.SearchBySkill)
	}

	// 需要登录的路由
	postGroup := r.Group("/posts")
	postGroup.Use(middleware.JWTAuth())
	{
		postGroup.PUT("/:id/tags", h.UpdatePostTags)
	}

	// 用户技能标签路由
	userGroup := r.Group("/users")
	userGroup.Use(middleware.JWTAuth())
	{
		userGroup.GET("/:id/skills", h.GetUserSkills)
		userGroup.POST("/:id/skills", h.AddUserSkill)
		userGroup.DELETE("/:id/skills", h.RemoveUserSkill)
	}
}

// TagHandler 标签处理器
type TagHandler struct {
	svc     *TagService
	skillSvc *SkillTagService
}

// GetHotTags 获取热门标签
func (h *TagHandler) GetHotTags(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	tags, err := h.svc.GetHotTags(limit)
	if err != nil {
		if bizErr, ok := err.(*bizerr.BizError); ok {
			response.Error(c, http.StatusInternalServerError, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, 50000, "服务器内部错误")
		}
		return
	}

	response.Success(c, tags)
}

// UpdatePostTagsRequest 更新帖子标签请求
type UpdatePostTagsRequest struct {
	Tags []string `json:"tags" binding:"required"`
}

// UpdatePostTags 更新帖子标签
func (h *TagHandler) UpdatePostTags(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的帖子ID")
		return
	}

	var req UpdatePostTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}

	// TODO: 验证是否是帖子作者

	tags, svcErr := h.svc.UpdatePostTags(postID, req.Tags)
	if svcErr != nil {
		if bizErr, ok := svcErr.(*bizerr.BizError); ok {
			response.Error(c, http.StatusInternalServerError, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, 50000, "服务器内部错误")
		}
		return
	}

	response.Success(c, gin.H{"tags": tags})
}

// GetUserSkills 获取用户技能标签
func (h *TagHandler) GetUserSkills(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的用户ID")
		return
	}

	skills, svcErr := h.skillSvc.GetUserSkills(userID)
	if svcErr != nil {
		if bizErr, ok := svcErr.(*bizerr.BizError); ok {
			response.Error(c, http.StatusInternalServerError, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, 50000, "服务器内部错误")
		}
		return
	}

	response.Success(c, skills)
}

// AddUserSkillRequest 添加技能标签请求
type AddUserSkillRequest struct {
	Tag string `json:"tag" binding:"required,max=50"`
}

// AddUserSkill 添加用户技能标签
func (h *TagHandler) AddUserSkill(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的用户ID")
		return
	}

	// 验证只能修改自己的技能标签
	currentUserID, _ := c.Get("userID")
	if currentUserID.(int64) != userID {
		response.Error(c, http.StatusForbidden, 40302, "无权操作")
		return
	}

	var req AddUserSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}

	skill, svcErr := h.skillSvc.AddUserSkill(userID, req.Tag)
	if svcErr != nil {
		if bizErr, ok := svcErr.(*bizerr.BizError); ok {
			response.Error(c, http.StatusInternalServerError, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, 50000, "服务器内部错误")
		}
		return
	}

	response.Success(c, skill)
}

// RemoveUserSkillRequest 移除技能标签请求
type RemoveUserSkillRequest struct {
	Tag string `json:"tag" binding:"required"`
}

// RemoveUserSkill 移除用户技能标签
func (h *TagHandler) RemoveUserSkill(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的用户ID")
		return
	}

	// 验证只能修改自己的技能标签
	currentUserID, _ := c.Get("userID")
	if currentUserID.(int64) != userID {
		response.Error(c, http.StatusForbidden, 40302, "无权操作")
		return
	}

	var req RemoveUserSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}

	svcErr := h.skillSvc.RemoveUserSkill(userID, req.Tag)
	if svcErr != nil {
		if bizErr, ok := svcErr.(*bizerr.BizError); ok {
			response.Error(c, http.StatusInternalServerError, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, 50000, "服务器内部错误")
		}
		return
	}

	response.Success(c, nil)
}

// SearchBySkill 根据技能搜索用户
func (h *TagHandler) SearchBySkill(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		response.Error(c, http.StatusBadRequest, 40002, "缺少搜索关键词")
		return
	}

	userIDs, svcErr := h.skillSvc.SearchBySkill(tag)
	if svcErr != nil {
		if bizErr, ok := svcErr.(*bizerr.BizError); ok {
			response.Error(c, http.StatusInternalServerError, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, 50000, "服务器内部错误")
		}
		return
	}

	response.Success(c, gin.H{"user_ids": userIDs})
}
