package handlers

import (
	"encoding/json"
	"errors"
	"flashquest/database"
	"flashquest/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetQuestions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting questions")

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	var questions []models.Question
	result := db.Find(&questions)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching questions: %v", result.Error),
			http.StatusInternalServerError)
		return
	}

	fmt.Printf("Found %d questions\n", len(questions))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(questions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func GetQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	var question models.Question
	result := db.Where("id = ?", id).First(&question)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Question not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error fetching question: %v", result.Error),
				http.StatusInternalServerError)
		}
		return
	}

	fmt.Printf("Found question %s \n", id)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(question); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	log.Printf("Attempting to delete question ID: %s", id)

	result := db.Delete(&models.Question{}, id)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error deleting question: %v", result.Error),
			http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	log.Printf("Question ID %s deleted successfully", id)
}
