package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shiyemo/server/internal/common/config"
)

// Claims JWT 声明
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	ErrTokenExpired = errors.New("token 已过期")
	ErrTokenInvalid = errors.New("token 无效")
)

// GenerateToken 生成 JWT token
func GenerateToken(userID int64) (string, error) {
	cfg := config.Load()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWT.ExpireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "shiyemo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.Load()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}
