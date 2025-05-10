package models

type Answer struct {
	ID         uint   `gorm:"primaryKey"`
	Answer string `gorm:"not null"`
	Is_correct bool   `gorm:"not null"`
	QuestionID uint   `gorm:"column:id_question; not null"`
}

func (Answer) TableName() string {
	return "answer"
}
