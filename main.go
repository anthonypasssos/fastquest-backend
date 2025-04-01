package main

import (
	"log"
	"flashquest/database"
	"fmt"
)

func main() {
	fmt.Println("Running backend")
	database.InitDB()

	srv := NewServer()
	log.Fatal(srv.ListenAndServe())
}

