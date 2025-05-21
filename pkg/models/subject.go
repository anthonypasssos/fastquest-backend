package models

type Subject struct {
	ID         uint   `gorm:"primaryKey"`
	Name	   string `gorm:"not null"`
}

func (Subject) TableName() string {
	return "subject"
}
