package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shiyemo/server/internal/common/errors"
	jwtutil "github.com/shiyemo/server/internal/common/jwt"
	"github.com/shiyemo/server/internal/common/response"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Error(c, http.StatusUnauthorized, errors.ErrUnauthorized.Code, errors.ErrUnauthorized.Message)
			c.Abort()
			return
		}

		// 去掉 Bearer 前缀
		token = strings.TrimPrefix(token, "Bearer ")

		// 解析 JWT
		claims, err := jwtutil.ParseToken(token)
		if err != nil {
			if err == jwtutil.ErrTokenExpired {
				response.Error(c, http.StatusUnauthorized, errors.ErrTokenExpired.Code, errors.ErrTokenExpired.Message)
			} else {
				response.Error(c, http.StatusUnauthorized, errors.ErrTokenInvalid.Code, errors.ErrTokenInvalid.Message)
			}
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
