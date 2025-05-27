package handlers

import (
	"encoding/json"
	"errors"
	"flashquest/database"
	"flashquest/pkg/models"
	"fmt"
	"log"
	"net/http"

	filters "flashquest/pkg"

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

	// Paginação
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Ordenação
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

	for param, handler := range filters.QuestionFilters {
		if value := query.Get(param); value != "" {
			queryBuilder = handler(value, queryBuilder)
		}
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

	switch detail {
	case "full":
		handleFullDetail(w, db, id, question)
	case "information":
		handleInformationDetail(w, db, id, question)
	default:
		handleDefaultDetail(w, question)
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



// Question functions //

func handleFullDetail(w http.ResponseWriter, db *gorm.DB, id string, question models.Question) {
	// Define a safe user response struct without sensitive fields
	type SafeUser struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	// Define the full question response struct
	type QuestionDetail struct {
		ID        uint             `json:"id"`
		CreatedAt time.Time        `json:"created_at"`
		UpdatedAt time.Time        `json:"updated_at"`
		Statement string           `json:"statement"`
		Subject   *models.Subject  `json:"subject,omitempty"`
		Topic     *models.Topic    `json:"topic,omitempty"`
		User      *SafeUser        `json:"user,omitempty"`
		Source    *models.Source   `json:"source,omitempty"`
		Answers   []models.Answer  `json:"answers"`
		Comments  []models.Comment `json:"comments"`
	}

	// Get related data
	subject, subjectErr := getQuestionSubject(db, question.SubjectID)
	topic, topicErr := getQuestionTopic(db, id)
	source, sourceErr := getQuestionSource(db, id)
	answers, _ := getQuestionAnswers(db, id)
	comments, _ := getQuestionComments(db, id)
	user, userErr := getQuestionUser(db, question.UserID)

	// Prepare the full response
	fullResponse := QuestionDetail{
		ID:        question.ID,
		CreatedAt: question.CreatedAt,
		UpdatedAt: question.UpdatedAt,
		Statement: question.Statement,
		Answers:   answers,
		Comments:  comments,
	}

	// Only include subject if it was found
	if subjectErr == nil {
		fullResponse.Subject = subject
	}

	// Only include topic if it was found
	if topicErr == nil {
		fullResponse.Topic = topic
	}

	// Only include source if it was found AND has metadata
	if sourceErr == nil && len(source.Metadata) > 0 {
		fullResponse.Source = source
	}

	// Only include user if it was found (without sensitive fields)
	if userErr == nil {
		fullResponse.User = &SafeUser{
			ID:   user.ID,
			Name: user.Name,
		}
	}

	fmt.Printf("Found question with full detail %s \n", id)
	sendResponse(w, fullResponse)
}

func handleInformationDetail(w http.ResponseWriter, db *gorm.DB, id string, question models.Question) {
	// Define a safe user response struct without sensitive fields
	type SafeUser struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	// Define the information question response struct
	type QuestionInfo struct {
		ID        uint            `json:"id"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedAt time.Time       `json:"updated_at"`
		Statement string          `json:"statement"`
		Subject   *models.Subject `json:"subject,omitempty"`
		Topic     *models.Topic   `json:"topic,omitempty"`
		User      *SafeUser       `json:"user,omitempty"`
		Source    *models.Source  `json:"source,omitempty"`
	}

	// Get related data (excluding answers and comments)
	subject, subjectErr := getQuestionSubject(db, question.SubjectID)
	topic, topicErr := getQuestionTopic(db, id)
	source, sourceErr := getQuestionSource(db, id)
	user, userErr := getQuestionUser(db, question.UserID)

	// Prepare the information response
	infoResponse := QuestionInfo{
		ID:        question.ID,
		CreatedAt: question.CreatedAt,
		UpdatedAt: question.UpdatedAt,
		Statement: question.Statement,
	}

	// Only include subject if it was found
	if subjectErr == nil {
		infoResponse.Subject = subject
	}

	// Only include topic if it was found
	if topicErr == nil {
		infoResponse.Topic = topic
	}

	// Only include source if it was found AND has metadata
	if sourceErr == nil && len(source.Metadata) > 0 {
		infoResponse.Source = source
	}

	// Only include user if it was found (without sensitive fields)
	if userErr == nil {
		infoResponse.User = &SafeUser{
			ID:   user.ID,
			Name: user.Name,
		}
	}

	fmt.Printf("Found question with information detail %s \n", id)
	sendResponse(w, infoResponse)
}

func handleDefaultDetail(w http.ResponseWriter, question models.Question) {
	fmt.Printf("Found question %s \n", question.ID)
	sendResponse(w, question)
}