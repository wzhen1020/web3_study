package main

import (
	"fmt"
	"log"
	"task3/datebeas"

	"gorm.io/gorm"
)

// 定义模型结构体
type User struct {
	ID        uint `gorm:"primaryKey"`
	Username  string
	Email     string
	PostCount int `gorm:"default:0"` // 用户文章数量统计字段
}

type Post struct {
	ID            uint      `gorm:"primaryKey"`
	Title         string    `gorm:"not null"`
	Content       string    `gorm:"type:text"`
	UserID        uint      `gorm:"not null;index"`
	CommentCount  int       `gorm:"default:0"`         // 文章评论数量统计字段
	CommentStatus string    `gorm:"default:'无评论'"`     // 文章评论状态
	Comments      []Comment `gorm:"foreignKey:PostID"` // 一对多关系
}

type Comment struct {
	ID      uint   `gorm:"primaryKey"`
	Content string `gorm:"type:text;not null"`
	PostID  uint   `gorm:"not null;index"`
	UserID  uint   `gorm:"not null;index"`
}

func main() {

	db := datebeas.DB

	// 自动迁移（创建表）
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatal("自动迁移失败:", err)
	}

	// 创建测试用户
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
	}
	db.Create(&user)

	// 创建测试文章
	post := Post{
		Title:   "测试文章",
		Content: "这是一篇测试文章的内容",
		UserID:  user.ID,
	}
	db.Create(&post)

	// 创建测试评论
	comment := Comment{
		Content: "这是一条测试评论",
		PostID:  post.ID,
		UserID:  user.ID,
	}
	db.Create(&comment)

	// 查询用户，检查文章数量是否更新
	var updatedUser User
	db.First(&updatedUser, user.ID)
	fmt.Printf("用户 %s 的文章数量: %d\n", updatedUser.Username, updatedUser.PostCount)

	// 查询文章，检查评论数量和状态
	var updatedPost Post
	db.Preload("Comments").First(&updatedPost, post.ID)
	fmt.Printf("文章 %s 的评论数量: %d, 评论状态: %s\n",
		updatedPost.Title, updatedPost.CommentCount, updatedPost.CommentStatus)

	// 删除评论，检查评论状态是否更新
	db.Delete(&comment)

	var finalPost Post
	db.First(&finalPost, post.ID)
	fmt.Printf("删除评论后，文章 %s 的评论数量: %d, 评论状态: %s\n",
		finalPost.Title, finalPost.CommentCount, finalPost.CommentStatus)
}

// Post 模型的钩子函数

// BeforeCreate 在创建文章前更新用户的文章数量
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	// 更新用户的文章数量
	result := tx.Model(&User{}).Where("id = ?", p.UserID).Update("post_count", gorm.Expr("post_count + 1"))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// BeforeDelete 在删除文章前更新用户的文章数量
func (p *Post) BeforeDelete(tx *gorm.DB) (err error) {
	// 更新用户的文章数量
	result := tx.Model(&User{}).Where("id = ?", p.UserID).Update("post_count", gorm.Expr("post_count - 1"))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Comment 模型的钩子函数

// AfterCreate 在创建评论后更新文章的评论数量和状态
func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	// 更新文章的评论数量
	result := tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_count", gorm.Expr("comment_count + 1"))
	if result.Error != nil {
		return result.Error
	}

	// 更新文章的评论状态
	result = tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "有评论")
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// AfterDelete 在删除评论后更新文章的评论数量和状态
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	// 更新文章的评论数量
	result := tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_count", gorm.Expr("comment_count - 1"))
	if result.Error != nil {
		return result.Error
	}

	// 检查文章的评论数量，如果为0则更新评论状态
	var post Post
	tx.First(&post, c.PostID)
	if post.CommentCount == 0 {
		result = tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "无评论")
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
