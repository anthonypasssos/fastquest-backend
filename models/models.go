package models

import "time"

type Question struct {
    ID        uint   `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    Statement string `gorm:"not null"`
    Subject   int    
    UserID    int    `gorm:"not null"`
}
func (Question) TableName() string {
	return "question"
}