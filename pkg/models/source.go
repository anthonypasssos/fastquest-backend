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
	ID   uint   `json:"id"`   // ID da Fonte (Ex: ID da OAB)
	Name string `json:"name"` // Ex: "OAB"
	Type string `json:"type"` // "EXAM" ou "CLASS"

	// Aqui está o pulo do gato: um campo formatado ou genérico
	Info     string                 `json:"info"`     // Ex: "XXXVIII Exame - 2023"
	Metadata map[string]interface{} `json:"metadata"` // Dados extras flexíveis
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
