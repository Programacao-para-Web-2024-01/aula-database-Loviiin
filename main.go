package main

import (
	"aula-database/db"
	"database/sql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
)

func main() {
	if err := createDatabase(); err != nil {
		log.Fatal(err)
	}

	if err := createServer(); err != nil {
		log.Fatal(err)
	}
}

func createDatabase() error {
	db, err := sql.Open("mysql", "root:1009@tcp(localhost:3306)/")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS web")
	if err != nil {
		return err
	}

	return nil
}

func connectDB() *sql.DB {
	config := mysql.NewConfig()
	config.User = "root"
	config.Passwd = "1009"
	config.DBName = "web"
	conn, err := mysql.NewConnector(config)
	if err != nil {
		panic(err)
	}
	return sql.OpenDB(conn)
}

func createServer() error {
	studentRepository := db.NewStudentRepository(connectDB())

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/students",
		func(w http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case "GET":
				students, err := studentRepository.List()
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}

				err = json.NewEncoder(w).Encode(students)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}

				w.WriteHeader(200)
			case "POST":
				// Leitura do corpo (INPUT)
				var student db.Student
				err := json.NewDecoder(req.Body).Decode(&student)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}

				// Lógica da função/Algoritmo
				id, err := studentRepository.Create(student)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}

				student.Id = id

				// Output / Resposta
				err = json.NewEncoder(w).Encode(student)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			default:
				http.Error(w, "method not supported", 405)
				return
			}

		})

	mux.HandleFunc(
		"/students/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			// Input
			idRaw := req.PathValue("id")

			id, err := strconv.Atoi(idRaw)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			switch req.Method {
			case "GET":
				// processamento
				student, err := studentRepository.Get(id)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}

				// output
				err = json.NewEncoder(w).Encode(student)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			case "PUT":
				var student db.Student
				err := json.NewDecoder(req.Body).Decode(&student)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}

				err = studentRepository.Update(id, student)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}

				updatedStudent, err := studentRepository.Get(id)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}

				err = json.NewEncoder(w).Encode(updatedStudent)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			case "DELETE":
				err := studentRepository.Delete(id)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				w.WriteHeader(204)
			default:
				http.Error(w, "method not supported", 405)
				return
			}

		})

	return http.ListenAndServe("localhost:8080", mux)
}
