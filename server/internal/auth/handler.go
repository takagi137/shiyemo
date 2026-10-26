package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	service *AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		service: NewAuthService(),
	}
}

// RegisterParams 注册参数（传递给 service）
type RegisterParams struct {
	Email     string
	Password  string
	Nickname  string
	City      string
	Industry  string
	LeftJobAt string // 格式: 2024-01-15
}

// RegisterByEmailRequest 邮箱注册请求
type RegisterByEmailRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required,min=2,max=50"`
	// 以下为可选字段
	City      string `json:"city"`
	Industry  string `json:"industry"`
	LeftJobAt string `json:"left_job_at"` // 格式: 2024-01-15
}

// RegisterByPhoneRequest 手机号注册请求
type RegisterByPhoneRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Nickname string `json:"nickname" binding:"required,min=2,max=50"`
}

// LoginByEmailRequest 邮箱登录请求
type LoginByEmailRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginByPhoneRequest 手机号登录请求
type LoginByPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// SendSMSRequest 发送验证码请求
type SendSMSRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// RegisterByEmail 邮箱注册
func (h *AuthHandler) RegisterByEmail(c *gin.Context) {
	var req RegisterByEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	token, user, err := h.service.RegisterByEmail(RegisterParams{
		Email:     req.Email,
		Password:  req.Password,
		Nickname:  req.Nickname,
		City:      req.City,
		Industry:  req.Industry,
		LeftJobAt: req.LeftJobAt,
	})
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
		"token": token,
		"user":  user,
	})
}

// RegisterByPhone 手机号注册
func (h *AuthHandler) RegisterByPhone(c *gin.Context) {
	var req RegisterByPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	token, user, err := h.service.RegisterByPhone(req.Phone, req.Code, req.Nickname)
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
		"token": token,
		"user":  user,
	})
}

// LoginByEmail 邮箱登录
func (h *AuthHandler) LoginByEmail(c *gin.Context) {
	var req LoginByEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	token, user, err := h.service.LoginByEmail(req.Email, req.Password)
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
		"token": token,
		"user":  user,
	})
}

// LoginByPhone 手机号登录
func (h *AuthHandler) LoginByPhone(c *gin.Context) {
	var req LoginByPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	token, user, err := h.service.LoginByPhone(req.Phone, req.Code)
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
		"token": token,
		"user":  user,
	})
}

// SendSMS 发送验证码
func (h *AuthHandler) SendSMS(c *gin.Context) {
	var req SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "参数错误: "+err.Error())
		return
	}

	if err := h.service.SendSMS(req.Phone); err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, "发送验证码失败")
		return
	}

	response.Success(c, gin.H{
		"message": "验证码已发送",
	})
}

// GetProfile 获取当前用户信息（需要认证）
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	user, err := h.service.GetProfile(userID.(int64))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrDB.Code, bizerr.ErrDB.Message)
		return
	}

	response.Success(c, user)
}
