package blog2

import (
	"fmt"
	"task3/datebeas"
)

// User 用户模型
type User struct {
	ID       uint
	Username string

	Posts []Post
}

// Post 文章模型
type Post struct {
	ID       uint
	Title    string
	Content  string
	UserID   uint `gorm:"not null;index"`
	Comments []Comment
}

// Comment 评论模型
type Comment struct {
	ID      uint
	Content string
	UserID  uint `gorm:"not null;index"`
	PostID  uint `gorm:"not null;index"`
}

var db = datebeas.DB

// 使用Gorm查询某个用户发布的所有文章及其对应的评论信息。

func Query1() {

	var user User
	db.Where("username = ?", "john_doe").Preload("Posts.Comments").Find(&user)
	fmt.Println(user)

}

type result struct {
	PostID uint
	Count  int
}

// 查询评论数量最多的文章信息
func Query2() {

	var result Post

	err := db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count DESC").
		First(&result).Error

	if err == nil {
		fmt.Println(result)
	}

}
