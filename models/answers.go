package models

type Answer struct {
	ID         uint   `gorm:"primaryKey"`
	AnswerDesc string `gorm:"not null"`
	Is_correct bool   `gorm:"not null"`
	QuestionID uint   `gorm:"not null"`
}

func (Answer) TableName() string {
	return "answer"
}
