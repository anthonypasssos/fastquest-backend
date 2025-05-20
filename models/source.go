package models

import "encoding/json"

type Source struct {
	ID         uint   `gorm:"primaryKey"`
	Name	   string `gorm:"not null"`
	Type  	   string `gorm:"not null"`
	Metadata   json.RawMessage `gorm:"type:json"`
}

func (Source) TableName() string {
	return "source"
}
