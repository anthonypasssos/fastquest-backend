package models

import "time"

type Question struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Statement string `gorm:"not null"`
	SubjectID int `gorm:"not null"`
	UserID    int `gorm:"not null"`

    User       User      `json:"user" gorm:"foreignKey:UserID"`
    TopicID    uint      `json:"-"`
    Topic      Topic     `json:"topic" gorm:"foreignKey:TopicID"`
    SourceID   uint      `json:"-"`
    Source     Source    `json:"source" gorm:"foreignKey:SourceID"`
    Answers    []Answer  `json:"answers"`
    Comments   []Comment `json:"comments"`
}

func (Question) TableName() string {
	return "question"
}

/*type Question struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Statement string `gorm:"not null"`
	SubjectID int `gorm:"not null"`
	UserID    int `gorm:"not null"`
}

func (Question) TableName() string {
	return "question"
}*/
