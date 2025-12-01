package main

import (
	"net/http"
	"strings"
	"time"

	//"FlashQuest/database"
	"flashquest/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	_ "flashquest/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewServer() *http.Server {
	r := mux.NewRouter()
	registerPaths(r)

	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		AllowedOrigins:   []string{"https://fastquest.vercel.app"},

		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost")
		},
	})

	handler := c.Handler(r)

	return &http.Server{
		Handler:      handler,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func registerPaths(r *mux.Router) {
	// Question Requests
	r.HandleFunc("/questions", handlers.CreateQuestion).Methods("POST") //Updated
	r.HandleFunc("/questions", handlers.GetQuestions).Methods("GET")
	r.HandleFunc("/questions/by-ids", handlers.GetQuestionsByArray).Methods("POST")
	r.HandleFunc("/questions/{id}", handlers.GetQuestion).Methods("GET")
	r.HandleFunc("/questions/{id}", handlers.DeleteQuestion).Methods("DELETE")

	// Answer Requests
	r.HandleFunc("/questions/{id}/answers", handlers.PostAnswers).Methods("POST")
	r.HandleFunc("/questions/{id}/answers", handlers.GetAnswers).Methods("GET")
	r.HandleFunc("/answers/by-ids", handlers.GetAnswersByIDArray).Methods("POST")

	//Question Set Requests
	r.HandleFunc("/question-sets", handlers.CreateQuestionSet).Methods("POST")
	r.HandleFunc("/question-sets", handlers.GetLists).Methods("GET")
	r.HandleFunc("/question-sets/{id}", handlers.GetQuestionSet).Methods("GET")
	r.HandleFunc("/question-sets/{id}/questions", handlers.GetQuestionsFromSet).Methods("GET")

	//AI requests
	r.HandleFunc("/ai/gen-question", handlers.PostAIGenQuestion).Methods("POST")
	r.HandleFunc("/ai/gen-questionset", handlers.PostAIGenQuestionSet).Methods("POST")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
