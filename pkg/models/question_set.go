package models

import (
	"time"
)

type QuestionSet struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Type         string    `json:"type"`
	UserID       int       `json:"user_id"`
	CreationDate time.Time `json:"creation_date" gorm:"autoCreateTime"`
	IsPrivate    bool      `json:"is_private"`
}

type QuestionSetQuestion struct {
	ID            int `json:"id" gorm:"primaryKey;autoIncrement"`
	QuestionSetID int `json:"question_set_id"`
	QuestionID    int `json:"question_id"`
	Position      int `json:"position"`
}

func (QuestionSet) TableName() string {
	return "question_set"
}

func (QuestionSetQuestion) TableName() string {
	return "question_set_question"
}
