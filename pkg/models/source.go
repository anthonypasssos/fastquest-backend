package models

import "encoding/json"

type SourceDoc struct {
	ID       uint                   `json:"id"`
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Source struct {
	ID       uint            `gorm:"primaryKey"`
	Name     string          `gorm:"not null"`
	Type     string          `gorm:"not null"`
	Metadata json.RawMessage `gorm:"type:json"`
}

func (Source) TableName() string {
	return "source"
}
