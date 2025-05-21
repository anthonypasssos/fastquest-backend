package models

type Topic struct {
	ID         uint   `gorm:"primaryKey"`
	Name	   string `gorm:"not null"`
	SubjectID  int `gorm:"not null"`
}

func (Topic) TableName() string {
	return "topic"
}
