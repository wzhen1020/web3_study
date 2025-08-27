package models

type User struct {
	ID        uint `gorm:"primaryKey"`
	Username  string
	Email     string
	Password  string
	PostCount int `gorm:"default:0"` // 用户文章数量统计字段
}
