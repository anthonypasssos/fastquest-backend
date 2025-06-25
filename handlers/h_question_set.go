package handlers

import (
	"encoding/json"
	"flashquest/database"
	"flashquest/pkg/models"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type NewList struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Desc      string `json:"desc"`
	IsPrivate bool   `json:"is_private"`
	UserID    int    `json:"user_id"`
	Questions []int  `json:"questions"`
}

func CreateQuestionSet(w http.ResponseWriter, r *http.Request) {
	var newList NewList

	// Decodifica o JSON enviado
	err := json.NewDecoder(r.Body).Decode(&newList)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	// Cria o question set
	questionSet := models.QuestionSet{
		Name:      newList.Name,
		Type:      newList.Type,
		Desc:      newList.Desc,
		UserID:    newList.UserID,
		IsPrivate: newList.IsPrivate,
	}

	// Inicia uma transação
	err = db.Transaction(func(tx *gorm.DB) error {
		// Cria o question_set
		if err := tx.Create(&questionSet).Error; err != nil {
			return err
		}

		// Cria as relações question_set_question
		for index, questionID := range newList.Questions {
			link := models.QuestionSetQuestion{
				QuestionSetID: questionSet.ID,
				QuestionID:    questionID,
				Position:      index + 1,
			}

			if err := tx.Create(&link).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating question set: %v", err), http.StatusInternalServerError)
		return
	}

	// Retorna o question set criado
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questionSet)
}
