package main

import (
	"net/http"
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
		AllowedOrigins:   []string{"http://localhost:5173", "https://fastquest.vercel.app"}, // Origem permitida
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
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
	r.HandleFunc("/question/create", handlers.CreateQuestion).Methods("POST") //Updated
	r.HandleFunc("/questions", handlers.GetQuestions).Methods("GET")
	r.HandleFunc("/questions/array", handlers.GetQuestionsByArray).Methods("POST")
	r.HandleFunc("/question/{id}", handlers.GetQuestion).Methods("GET")
	r.HandleFunc("/question/{id}/delete", handlers.DeleteQuestion).Methods("DELETE")

	// Answer Requests
	r.HandleFunc("/question/{id}/answer/create", handlers.PostAnswers).Methods("POST")
	r.HandleFunc("/question/{id}/answers", handlers.GetAnswers).Methods("GET")

	//Question Set Requests
	r.HandleFunc("/question-set", handlers.CreateQuestionSet).Methods("POST")
	r.HandleFunc("/question-sets", handlers.GetLists).Methods("GET")
	r.HandleFunc("/question-set/{id}", handlers.GetQuestionSet).Methods("GET")
	r.HandleFunc("/question-set/{id}/questions", handlers.GetQuestionsFromSet).Methods("GET")
	r.HandleFunc("/question-set/{id}/question-ids", handlers.GetQuestionIDsFromSet).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
