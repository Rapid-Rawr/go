package handlers

import (
	"database/sql"
	"encoding/json"
	"go-crud/models"
	"go-crud/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, hashedPassword)
	if err != nil {
		if err.Error() == "Error 1062: Duplicate entry" { // Gantilah dengan error yang sesuai
			http.Error(w, "User already exists", http.StatusConflict)
		} else {
			http.Error(w, "Error registering user", http.StatusInternalServerError)
		}
		return
	}

	// Successful registration
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var reqUser models.User
	err := json.NewDecoder(r.Body).Decode(&reqUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var dbUser models.User
	row := db.QueryRow("SELECT id, username, password FROM users WHERE username=?", reqUser.Username)
	err = row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error querying user", http.StatusInternalServerError)
		return
	}

	// Check password
	if !utils.CheckPasswordHash(reqUser.Password, dbUser.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := utils.GenerateToken(dbUser.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Successful login
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func GetUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	row := db.QueryRow("SELECT id, username, created_at FROM users WHERE id = ?", userID)
	err = row.Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
