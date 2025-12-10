package handlers

import (
	"encoding/json"
	"flashquest/pkg/models"
	"net/http"

	"gorm.io/gorm"
)

// Common database requests //

func getQuestionSubject(db *gorm.DB, subjectID int) (*models.Subject, error) {
	var subject models.Subject
	err := db.Where("id = ?", subjectID).First(&subject).Error
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

func getQuestionTopic(db *gorm.DB, questionID string) (*models.Topic, error) {
	var topic models.Topic
	err := db.Table("question_topic").
		Select("topic.id, topic.name, topic.subject_id").
		Joins("JOIN topic ON question_topic.topic_id = topic.id").
		Where("question_topic.question_id = ?", questionID).
		Scan(&topic).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func getQuestionAnswers(db *gorm.DB, questionID string) ([]models.Answer, error) {
	var answers []models.Answer
	err := db.Where("id_question = ?", questionID).Find(&answers).Error
	if err != nil {
		return nil, err
	}
	return answers, nil
}

func getQuestionComments(db *gorm.DB, questionID string) ([]models.Comment, error) {
	var comments []models.Comment
	err := db.Table("comment_relationship").
		Select("comment.id, comment.text, comment.creation_date, comment.user_id").
		Joins("JOIN comment ON comment_relationship.id_comment = comment.id").
		Where("comment_relationship.type_reference = ? AND comment_relationship.id_reference = ?", "question", questionID).
		Scan(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func getQuestionUser(db *gorm.DB, userID int) (*models.User, error) {
	var user models.User
	err := db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func sendResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
