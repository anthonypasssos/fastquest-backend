package main

import (
	"encoding/json"
	"fmt"
	"flashquest/database"
	"flashquest/models"
	"net/http"
)

func getQuestions(w http.ResponseWriter, r *http.Request) {
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
	for i, q := range questions {
		fmt.Printf("%d: %s\n", i+1, q.Statement)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(questions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}