package handlers

import (
	"encoding/json"
	"flashquest/database"
	"flashquest/pkg/models"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func CreateSource(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	var body models.SourceExamBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.Name == "" || body.Year == 0 {
		http.Error(w, "Name and Year are required", http.StatusBadRequest)
		return
	}

	source := models.Source{
		Name: body.Name,
		Type: body.Type,
	}

	var generatedID uint

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&source).Error; err != nil {
			return err
		}

		examInstance := models.ExamInstance{
			SourceId: source.ID,
			Edition:  body.Edition,
			Phase:    body.Phase,
			Year:     body.Year,
		}

		if err := tx.Create(&examInstance).Error; err != nil {
			return err
		}

		generatedID = source.ID
		return nil
	})

	if err != nil {
		fmt.Printf("Error creating source/instance: %v\n", err)

		http.Error(w, "Failed to create source instance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Source Instance created successfully",
		"id":      generatedID,
	})
}
