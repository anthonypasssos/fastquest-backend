package models

import (
	"time"
)

type QuestionDoc struct {
	ID        uint        `json:"id"`
	Statement string      `json:"statement"`
	SubjectID int         `json:"subject_id"`
	UserID    int         `json:"user_id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Sources   []SourceDoc `json:"sources"`
}

type Question struct {
	Sources   []Source `gorm:"many2many:question_sources;"`
	ID        uint     `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Statement string `gorm:"not null"`
	SubjectID int    `gorm:"not null"`
	UserID    int    `gorm:"not null"`
}

type QuestionListResponse struct {
	Data       []QuestionDoc `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

func (Question) TableName() string {
	return "question"
}
