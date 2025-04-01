package main

import (
	"net/http"
	"encoding/json"
	"flashquest/database"
	"flashquest/models"
	"fmt"
)

func getQuestions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting questions")
	var questions []models.Question
	db := database.GetDB()

	if err := db.Find(&questions).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}