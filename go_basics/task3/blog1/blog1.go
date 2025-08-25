package blog1

import (
	"fmt"
	"task3/datebeas"
)

// User 用户模型
type User struct {
	ID       uint
	Username string
}

// Post 文章模型
type Post struct {
	ID      uint
	Title   string
	Content string
	UserID  uint `gorm:"not null;index"`
}

// Comment 评论模型
type Comment struct {
	ID      uint
	Content string
	UserID  uint `gorm:"not null;index"`
	PostID  uint `gorm:"not null;index"`
}

var db = datebeas.DB

func CreateTable() {
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
}

// 创建示例数据
func CreateSampleData() {
	// 创建用户
	user := User{
		Username: "john_doe",
	}
	db.Create(&user)

	// 创建文章
	post := Post{
		Title:   "我的第一篇博客文章",
		Content: "这是文章内容...",
		UserID:  user.ID,
	}
	db.Create(&post)

	// 创建评论
	comment := Comment{
		Content: "这是一条评论",
		UserID:  user.ID,
		PostID:  post.ID,
	}
	db.Create(&comment)

	fmt.Println("示例数据创建成功!")
}
