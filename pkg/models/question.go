package models

import (
	"time"

	"gorm.io/gorm"
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
	ID                   uint `gorm:"primaryKey"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Statement            string   `gorm:"not null"`
	SubjectID            int      `gorm:"not null"`
	Subject              *Subject `gorm:"foreignKey:SubjectID"`
	UserID               int      `gorm:"not null"`
	User                 *User    `gorm:"foreignKey:UserID"`
	SourceExamInstanceID *int
	SourceExamInstance   *ExamInstance `gorm:"foreignKey:SourceExamInstanceID"`
	Answers              *[]Answer
}

func ApplyIncludes(includes []string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, include := range includes {
			switch include {
			case "source":
				// Carrega ambos os caminhos possíveis de source + os dados da fonte
				db = db.Preload("SourceExamInstance.Source")
			case "answers":
				db = db.Preload("Answers")
			case "comments":
				db = db.Preload("Comments")
			case "user":
				db = db.Preload("User")
			case "subject":
				db = db.Preload("Subject")
			case "answers.author": // Exemplo avançado: Autor da resposta
				db = db.Preload("Answers.User")
			}
		}
		return db
	}
}

type QuestionListResponse struct {
	Data       []QuestionDoc `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

func (Question) TableName() string {
	return "question"
}
