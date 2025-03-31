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

	
	/*r.HandleFunc("/user/create", handleUserCreate).Methods("POST")
	r.HandleFunc("/user/login", handleUserLogin).Methods("POST")
	r.HandleFunc("/user/login", handleUserFind).Methods("GET")
	r.HandleFunc("/register/{key}", validadeRegisterKey).Methods("GET")
	r.HandleFunc("/user/checkToken", checkToken).Methods("POST")
	r.HandleFunc("/products/get", getAllProducts).Methods("GET")
	r.HandleFunc("/product/get/{id}", getOneProduct).Methods("GET")
	r.HandleFunc("/cart/insert", handleCartInsert).Methods("POST")
	r.HandleFunc("cart/delete/{userID}/{productID}", handleCartDeleteOne).Methods("POST")
	r.HandleFunc("cart/delete/{userID}", handleCartDeleteAll).Methods("DELETE")*/
}
