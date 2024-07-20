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
