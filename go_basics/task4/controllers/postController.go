package controllers

import (
	"net/http"
	"task4/datebeas"
	"task4/models"
	"task4/utils"

	"github.com/gin-gonic/gin"
)

type PostController struct {
}

var db = datebeas.DB
var logger = utils.GetLogger()

// 创建文章
func (con PostController) Create(c *gin.Context) {

	var post models.Post
	logger.Info("创建文章", map[string]interface{}{
		"post": post,
	})
	userID, _ := c.Get("userId")
	// fmt.Printf("%v---", userID)

	if err := c.ShouldBindJSON(&post); err != nil {
		logger.Error("创建文章失败", map[string]interface{}{
			"post": post,
		})
		c.JSON(http.StatusBadRequest, gin.H{

			"error":   "参数错误",
			"details": err.Error(),
		})
		return
	}
	post.UserID = userID.(uint)
	// 5. 数据库操作（需要访问控制器持有的数据库实例）
	if err := db.Create(&post).Error; err != nil {
		logger.Error("创建文章数据库存储失败", map[string]interface{}{
			"post": post,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "创建文章失败",
			"details": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": "创建文章成功",
		})
		return
	}

}

//修改文章

func (con PostController) Edit(c *gin.Context) {
	var post models.Post
	userID, _ := c.Get("userId")
	// fmt.Printf("%v---", userID)

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"details": err.Error(),
		})
		return
	}
	logger.Info("修改文章", map[string]interface{}{
		"post": post,
	})
	if post.UserID == userID.(uint) {

		if err := db.Save(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "修改文章失败",
				"details": err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"error": "修改文章成功",
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": "无权限",
		})
		return
	}

}

//删除文章

func (con PostController) Delete(c *gin.Context) {

	postId := c.Query("postId")
	userID, _ := c.Get("userId")

	logger.Info("删除文章", map[string]interface{}{
		"postId": postId,
	})

	if err := db.Where("id = ? and user_id = ?", postId, userID.(uint)).Delete(&models.Post{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "删除文章失败",
			"details": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": "删除文章成功",
		})
		return
	}
}

// 读取文章列表
func (con PostController) QueryList(c *gin.Context) {

	postList := []models.Post{}

	db.Order("id desc").Find(&postList)
	c.JSON(200, gin.H{
		"result": postList,
	})
}

// 读取文章详情
func (con PostController) QueryPostInfo(c *gin.Context) {

	postId := c.Query("postId")

	var post models.Post
	db.Where("id = ?", postId).Find(&post)

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}
