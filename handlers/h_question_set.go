package handlers

import (
	"encoding/json"
	"errors"
	"flashquest/database"
	"flashquest/pkg/models"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type NewList struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
	UserID      int    `json:"user_id"`
	Questions   []int  `json:"questions"`
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
		Name:        newList.Name,
		Type:        newList.Type,
		Description: newList.Description,
		UserID:      newList.UserID,
		IsPrivate:   newList.IsPrivate,
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

func SendQuestionSets(qs ...*models.QuestionSet) error {
	for i, a := range qs {
		if a.Name == "" {
			return fmt.Errorf("questionSet.Name at index %d cannot be empty", i)
		}
		if a.Description == "" {
			return fmt.Errorf("questionSet.Description at index %d cannot be empty", i)
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

func sendQuestionSetQuestion(qqs ...*models.QuestionSetQuestion) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := db.Create(qqs).Error; err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

func GetQuestionSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	query := r.URL.Query()
	includeParam := query.Get("include")
	var includes []string
	if includeParam != "" {
		includes = strings.Split(includeParam, ",")
	}

	db := database.GetDB()

	var questionSet models.QuestionSet

	result := db.Scopes(models.ApplyQuestionSetIncludes(includes)).First(&questionSet, id)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching question set: %v", result.Error), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questionSet.ToResponse())
}

func GetQuestionsFromSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db := database.GetDB()

	query := r.URL.Query()
	returnIDs := query.Get("fields") == "id"

	var links []models.QuestionSetQuestion
	result := db.Where("question_set_id = ?", id).Order("position ASC").Find(&links)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching question set links: %v", result.Error), http.StatusInternalServerError)
		return
	}

	// Se não tiver questões associadas
	if len(links) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	// 2. Extrair os IDs das questões
	var questionIDs []int
	for _, link := range links {
		questionIDs = append(questionIDs, link.QuestionID)
	}

	w.Header().Set("Content-Type", "application/json")
	if returnIDs {
		json.NewEncoder(w).Encode(questionIDs)
	} else {
		// 3. Buscar as questões no banco
		var questions []models.Question
		result = db.Where("id IN ?", questionIDs).Find(&questions)
		if result.Error != nil {
			http.Error(w, fmt.Sprintf("Error fetching questions: %v", result.Error), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(questions)
	}
}

/*func GetQuestionIDsFromSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db := database.GetDB()

	// Buscar as relações question_set_question ordenadas pela posição
	var links []models.QuestionSetQuestion
	result := db.Where("question_set_id = ?", id).Order("position ASC").Find(&links)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching question set links: %v", result.Error), http.StatusInternalServerError)
		return
	}

	// Extrair os IDs das questões
	var questionIDs []int
	for _, link := range links {
		questionIDs = append(questionIDs, link.QuestionID)
	}

	// Retornar como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questionIDs)
}*/

func GetLists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("perPage"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	orderBy := query.Get("orderBy")
	allowedOrders := map[string]bool{
		"created_at desc": true,
		"created_at asc":  true,
		"name asc":        true,
		"name desc":       true,
	}
	if !allowedOrders[orderBy] {
		orderBy = "created_at desc"
	}

	db := database.GetDB()
	if db == nil {
		http.Error(w, "Database connection not established", http.StatusInternalServerError)
		return
	}

	queryBuilder := db.Model(&models.QuestionSet{})

	// 🔸 Filtro user_id
	if userId := query.Get("userId"); userId != "" {
		uid, err := strconv.Atoi(userId)
		if err != nil {
			http.Error(w, "Invalid userId", http.StatusBadRequest)
			return
		}
		queryBuilder = queryBuilder.Where("userId = ?", uid)
	}

	// 🔸 Filtro is_private
	if isPrivate := query.Get("isPrivate"); isPrivate != "" {
		private, err := strconv.ParseBool(isPrivate)
		if err != nil {
			http.Error(w, "Invalid isPrivate value", http.StatusBadRequest)
			return
		}
		queryBuilder = queryBuilder.Where("isPrivate = ?", private)
	}

	// 🔍 Filtro de busca por nome ou descrição
	if search := query.Get("statement"); search != "" {
		likeSearch := fmt.Sprintf("%%%s%%", search)
		queryBuilder = queryBuilder.Where(
			"(LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?))",
			likeSearch, likeSearch,
		)
	}

	// Conta total para paginação
	var total int64
	if err := queryBuilder.Count(&total).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error counting lists: %v", err), http.StatusInternalServerError)
		return
	}

	includeParam := query.Get("include")
	var includes []string
	if includeParam != "" {
		includes = strings.Split(includeParam, ",")
	}

	offset := (page - 1) * limit
	var lists []models.QuestionSet
	result := queryBuilder.Scopes(models.ApplyQuestionSetIncludes(includes)).Order(orderBy).Offset(offset).Limit(limit).Find(&lists)

	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching lists: %v", result.Error), http.StatusInternalServerError)
		return
	}

	var responseLists []models.QuestionSetResponse

	for _, qs := range lists {
		responseLists = append(responseLists, qs.ToResponse())
	}

	// Resposta
	response := map[string]interface{}{
		"data": responseLists,
		"pagination": map[string]interface{}{
			"total":        total,
			"per_page":     limit,
			"current_page": page,
			"last_page":    int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
