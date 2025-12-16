package models

import (
	"time"

	"gorm.io/gorm"
)

type QuestionSetResponse struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Type        string             `json:"type"`
	User        *UserResponse      `json:"user,omitempty"`
	Questions   []QuestionResponse `json:"questions,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	IsPrivate   bool               `json:"is_private"`
}

type QuestionSet struct {
	ID          int `gorm:"primaryKey;autoIncrement"`
	Name        string
	Description string
	Type        string `gorm:"not null"`
	UserID      int
	User        *User      `gorm:"foreignKey:UserID"`
	Questions   []Question `gorm:"many2many:question_set_question;"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	IsPrivate   bool       `gorm:"not null"`
}

type QuestionSetQuestion struct {
	ID            int `json:"id" gorm:"primaryKey;autoIncrement"`
	QuestionSetID int `json:"question_set_id"`
	QuestionID    int `json:"question_id"`
	Position      int `json:"position"`
}

func (qs QuestionSet) ToResponse() QuestionSetResponse {
	set := QuestionSetResponse{
		ID:          qs.ID,
		Name:        qs.Name,
		Description: qs.Description,
		Type:        qs.Type,
		CreatedAt:   qs.CreatedAt,
		IsPrivate:   qs.IsPrivate,
	}

	if qs.User != nil {
		resUser := qs.User.ToResponse()
		set.User = &resUser
	}

	if len(qs.Questions) > 0 {
		var resQuestions []QuestionResponse
		for _, q := range qs.Questions {
			resQuestions = append(resQuestions, q.ToResponse())
		}
		set.Questions = resQuestions
	}

	return set

}

func ApplyQuestionSetIncludes(includes []string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, include := range includes {
			switch include {
			case "questions":
				db = db.Preload("Questions")
			case "user":
				db = db.Preload("User")
			}
		}
		return db
	}
}

func (QuestionSet) TableName() string {
	return "question_set"
}

func (QuestionSetQuestion) TableName() string {
	return "question_set_question"
}
