package main

import (
	"net/http"
	"time"
	//"FlashQuest/database"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func NewServer() *http.Server {
	r := mux.NewRouter()
	registerPaths(r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"}, // Origem permitida
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
	r.HandleFunc("/questions", getQuestions).Methods("GET")
	r.HandleFunc("/question/{id}", getQuestion).Methods("GET")
	//r.HandleFunc("/question/{id}/delete", deleteQuestion).Methods("DELETE")

	// Answer Requests
	/*r.HandleFunc("/questions/{id}/answer/create", postAnswers).Methods("POST")
	r.HandleFunc("/questions/{id}/answers", getAnswers).Methods("GET")*/
}
