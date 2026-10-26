package comment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	service *CommentService
}

// NewCommentHandler 创建评论处理器
func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		service: NewCommentService(),
	}
}

// Create 创建评论
func (h *CommentHandler) Create(c *gin.Context) {
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

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	comment, err := h.service.Create(userID.(int64), postID, req)
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, comment)
}

// List 获取评论列表
func (h *CommentHandler) List(c *gin.Context) {
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	comments, total, err := h.service.List(postID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		return
	}

	response.SuccessPage(c, comments, page, size, total)
}

// Delete 删除评论
func (h *CommentHandler) Delete(c *gin.Context) {
	commentIDStr := c.Param("id")
	if commentIDStr == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamMissing.Code, "评论ID不能为空")
		return
	}

	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "无效的评论ID")
		return
	}

	// 从 context 获取 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	if err := h.service.Delete(userID.(int64), commentID); err != nil {
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
