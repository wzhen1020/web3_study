package middlewares

import (
	"net/http"
	"task4/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 token
		tokenString := c.GetHeader("Authorization")
		// 支持从 Cookie 读取（可选）
		if tokenString == "" {
			if cookie, err := c.Cookie("jwt_token"); err == nil {
				tokenString = cookie
			}
		}

		// 检查 token 是否存在
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求未携带token，无权限访问"})
			c.Abort()
			return
		}

		// 解析和验证 token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token或token已过期"})
			c.Abort()
			return
		}

		// 将解析出的 claims 信息存入上下文，方便后续处理器使用
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)

		c.Next() // 处理下一个中间件或请求
	}
}
