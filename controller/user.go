package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"rest-api/model"
	"strconv"
)

// UserController handles user-related requests
type UserController struct {
	DB *sql.DB
}

// GetUserHandler handles GET requests to retrieve a user by ID
func (uc *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := model.GetUserByID(uc.DB, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error encoding user data", http.StatusInternalServerError)
	}
}
