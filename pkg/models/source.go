package models

import (
	"time"
)

type SourceDoc struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type Source struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Type      string `gorm:"not null"`
	CreatedAt time.Time
}

type UnifiedSource struct {
	ID       uint                   `json:"id"`
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
}

type SourceExamBody struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Edition uint   `json:"edition"`
	Phase   uint   `json:"phase"`
	Year    uint   `json:"year"`
}

type ExamInstanceBody struct {
	CreatedAt time.Time `json:"created_at"`
}

func (Source) TableName() string {
	return "source"
}
