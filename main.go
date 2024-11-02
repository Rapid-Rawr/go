package main

import (
	"database/sql"
	"go-crud/database"
	"go-crud/handlers"
	"go-crud/middleware"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/golang_test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.InitMigration(db)

	r := mux.NewRouter()
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) { handlers.Register(w, r, db) }).Methods("POST")
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) { handlers.Login(w, r, db) }).Methods("POST")

	auth := r.PathPrefix("/user").Subrouter()
	auth.Use(middleware.Auth)
	auth.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) { handlers.GetUser(w, r, db) }).Methods("GET")

	log.Println("Server running on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
