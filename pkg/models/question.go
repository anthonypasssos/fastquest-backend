package models

import (
	"time"

	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	Sources   []Source `gorm:"many2many:question_sources;"`
	ID        uint     `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Statement string `gorm:"not null"`
	SubjectID int    `gorm:"not null"`
	UserID    int    `gorm:"not null"`
}

func (Question) TableName() string {
	return "question"
}
