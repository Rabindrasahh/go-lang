package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"rest-api/model"
)

type UserController struct {
	DB *sql.DB
}

func (uc *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetUserHandler called")
	users, err := model.GetAllUsers(uc.DB)
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
