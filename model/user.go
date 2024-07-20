package model

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Class    string `json:"class"`
	Password string `json:"password"`
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	log.Println("GetAllUsers called")
	rows, err := db.Query("SELECT id, name, email, class, password FROM users")
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Class, &user.Password); err != nil {
			return nil, fmt.Errorf("GetAllUsers: %v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}
	log.Println("GetAllUsers completed successfully")
	return users, nil
}
