package models

type Answer struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Text       string `gorm:"not null" json:"text"`
	Is_correct bool   `gorm:"not null" json:"is_correct"`
	QuestionID uint   `gorm:"column:id_question; not null" json:"question_id"`
}

func (Answer) TableName() string {
	return "answer"
}
