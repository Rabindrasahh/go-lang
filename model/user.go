package model

import (
	"database/sql"
	"fmt"
)

// User represents a user in the database
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Class string `json:"class"`
}

// GetUserByID retrieves a user by ID from the database
func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := "SELECT id, name, email, class FROM users WHERE id = $1"
	row := db.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Class)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User not found")
		}
		return nil, fmt.Errorf("Error querying user: %v", err)
	}

	return &user, nil
}
