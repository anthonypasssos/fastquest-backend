package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	Statement string `gorm:"column:statemente;not null"`
	Subject   int    `gorm:"column:subject"`
	UserID    int    `gorm:"column:user_id;not null"` // Using standard naming convention
}

func (Question) TableName() string {
	return "questions"
}