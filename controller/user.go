package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"rest-api/config"
	"rest-api/model"
	"strconv"
)

type UserController struct {
	DB *sql.DB
}

func (uc *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetUserHandler called")

	// Parse page number and page size from query parameters
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	page := config.DefaultPage
	pageSize := config.DefaultPageSize

	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)

		if err != nil || page < 1 {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}

	if pageSizeStr != "" {
		var err error
		pageSize, err = strconv.Atoi(pageSizeStr)

		if err != nil || pageSize < 1 {
			http.Error(w, "Invalid page size", http.StatusBadRequest)
			return
		}
	}

	users, err := model.GetAllUsers(uc.DB, page, pageSize)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		http.Error(w, "Error getting users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("Error encoding user data: %v", err)
		http.Error(w, "Error encoding user data", http.StatusInternalServerError)
	}
}

func (uc *UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateUserHandler called")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create user in the database
	createdUser, err := model.CreateUser(uc.DB, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Respond with the created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		log.Printf("Error encoding user data: %v", err)
		http.Error(w, "Error encoding user data", http.StatusInternalServerError)
	}
}
