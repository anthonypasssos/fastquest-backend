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

	// Pagination
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Ordering
	orderBy := query.Get("order_by")
	if orderBy == "" {
		orderBy = "created_at desc"
	}

	// Detail level
	detail := query.Get("detail")

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

	// Handle response based on detail level
	w.Header().Set("Content-Type", "application/json")

	if detail != "" {
		// For detailed responses, we need to handle pagination separately
		switch detail {
		case "full":
			handleFullDetail(w, db, questions)
			return
		case "information":
			infoQuestions := make([]interface{}, len(questions))
			for i, question := range questions {
				infoQuestions[i] = getInformationQuestionDetail(db, question)
			}
		
			response := map[string]interface{}{
				"data": infoQuestions,
				"pagination": map[string]interface{}{
					"total":        total,
					"per_page":     limit,
					"current_page": page,
					"last_page":    int(math.Ceil(float64(total) / float64(limit))),
				},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Default response with pagination
	response := map[string]interface{}{
		"data": questions,
		"pagination": map[string]interface{}{
			"total":        total,
			"per_page":     limit,
			"current_page": page,
			"last_page":    int(math.Ceil(float64(total) / float64(limit))),
		},
	}

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
		handleFullDetail(w, db, question)
	case "information":
		handleInformationDetail(w, db, question)
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

// Updated handleFullDetail to work with either single or multiple questions
func handleFullDetail(w http.ResponseWriter, db *gorm.DB, questions interface{}) {
	switch q := questions.(type) {
	case models.Question:
		// Single question case
		detail := getFullQuestionDetail(db, q)
		sendResponse(w, detail)
	case []models.Question:
		// Multiple questions case
		detailedQuestions := make([]interface{}, len(q))
		for i, question := range q {
			detailedQuestions[i] = getFullQuestionDetail(db, question)
		}
		sendResponse(w, detailedQuestions)
	default:
		http.Error(w, "Invalid question type", http.StatusInternalServerError)
	}
}

// Similarly update handleInformationDetail
func handleInformationDetail(w http.ResponseWriter, db *gorm.DB, questions interface{}) {
	switch q := questions.(type) {
	case models.Question:
		// Single question case
		info := getInformationQuestionDetail(db, q)
		sendResponse(w, info)
	case []models.Question:
		// Multiple questions case
		infoQuestions := make([]interface{}, len(q))
		for i, question := range q {
			infoQuestions[i] = getInformationQuestionDetail(db, question)
		}
		sendResponse(w, infoQuestions)
	default:
		http.Error(w, "Invalid question type", http.StatusInternalServerError)
	}
}

func handleDefaultDetail(w http.ResponseWriter, question models.Question) {
	fmt.Printf("Found question %s \n", question.ID)
	sendResponse(w, question)
}

// Reusable functions from GetQuestion handler
func getFullQuestionDetail(db *gorm.DB, question models.Question) interface{} {
	type SafeUser struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

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

	questionID := strconv.FormatUint(uint64(question.ID), 10)
	subject, _ := getQuestionSubject(db, question.SubjectID)
	topic, _ := getQuestionTopic(db, questionID)
	source, _ := getQuestionSource(db, questionID)
	answers, _ := getQuestionAnswers(db, questionID)
	comments, _ := getQuestionComments(db, questionID)
	user, _ := getQuestionUser(db, question.UserID)

	detail := QuestionDetail{
		ID:        question.ID,
		CreatedAt: question.CreatedAt,
		UpdatedAt: question.UpdatedAt,
		Statement: question.Statement,
		Answers:   answers,
		Comments:  comments,
	}

	if subject != nil {
		detail.Subject = subject
	}
	if topic != nil {
		detail.Topic = topic
	}
	if source != nil && len(source.Metadata) > 0 {
		detail.Source = source
	}
	if user != nil {
		detail.User = &SafeUser{
			ID:   user.ID,
			Name: user.Name,
		}
	}

	return detail
}

func getInformationQuestionDetail(db *gorm.DB, question models.Question) interface{} {
	type SafeUser struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

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

	questionID := strconv.FormatUint(uint64(question.ID), 10)
	subject, _ := getQuestionSubject(db, question.SubjectID)
	topic, _ := getQuestionTopic(db, questionID)
	source, _ := getQuestionSource(db, questionID)
	user, _ := getQuestionUser(db, question.UserID)

	info := QuestionInfo{
		ID:        question.ID,
		CreatedAt: question.CreatedAt,
		UpdatedAt: question.UpdatedAt,
		Statement: question.Statement,
	}

	if subject != nil {
		info.Subject = subject
	}
	if topic != nil {
		info.Topic = topic
	}
	if source != nil && len(source.Metadata) > 0 {
		info.Source = source
	}
	if user != nil {
		info.User = &SafeUser{
			ID:   user.ID,
			Name: user.Name,
		}
	}

	return info
}
