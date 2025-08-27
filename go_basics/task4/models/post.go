package models

type Post struct {
	ID            uint      `gorm:"primaryKey"`
	Title         string    `gorm:"not null"`
	Content       string    `gorm:"type:text"`
	UserID        uint      `gorm:"not null;index"`
	CommentCount  int       `gorm:"default:0"`         // 文章评论数量统计字段
	CommentStatus string    `gorm:"default:'无评论'"`     // 文章评论状态
	Comments      []Comment `gorm:"foreignKey:PostID"` // 一对多关系
}
