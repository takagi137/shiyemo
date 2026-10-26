package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// Handler 系统设置处理器
type Handler struct {
	service SettingServiceInterface
}

// NewHandler 创建系统设置处理器
func NewHandler() *Handler {
	return &Handler{service: NewSettingService()}
}

// NewHandlerWithService 创建带依赖注入的系统设置处理器（用于测试）
func NewHandlerWithService(service SettingServiceInterface) *Handler {
	return &Handler{service: service}
}

// GetSMTP 获取 SMTP 配置（密码不回显）
func (h *Handler) GetSMTP(c *gin.Context) {
	cfg, err := h.service.GetSMTP()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrDB.Code, bizerr.ErrDB.Message)
		return
	}
	response.Success(c, cfg)
}

// UpdateSMTP 更新 SMTP 配置
func (h *Handler) UpdateSMTP(c *gin.Context) {
	var req UpdateSMTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}
	if err := h.service.UpdateSMTP(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, bizerr.ErrDB.Code, bizerr.ErrDB.Message)
		return
	}
	response.Success(c, gin.H{"updated": true})
}

// TestSMTP 发送测试邮件
func (h *Handler) TestSMTP(c *gin.Context) {
	var req TestSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}
	if err := h.service.TestSend(req.To); err != nil {
		// 把 SMTP 的原始错误透给管理员，否则无从排查配置问题
		response.Error(c, http.StatusBadGateway, 50002, "发送失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"sent": true})
}
