package main

import (
	"log"
)

func main() {
	srv := NewServer()
	log.Fatal(srv.ListenAndServe())
}

