package verify

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// SendResetEmailRequest 发送重置邮件请求
// @Description 发送密码重置邮件
// @Accept json
// @Produce json
// @Param email body SendResetEmailRequest true "邮箱"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/forgot-password [post]

// VerifyHandler 验证处理器
type VerifyHandler struct {
	service *VerifyService
}

// NewVerifyHandler 创建验证处理器
func NewVerifyHandler() *VerifyHandler {
	return &VerifyHandler{
		service: NewVerifyService(),
	}
}

// SendVerificationEmail 发送验证邮件
func (h *VerifyHandler) SendVerificationEmail(c *gin.Context) {
	// 从 JWT 中获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	// 发送验证邮件
	if err := h.service.SendVerificationEmail(userID.(int64)); err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, gin.H{
		"message": "验证邮件已发送",
	})
}

// SendResetEmail 发送密码重置邮件
func (h *VerifyHandler) SendResetEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "请输入有效的邮箱地址")
		return
	}

	if err := h.service.SendResetEmail(req.Email); err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, gin.H{
		"message": "重置邮件已发送",
	})
}

// ResetPassword 重置密码
func (h *VerifyHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误：token 和新密码（至少6位）必填")
		return
	}

	if err := h.service.ResetPassword(req.Token, req.NewPassword); err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, gin.H{
		"message": "密码重置成功",
	})
}

// VerifyEmail 验证邮箱
func (h *VerifyHandler) VerifyEmail(c *gin.Context) {
	// 从 URL 参数获取 token
	token := c.Param("token")
	if token == "" {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "验证 token 不能为空")
		return
	}

	// 验证邮箱
	if err := h.service.VerifyEmail(token); err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, gin.H{
		"message": "邮箱验证成功",
	})
}
