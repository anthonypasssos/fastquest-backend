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
    if detail == "full" {
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

        // Get the subject for this question
        var subject models.Subject
        subjectErr := db.Where("id = ?", question.SubjectID).First(&subject).Error
        if subjectErr != nil && !errors.Is(subjectErr, gorm.ErrRecordNotFound) {
            http.Error(w, fmt.Sprintf("Error fetching subject: %v", subjectErr),
                http.StatusInternalServerError)
            return
        }

        // Get the topic for this question through question_topic join table
        var topic models.Topic
        topicErr := db.Table("question_topic").
            Select("topic.id, topic.name, topic.subject_id").
            Joins("JOIN topic ON question_topic.topic_id = topic.id").
            Where("question_topic.question_id = ?", id).
            Scan(&topic).Error

        // Get the source for this question
        var source models.Source
        sourceErr := db.Table("question_source").
            Select("source.id, source.name, source.type, source.metadata").
            Joins("JOIN source ON question_source.source_id = source.id").
            Where("question_source.question_id = ?", id).
            Scan(&source).Error

        // Get answers for this question
        var answers []models.Answer
        answersErr := db.Where("id_question = ?", id).Find(&answers).Error
        if answersErr != nil && !errors.Is(answersErr, gorm.ErrRecordNotFound) {
            http.Error(w, fmt.Sprintf("Error fetching answers: %v", answersErr),
                http.StatusInternalServerError)
            return
        }

        // Get comments for this question through comment_relationship
        var comments []models.Comment
        commentsErr := db.Table("comment_relationship").
            Select("comment.id, comment.text, comment.creation_date, comment.user_id").
            Joins("JOIN comment ON comment_relationship.id_comment = comment.id").
            Where("comment_relationship.type_reference = ? AND comment_relationship.id_reference = ?", "question", id).
            Scan(&comments).Error

        if commentsErr != nil && !errors.Is(commentsErr, gorm.ErrRecordNotFound) {
            http.Error(w, fmt.Sprintf("Error fetching comments: %v", commentsErr),
                http.StatusInternalServerError)
            return
        }

        // Get user information
        var user models.User
        userErr := db.Where("id = ?", question.UserID).First(&user).Error
        if userErr != nil && !errors.Is(userErr, gorm.ErrRecordNotFound) {
            http.Error(w, fmt.Sprintf("Error fetching user: %v", userErr),
                http.StatusInternalServerError)
            return
        }

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
            fullResponse.Subject = &subject
        }

        // Only include topic if it was found
        if topicErr == nil {
            fullResponse.Topic = &topic
        }

        // Only include source if it was found AND has metadata
        if sourceErr == nil && len(source.Metadata) > 0 {
            fullResponse.Source = &source
        }

        // Only include user if it was found (without sensitive fields)
        if userErr == nil {
            fullResponse.User = &SafeUser{
                ID:   user.ID,
                Name: user.Name,
            }
        }

        fmt.Printf("Found question with full detail %s \n", id)
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(fullResponse); err != nil {
            http.Error(w, "Error encoding response", http.StatusInternalServerError)
        }
    } else {
        fmt.Printf("Found question %s \n", id)
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(question); err != nil {
            http.Error(w, "Error encoding response", http.StatusInternalServerError)
        }
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
