package handlers

import (
	"encoding/json"
	"errors"
	"flashquest/database"
	"flashquest/models"
	"fmt"
	"log"
	"net/http"

	"time"

	"math"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating question")

	var questionInput struct {
		Statement string `json:"statement"`
		SubjectID int    `json:"subject_id"`
		UserID    int    `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&questionInput); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if questionInput.Statement == "" {
		http.Error(w, "Statement is required", http.StatusBadRequest)
		return
	}

	if questionInput.UserID == 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	question := models.Question{
		Statement: questionInput.Statement,
		SubjectID: questionInput.SubjectID,
		UserID:    questionInput.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	result := db.Create(&question)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error creating question: %v", result.Error),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(question)
}

func GetQuestions(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	filter := query.Get("filter")

	orderBy := query.Get("order_by")
	if orderBy == "" {
		orderBy = "created_at desc"
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	queryBuilder := db.Model(&models.Question{})

	if filter != "" {
		queryBuilder = queryBuilder.Where("statement LIKE ?", "%"+filter+"%").
			Or("subject::text LIKE ?", "%"+filter+"%")
	}

	queryBuilder = queryBuilder.Order(orderBy)

	var total int64
	if err := queryBuilder.Count(&total).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error counting questions: %v", err),
			http.StatusInternalServerError)
		return
	}

	offset := (page - 1) * limit
	var questions []models.Question
	result := queryBuilder.Offset(offset).Limit(limit).Find(&questions)

	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching questions: %v", result.Error),
			http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data": questions,
		"pagination": map[string]interface{}{
			"total":        total,
			"per_page":     limit,
			"current_page": page,
			"last_page":    int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// GET Question
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

	query := r.URL.Query()
	detail := query.Get("detail")
	fmt.Println(detail)

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
