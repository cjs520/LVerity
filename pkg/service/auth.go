package service

import (
	"errors"
	"LVerity/pkg/config"
	"github.com/golang-jwt/jwt"
	"time"
)

// ErrInvalidToken 无效的令牌
var ErrInvalidToken = errors.New("invalid token")

// Claims JWT 声明结构体
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	RoleID   string `json:"role_id"`
	jwt.StandardClaims
}

// ValidateToken 验证JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GenerateToken 生成JWT token
func GenerateToken(userID, username, roleID string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 24小时后过期
			IssuedAt:  time.Now().Unix(),
			Issuer:    "LVerity",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
}

// CheckPermission 检查权限
func CheckPermission(roleID string, resource string, action string) bool {
	// TODO: 实现权限检查逻辑
	return true
}

// HasPermission 检查用户权限
func HasPermission(userID string, resource string, action string) bool {
	user, err := GetUserByID(userID)
	if err != nil {
		return false
	}
	return CheckPermission(user.RoleID, resource, action)
}
