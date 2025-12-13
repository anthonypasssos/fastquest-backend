package handlers

import (
	"encoding/json"
	"errors"
	"flashquest/database"
	"flashquest/pkg/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	filters "flashquest/pkg"
	"time"

	"math"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type SafeUser struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type QuestionInput struct {
	Statement            string `json:"statement"`
	SubjectID            int    `json:"subject_id"`
	UserID               int    `json:"user_id"`
	SourceExamInstanceID *int   `json:"source_exam_instance_id"`
}

func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var questionsToProcess []QuestionInput
	var createdQuestions []models.Question

	var questionArray []QuestionInput
	errArray := json.Unmarshal(body, &questionArray)

	if errArray == nil && len(questionArray) > 0 {
		questionsToProcess = questionArray
	} else {
		var singleQuestion QuestionInput
		errSingle := json.Unmarshal(body, &singleQuestion)

		if errSingle == nil && (singleQuestion.Statement != "" || singleQuestion.UserID != 0) {
			questionsToProcess = []QuestionInput{singleQuestion}
		} else {
			http.Error(w, "Invalid request body format: expected single question object or non-empty array of objects", http.StatusBadRequest)
			return
		}
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	for _, input := range questionsToProcess {
		if input.Statement == "" || input.UserID == 0 {
			if len(questionsToProcess) > 1 {
				http.Error(w, "One or more questions are missing Statement or User ID in the batch request", http.StatusBadRequest)
				return
			}
			http.Error(w, "Statement and User ID are required", http.StatusBadRequest)
			return
		}

		question := models.Question{
			Statement:            input.Statement,
			SubjectID:            input.SubjectID,
			UserID:               input.UserID,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			SourceExamInstanceID: input.SourceExamInstanceID,
		}

		if result := db.Create(&question); result.Error != nil {
			http.Error(w, fmt.Sprintf("Error creating question: %v", result.Error), http.StatusInternalServerError)
			return
		}

		createdQuestions = append(createdQuestions, question)
	}

	if len(createdQuestions) == 1 && len(questionsToProcess) == 1 {
		sendJSON(w, createdQuestions[0], http.StatusCreated)
	} else {
		sendJSON(w, createdQuestions, http.StatusCreated)
	}
}

func SendQuestions(qs ...*models.Question) error {
	for i, q := range qs {
		if q.Statement == "" {
			return fmt.Errorf("question.Statement at index %d cannot be empty", i)
		}
	}

	db := database.GetDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := db.Create(qs).Error; err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

func SendAnswers(a *[]models.Answer) error {
	for i, a := range *a {
		if a.Text == "" {
			return fmt.Errorf("answer.Text at index %d cannot be empty", i)
		}
	}

	db := database.GetDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := db.Create(a).Error; err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

/*
func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Statement string `json:"statement"`
		SubjectID int    `json:"subject_id"`
		UserID    int    `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Statement == "" || input.UserID == 0 {
		http.Error(w, "Statement and User ID are required", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	question := models.Question{
		Statement: input.Statement,
		SubjectID: input.SubjectID,
		UserID:    input.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if result := db.Create(&question); result.Error != nil {
		http.Error(w, fmt.Sprintf("Error creating question: %v", result.Error), http.StatusInternalServerError)
		return
	}

	sendJSON(w, question, http.StatusCreated)
}*/

func GetQuestions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := parseInt(query.Get("page"), 1)
	limit := parseInt(query.Get("perPage"), 10)
	if limit > 100 {
		limit = 100
	}

	orderBy := query.Get("order_by")
	if orderBy == "" {
		orderBy = "created_at desc"
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	queryBuilder := applyFilters(db.Model(&models.Question{}), query)
	queryBuilder = queryBuilder.Order(orderBy)

	var total int64
	if err := queryBuilder.Count(&total).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error counting questions: %v", err), http.StatusInternalServerError)
		return
	}

	includeParam := query.Get("include")
	var includes []string
	if includeParam != "" {
		includes = strings.Split(includeParam, ",")
	}

	queryBuilder = queryBuilder.Scopes(models.ApplyIncludes(includes))

	var questions []models.Question
	offset := (page - 1) * limit
	if result := queryBuilder.Offset(offset).Limit(limit).Find(&questions); result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching questions: %v", result.Error), http.StatusInternalServerError)
		return
	}

	var questionsResp []models.QuestionResponse
	for _, q := range questions {
		questionsResp = append(questionsResp, q.ToResponse())
	}

	sendPaginatedResponse(w, questionsResp, total, limit, page)
}

func GetQuestion(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	includeParam := query.Get("include")
	var includes []string
	if includeParam != "" {
		includes = strings.Split(includeParam, ",")
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	var question models.Question
	if result := db.Scopes(models.ApplyIncludes(includes)).Where("id = ?", id).First(&question); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Question not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Error fetching question: %v", result.Error), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(question.ToResponse()); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func GetQuestionsByArray(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	includeParam := query.Get("include")
	var includes []string
	if includeParam != "" {
		includes = strings.Split(includeParam, ",")
	}

	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.IDs) == 0 {
		http.Error(w, "Invalid JSON body or empty IDs array", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	var questions []models.Question

	if result := db.Scopes(models.ApplyIncludes(includes)).Where("id IN ?", req.IDs).Find(&questions); result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching questions: %v", result.Error), http.StatusInternalServerError)
		return
	}

	var questionsResp []models.QuestionResponse
	for _, q := range questions {
		questionsResp = append(questionsResp, q.ToResponse())
	}

	sendJSON(w, questionsResp, http.StatusOK)
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
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
		http.Error(w, fmt.Sprintf("Error deleting question: %v", result.Error), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	log.Printf("Question ID %s deleted successfully", id)
}

func applyFilters(query *gorm.DB, params map[string][]string) *gorm.DB {
	for param, handler := range filters.QuestionFilters {
		if values, exists := params[param]; exists && len(values) > 0 && values[0] != "" {
			query = handler(values[0], query)
		}
	}
	return query
}

func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func sendPaginatedResponse(w http.ResponseWriter, data interface{}, total int64, limit, page int) {
	response := map[string]interface{}{
		"data": data,
		"pagination": map[string]interface{}{
			"total":        total,
			"per_page":     limit,
			"current_page": page,
			"last_page":    int(math.Ceil(float64(total) / float64(limit))),
		},
	}
	sendJSON(w, response, http.StatusOK)
}

func parseInt(value string, defaultValue int) int {
	if result, err := strconv.Atoi(value); err == nil && result > 0 {
		return result
	}
	return defaultValue
}
