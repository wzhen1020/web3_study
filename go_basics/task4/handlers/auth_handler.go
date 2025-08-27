package handlers

import (
	"net/http"
	"task4/datebeas"
	"task4/models"
	"task4/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var db = datebeas.DB

// Login 登录处理器
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 1. 验证用户名和密码（这里只是示例，实际应查询数据库）
	// 假设验证通过，获取用户ID
	userID := int64(1)
	expectedUsername := req.Username

	user := models.User{}
	db.Where("username = ?", expectedUsername).Find(&user)

	expectedPassword := user.Password
	if req.Username != expectedUsername || req.Password != expectedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 2. 生成 JWT Token
	token, err := utils.GenerateToken(userID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成Token失败"})
		return
	}

	// 3. 返回 token 给客户端（也可以设置为 Cookie）
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": gin.H{
			"id":       userID,
			"username": req.Username,
		},
	})
}

// GetUserProfile 获取用户信息（受保护路由示例）
func GetUserProfile(c *gin.Context) {
	// 从上下文中获取中间件设置的 claims 信息
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	c.JSON(http.StatusOK, gin.H{
		"message": "访问成功",
		"user": gin.H{
			"id":       userID,
			"username": username,
		},
	})
}
