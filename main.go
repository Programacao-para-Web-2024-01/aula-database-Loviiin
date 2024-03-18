package main

import (
	"aula-database/db"
	"log"
	"net/http"
)

func main() {
	if err := createServer(); err != nil {
		log.Panic(err)
	}
}

func createServer() error {
	studentRepository := db.NewStudentRepository()

	mux := http.NewServeMux()

	// create routes

	return http.ListenAndServe("localhost:8080", mux)
}
