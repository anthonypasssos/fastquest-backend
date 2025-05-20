package models

import "time"

type Comment struct {
	ID           uint   `gorm:"primaryKey"`
	Text	     string `gorm:"not null"`
	CreationDate time.Time
	UserID       int `gorm:"not null"`
}

func (Comment) TableName() string {
	return "comment"
}
