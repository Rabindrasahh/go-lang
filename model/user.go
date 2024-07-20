package model

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Class    string `json:"class"`
	Password string `json:"password"`
}

func GetAllUsers(db *sql.DB, page int, pageSize int) ([]User, error) {
	offset := (page - 1) * pageSize
	query := "SELECT id, name, email, class, password FROM users LIMIT $1 OFFSET $2"
	rows, err := db.Query(query, pageSize, offset)
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
	return users, nil
}
