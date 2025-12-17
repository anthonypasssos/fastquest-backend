package handlers

import (
	"encoding/json"
	"errors"
	"flashquest/database"
	"flashquest/helpers"
	"flashquest/pkg/models"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func SendExamInstance(ei ...*models.ExamInstance) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	if err := db.Create(ei).Error; err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	return nil
}

func GetInstanceWithSource(eiID int, ei *models.ExamInstance) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database connection not established")
	}

	result := db.Preload("SourceExamInstance.Source").Where("id = ?", eiID).Find(&ei)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("Source Exam Instance not found")
		} else {
			return errors.New("Error fetching Source Exam Instance")
		}
	}

	return nil
}

type NewExam struct {
	Exam models.ExamInstance  `json:"exam"`
	List NewListWithQuestions `json:"list"`
}

func CreateExam(w http.ResponseWriter, r *http.Request) {
	var newExam NewExam

	// Decodifica o JSON enviado
	err := json.NewDecoder(r.Body).Decode(&newExam)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		println(err.Error())
		return
	}

	exam := newExam.Exam

	errSendE := SendExamInstance(&exam)

	if errSendE != nil {
		http.Error(w, fmt.Sprintf("Error creating exam: %v", err), http.StatusInternalServerError)
		return
	}

	questionSet := models.QuestionSet{
		Name:        newExam.List.Name,
		Description: newExam.List.Description,
		UserID:      1,
		CreatedAt:   time.Now(),
		IsPrivate:   false,
		Type:        "list",
	}

	errSendQS := SendQuestionSets(&questionSet)

	if errSendQS != nil {
		http.Error(w, fmt.Sprintf("Error creating question set: %v", err), http.StatusInternalServerError)
		return
	}

	var questions []models.Question

	for _, q := range newExam.List.Questions {
		questions = append(questions, models.Question{
			Statement:            q.Statement,
			SubjectID:            q.SubjectID,
			UserID:               1,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			SourceExamInstanceID: &exam.ID,
		})
	}

	errSendQ := SendQuestions(helpers.PtrSlice(questions)...)

	if errSendQ != nil {
		http.Error(w, fmt.Sprintf("Error creating questions: %v", err), http.StatusInternalServerError)
		return
	}

	answers := make([]models.Answer, 0, len(questions)*4)

	for i, q := range newExam.List.Questions {
		for _, a := range *q.Answers {
			answers = append(answers, models.Answer{
				Text:       a.Text,
				Is_correct: a.Is_correct,
				QuestionID: questions[i].ID,
			})
		}
	}

	errSendA := SendAnswers(&answers)

	if errSendA != nil {
		http.Error(w, fmt.Sprintf("Error creating answers: %v", err), http.StatusInternalServerError)
		return
	}

	questionSetQuestion := make([]models.QuestionSetQuestion, 0, len(questions))

	for i, q := range questions {
		questionSetQuestion = append(questionSetQuestion, models.QuestionSetQuestion{
			QuestionSetID: questionSet.ID,
			QuestionID:    int(q.ID),
			Position:      i + 1,
		})
	}

	errSendQSQ := sendQuestionSetQuestion(helpers.PtrSlice(questionSetQuestion)...)

	if errSendQSQ != nil {
		http.Error(w, fmt.Sprintf("Error creating relation between question and question set: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questionSet.ToResponse())
}
