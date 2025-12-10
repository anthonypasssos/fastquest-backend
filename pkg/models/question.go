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

type QuestionResponse struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Statement string         `json:"statement"`
	Subject   *Subject       `json:"subject,omitempty"`
	User      *User          `json:"user,omitempty"`
	Source    *UnifiedSource `json:"source,omitempty"`
	Answers   *[]Answer      `json:"answers,omitempty"`
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

func (q Question) ToResponse() QuestionResponse {
	resp := QuestionResponse{
		ID:        q.ID,
		CreatedAt: q.CreatedAt,
		UpdatedAt: q.UpdatedAt,
		Statement: q.Statement,
	}

	if q.Subject != nil {
		resp.Subject = q.Subject
	}

	if q.User != nil {
		resp.User = q.User
	}

	if q.Answers != nil && len(*q.Answers) > 0 {
		resp.Answers = q.Answers
	}

	if q.SourceExamInstance != nil {
		resp.Source = &UnifiedSource{
			ID:   q.SourceExamInstance.Source.ID,
			Name: q.SourceExamInstance.Source.Name,
			Type: q.SourceExamInstance.Source.Type,
			Metadata: map[string]interface{}{
				"year":    q.SourceExamInstance.Year,
				"edition": q.SourceExamInstance.Edition,
				"phase":   q.SourceExamInstance.Phase,
			},
		}
	}

	return resp
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
			case "user":
				db = db.Preload("User")
			case "subject":
				db = db.Preload("Subject")
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
