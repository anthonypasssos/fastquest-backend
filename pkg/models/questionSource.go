package models

type QuestionSource struct {
	QuestionID uint `gorm:"primaryKey"`
	SourceID   uint `gorm:"primaryKey"`
}

func (QuestionSource) TableName() string {
	return "question_source"
}
