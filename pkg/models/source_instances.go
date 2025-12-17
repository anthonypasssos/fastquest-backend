package models

import "time"

type Origin struct {
	Exam ExamInstance
}

type ExamInstance struct {
	ID        uint   `gorm:primaryKey json:"id"`
	SourceId  uint   `gorm:"not null" json:"source_id"`
	Source    Source `gorm:"foreignKey:SourceId" json:"source"`
	Edition   uint   `gorm:"not null" json:"edition"`
	Phase     uint   `gorm:"not null" json:"phase"`
	Year      uint   `gorm:"not null" json:"year"`
	CreatedAt time.Time
}

func (ExamInstance) TableName() string {
	return "source_exam_instance"
}
