package models

import "time"

type Origin struct {
	Exam ExamInstance
}

type ExamInstance struct {
	ID        uint   `gorm:primaryKey`
	SourceId  uint   `gorm:"not null"`
	Source    Source `gorm:"foreignKey:SourceId"`
	Edition   uint   `gorm:"not null"`
	Phase     uint   `gorm:"not null"`
	Year      uint   `gorm:"not null"`
	CreatedAt time.Time
}

func (ExamInstance) TableName() string {
	return "source_exam_instance"
}
