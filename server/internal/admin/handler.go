package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/common/response"
)

// Handler 管理员处理器
type Handler struct {
	service *Service
}

// NewHandler 创建管理员处理器
func NewHandler() *Handler {
	return &Handler{
		service: NewService(),
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin")
	// 管理员中间件会在 cmd/main.go 中添加
	{
		admin.GET("/posts", h.ListPosts)
		admin.DELETE("/posts/:id", h.DeletePost)
		admin.GET("/users", h.ListUsers)
		admin.PUT("/users/:id/status", h.UpdateUserStatus)
		admin.GET("/stats", h.GetStats)
	}
}

// ListPosts 获取帖子列表
func (h *Handler) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.DefaultQuery("status", "active")

	posts, total, err := h.service.ListPosts(page, pageSize, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取帖子列表失败")
		return
	}

	response.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     posts,
	})
}

// DeletePost 删除帖子
func (h *Handler) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的帖子ID")
		return
	}

	if err := h.service.DeletePost(postID); err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "删除帖子失败")
		return
	}

	response.Success(c, nil)
}

// ListUsers 获取用户列表
func (h *Handler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.DefaultQuery("search", "")
	var status *int
	if s := c.Query("status"); s != "" {
		v, err := strconv.Atoi(s)
		if err == nil {
			status = &v
		}
	}

	users, total, err := h.service.ListUsers(page, pageSize, search, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取用户列表失败")
		return
	}

	// 清除敏感信息
	var safeUsers []map[string]interface{}
	for _, user := range users {
		safeUsers = append(safeUsers, map[string]interface{}{
			"id":         user.ID,
			"nickname":   user.Nickname,
			"email":      user.Email,
			"phone":      user.Phone,
			"is_admin":   user.IsAdmin,
			"status":     user.Status,
			"created_at": user.CreatedAt,
		})
	}

	response.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     safeUsers,
	})
}

// UpdateUserStatus 更新用户状态
func (h *Handler) UpdateUserStatus(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的用户ID")
		return
	}

	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}

	if req.Status != 0 && req.Status != 1 {
		response.Error(c, http.StatusBadRequest, 40001, "状态值无效")
		return
	}

	if err := h.service.UpdateUserStatus(userID, req.Status); err != nil {
		if err.Error() == "用户不存在" {
			response.Error(c, http.StatusNotFound, 40401, "用户不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, 50001, "更新用户状态失败")
		return
	}

	response.Success(c, nil)
}

// GetStats 获取平台统计数据
func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

// SetPostTop 设置帖子置顶
func (h *Handler) SetPostTop(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的帖子ID")
		return
	}

	var req struct {
		IsTop bool `json:"is_top"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}

	if err := h.service.SetPostTop(postID, req.IsTop); err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "设置置顶失败")
		return
	}

	response.Success(c, nil)
}

// SetPostFeatured 设置帖子加精
func (h *Handler) SetPostFeatured(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的帖子ID")
		return
	}

	var req struct {
		IsFeatured bool `json:"is_featured"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误")
		return
	}

	if err := h.service.SetPostFeatured(postID, req.IsFeatured); err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "设置加精失败")
		return
	}

	response.Success(c, nil)
}

// ListPendingPosts 获取待审核帖子列表
func (h *Handler) ListPendingPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	posts, total, err := h.service.ListPendingPosts(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取待审核帖子失败")
		return
	}

	response.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     posts,
	})
}

// ApprovePost 审核通过帖子
func (h *Handler) ApprovePost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的帖子ID")
		return
	}

	reviewerID := c.GetInt64("user_id")
	if reviewerID == 0 {
		response.Error(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}

	if err := h.service.ApprovePost(postID, reviewerID); err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "审核通过失败")
		return
	}

	response.Success(c, nil)
}

// RejectPost 审核拒绝帖子
func (h *Handler) RejectPost(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的帖子ID")
		return
	}

	reviewerID := c.GetInt64("user_id")
	if reviewerID == 0 {
		response.Error(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}

	if err := h.service.RejectPost(postID, reviewerID, req.Reason); err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "审核拒绝失败")
		return
	}

	response.Success(c, nil)
}

// ListComments 获取评论列表（管理员）
func (h *Handler) ListComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	postID, _ := strconv.ParseInt(c.Query("post_id"), 10, 64)

	comments, total, err := h.service.ListComments(page, pageSize, postID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "获取评论列表失败")
		return
	}

	response.Success(c, gin.H{
		"items":     comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// HideComment 屏蔽评论
func (h *Handler) HideComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的评论ID")
		return
	}

	if err := h.service.HideComment(commentID); err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "屏蔽评论失败")
		return
	}

	response.Success(c, nil)
}

// RestoreComment 恢复评论
func (h *Handler) RestoreComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "无效的评论ID")
		return
	}

	if err := h.service.RestoreComment(commentID); err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "恢复评论失败")
		return
	}

	response.Success(c, nil)
}
