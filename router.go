package main

import (
	"net/http"
	"time"

	//"FlashQuest/database"
	"flashquest/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func NewServer() *http.Server {
	r := mux.NewRouter()
	registerPaths(r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Origem permitida
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	return &http.Server{
		Handler:      handler,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func registerPaths(r *mux.Router) {
	// Question Requests
	r.HandleFunc("/question/create", handlers.CreateQuestion).Methods("POST") //Updated
	r.HandleFunc("/questions", handlers.GetQuestions).Methods("GET")
	r.HandleFunc("/question/{id}", handlers.GetQuestion).Methods("GET")
	r.HandleFunc("/question/{id}/delete", handlers.DeleteQuestion).Methods("DELETE")

	// Answer Requests
	r.HandleFunc("/questions/{id}/answer/create", handlers.PostAnswers).Methods("POST")
	r.HandleFunc("/questions/{id}/answers", handlers.GetAnswers).Methods("GET")
}
