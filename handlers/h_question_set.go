package handlers

import (
	"encoding/json"
	"flashquest/database"
	"flashquest/pkg/models"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

func GetQuestionSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db := database.GetDB()

	var questionSet models.QuestionSet

	result := db.First(&questionSet, id)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching question set: %v", result.Error), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questionSet)
}

func GetQuestionsFromSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db := database.GetDB()

	// 1. Buscar as relações question_set_question
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

	// 3. Buscar as questões no banco
	var questions []models.Question
	result = db.Where("id IN ?", questionIDs).Find(&questions)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error fetching questions: %v", result.Error), http.StatusInternalServerError)
		return
	}

	// 4. Retornar as questões
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func GetQuestionIDsFromSet(w http.ResponseWriter, r *http.Request) {
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
}
