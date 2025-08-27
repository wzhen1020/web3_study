package models

type Comment struct {
	ID      uint   `gorm:"primaryKey"`
	Content string `gorm:"type:text;not null"`
	PostID  uint   `gorm:"not null;index"`
	UserID  uint   `gorm:"not null;index"`
}
