package models

type SubjectResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Subject struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

func (s Subject) ToResponse() SubjectResponse {
	return SubjectResponse{
		ID:   s.ID,
		Name: s.Name,
	}
}

func (Subject) TableName() string {
	return "subject"
}
