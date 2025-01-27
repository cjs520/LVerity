package middleware

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"LVerity/pkg/config"
	"LVerity/pkg/service"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}

		claims := &service.Claims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GlobalConfig.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证信息"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleID", claims.RoleID)
		c.Next()
	}
}

// RequirePermission 权限检查中间件
func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("roleID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证的用户"})
			c.Abort()
			return
		}

		if !service.CheckPermission(roleID.(string), resource, action) {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有操作权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}
