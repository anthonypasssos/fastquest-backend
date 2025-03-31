package main

import (
	"log"
	"FlashQuest/database"
)

func main() {
	database.InitDB()

	srv := NewServer()
	log.Fatal(srv.ListenAndServe())
}

