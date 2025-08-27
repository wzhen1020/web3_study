package controllers

import (
	"net/http"
	"task4/models"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
}

// 添加评论
func (con CommentController) Add(c *gin.Context) {
	var comment models.Comment
	userID, _ := c.Get("userId")
	// fmt.Printf("%v---", userID)

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"details": err.Error(),
		})
		return
	}
	comment.UserID = userID.(uint)
	// 5. 数据库操作（需要访问控制器持有的数据库实例）
	if err := db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "创建评论失败",
			"details": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": "创建评论成功",
		})
		return
	}
}

func (con CommentController) QueryList(c *gin.Context) {
	postId := c.Query("postId")
	commentList := []models.Comment{}

	db.Where("post_id = ?", postId).Find(&commentList)
	c.JSON(200, gin.H{
		"result": commentList,
	})
}
