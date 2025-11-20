package handlers

import (
	"context"
	"encoding/json"
	"flashquest/pkg/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"google.golang.org/genai"
)

type TestBody struct {
	Text string `json:"text"`
}

type AIAnswerResponse struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type AIQuestionResponse struct {
	Statement string             `json:"statement"`
	Answers   []AIAnswerResponse `json:"answers"`
}

type AIQuestionSetResponse struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Questions   []AIQuestionResponse `json:"questions"`
}

var aiClient *genai.Client

func InitGemini() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Println("chatGemini error:", err)
		log.Fatal(err)
	}

	aiClient = client
}

func chatGemini(message string) (string, error) {
	ctx := context.Background()

	result, err := aiClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(message),
		nil,
	)
	if err != nil {
		log.Println("chatGemini error:", err)
		return "", err
	}

	return result.Text(), nil
}

func genQuestion(text string) AIQuestionResponse {
	prompt := fmt.Sprintf(`
		Você é um assistente que fala somente em JSON focado em criar questões com 4 alternativas. Não escreva texto normal. Não use markdown. Sempre responda no formato JSON:

		{
			statement: "Texto para questão",
			answers: [
				{
					text: "Alternativa 1",
					is_correct: false
				},
				{
					text: "Alternativa 2",
					is_correct: false
				},
				{
					text: "Alternativa 3",
					is_correct: true
				},
				{
					text: "Alternativa 4",
					is_correct: false
				}
			]
		}

		Seguindo o formato a cima cria uma questão sobre: %s
	`, text)

	var question AIQuestionResponse

	aiResponse, _ := chatGemini(prompt)

	err := json.Unmarshal([]byte(aiResponse), &question)

	if err != nil {
		log.Println("Error on convert response")
	} else {
		log.Println("Successful Generated")
	}

	return question
}

func genQuestionSet(text string) AIQuestionSetResponse {
	prompt := fmt.Sprintf(`
		Você é um assistente que fala somente em JSON focado em criar uma lista de questões com 10 questões com 4 alternativas cada. Não escreva texto normal. Não use markdown. Sempre responda no formato JSON:

		{
			"name": "Nome da lista",
			"description": "Descrição da lista",
			questions: [
				{
					statement: "Texto para questão",
					answers: [
						{
							text: "Alternativa 1",
							is_correct: false
						},
						{
							text: "Alternativa 2",
							is_correct: false
						},
						{
							text: "Alternativa 3",
							is_correct: true
						},
						{
							text: "Alternativa 4",
							is_correct: false
						}
					]
				}
			]
		}

		Seguindo o formato a cima cria uma lista com 10 questões sobre: %s
	`, text)

	var questionSet AIQuestionSetResponse

	aiResponse, _ := chatGemini(prompt)

	err := json.Unmarshal([]byte(aiResponse), &questionSet)

	if err != nil {
		log.Println("Error on convert response")
	} else {
		log.Println("Successful Generated")
	}

	return questionSet
}

func formatQuestions(aiQuestions ...AIQuestionResponse) []models.Question {
	questions := make([]models.Question, 0, len(aiQuestions))

	for _, q := range aiQuestions {
		questions = append(questions, models.Question{
			Statement: q.Statement,
			SubjectID: 7,
			UserID:    5,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return questions
}

func formatAnswer(questionId uint, aiAnswers ...AIAnswerResponse) []models.Answer {
	answers := make([]models.Answer, 0, len(aiAnswers))

	for _, a := range aiAnswers {
		answers = append(answers, models.Answer{
			Text:       a.Text,
			Is_correct: a.IsCorrect,
			QuestionID: questionId,
		})
	}

	return answers
}

func addAIQuestion(aiQuestion AIQuestionResponse) {
	question := formatQuestions(aiQuestion)[0]

	SendQuestions(&question)

	log.Println("Successful Question Insert")

	answers := formatAnswer(question.ID, aiQuestion.Answers...)

	SendAnswers(&answers)
	log.Println("Successful Answer Insert")
}

func addAIQuestionSet(aiQuestionSet AIQuestionSetResponse) {
	questionSet := models.QuestionSet{
		Name:         aiQuestionSet.Name,
		Description:  aiQuestionSet.Description,
		UserID:       5,
		CreationDate: time.Now(),
		IsPrivate:    false,
		Type:         "list",
	}

	SendQuestionSets(&questionSet)

	log.Println("Successful QuestionSet Insert")

	questions := formatQuestions(aiQuestionSet.Questions...)

	slicePtr := make([]*models.Question, 0, len(questions))

	for i := range questions {
		slicePtr = append(slicePtr, &questions[i])
	}

	SendQuestions(slicePtr...)

	log.Println("Successful Question Insert")

	answers := make([]models.Answer, 0, len(aiQuestionSet.Questions)*4)

	for i, q := range aiQuestionSet.Questions {
		answers = append(answers, formatAnswer(questions[i].ID, q.Answers...)...)
	}

	SendAnswers(&answers)

	log.Println("Successful Answer Insert")

	questionSetQuestion := make([]models.QuestionSetQuestion, 0, len(questions))

	for i, q := range questions {
		questionSetQuestion = append(questionSetQuestion, models.QuestionSetQuestion{
			QuestionSetID: questionSet.ID,
			QuestionID:    int(q.ID),
			Position:      i + 1,
		})
	}

	questionSetQuestionPrt := make([]*models.QuestionSetQuestion, 0, len(questionSetQuestion))

	for i := range questionSetQuestion {
		questionSetQuestionPrt = append(questionSetQuestionPrt, &questionSetQuestion[i])
	}

	sendQuestionSetQuestion(questionSetQuestionPrt...)

	log.Println("Successful Answer Insert")
}

func PostAIGenQuestion(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var test TestBody

	errConvert := json.Unmarshal(body, &test)

	if errConvert != nil {
		http.Error(w, "Invalid body", http.StatusInternalServerError)
		return
	}

	log.Println("Successful POST")
	addAIQuestion(genQuestion(test.Text))
}

func PostAIGenQuestionSet(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var test TestBody

	errConvert := json.Unmarshal(body, &test)

	if errConvert != nil {
		http.Error(w, "Invalid body", http.StatusInternalServerError)
		return
	}

	log.Println("Successful POST")
	addAIQuestionSet(genQuestionSet(test.Text))
}
