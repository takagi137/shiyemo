package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/auth"
	"github.com/shiyemo/server/internal/common/database"
	"github.com/shiyemo/server/internal/common/response"
)

// AdminAuth 管理员权限中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户 ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusUnauthorized, 40101, "未登录")
			c.Abort()
			return
		}

		// 查询用户信息
		var user auth.User
		if err := database.GetDB().Where("id = ?", userID).First(&user).Error; err != nil {
			response.Error(c, http.StatusUnauthorized, 40101, "用户不存在")
			c.Abort()
			return
		}

		// 检查是否是管理员
		if !user.IsAdmin {
			response.Error(c, http.StatusForbidden, 40310, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
