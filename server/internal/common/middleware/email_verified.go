package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/auth"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
	"github.com/shiyemo/server/internal/setting"
)

// EmailVerified 邮箱验证中间件
// 未验证邮箱的用户发帖时返回提示。
//
// 是否强制由管理后台的「强制邮箱验证」开关决定（system_settings 表）。
// 关闭时直接放行 —— 否则在 SMTP 未配置的环境下，
// 新用户收不到验证邮件，就永远无法发帖（死锁）。
func EmailVerified() gin.HandlerFunc {
	settingService := setting.NewSettingService()

	return func(c *gin.Context) {
		if !settingService.IsEmailVerifyRequired() {
			c.Next()
			return
		}

		// 从 JWT 中获取用户 ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
			c.Abort()
			return
		}

		// 获取用户信息
		userRepo := auth.NewAuthRepo()
		user, err := userRepo.GetByID(userID.(int64))
		if err != nil {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrDB.Code, bizerr.ErrDB.Message)
			c.Abort()
			return
		}

		// 检查邮箱是否已验证
		if !user.IsEmailVerified() {
			response.Error(c, http.StatusForbidden, 40301, "请先验证邮箱后再进行此操作")
			c.Abort()
			return
		}

		c.Next()
	}
}
