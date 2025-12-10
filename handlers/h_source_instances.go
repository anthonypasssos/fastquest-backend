package handlers

import (
	"errors"
	"flashquest/database"
	"flashquest/pkg/models"
	"fmt"

	"gorm.io/gorm"
)

func SendExamInstance(ei ...*models.ExamInstanceBody) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := db.Create(ei).Error; err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

func GetInstanceWithSource(eiID int, ei *models.ExamInstance) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	result := db.Preload("SourceExamInstance.Source").Where("id = ?", eiID).Find(&ei)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("Source Exam Instance not found")
		} else {
			return errors.New("Error fetching Source Exam Instance")
		}
	}

	return nil
}
