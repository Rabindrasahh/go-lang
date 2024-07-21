package model

import (
	"database/sql"
	"fmt"

	"rest-api/helper"
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

func CreateUser(db *sql.DB, user User) (User, error) {
	hashedPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %v", err)
	}
	user.Password = hashedPassword

	query := "INSERT INTO users (name, email, class, password) VALUES ($1, $2, $3, $4) RETURNING id"
	err = db.QueryRow(query, user.Name, user.Email, user.Class, user.Password).Scan(&user.ID)
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}
