package model

import (
	"database/sql"
	"log"
)

// Profile represents a user profile.
type Profile struct {
	ID    string
	Name  string
	Email string
}

func GetUserByID(db *sql.DB, userID string) (*Profile, error) {
	var profile Profile
	query := "SELECT id, name, email FROM profiles WHERE id = ?"
	err := db.QueryRow(query, userID).Scan(&profile.ID, &profile.Name, &profile.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println("Error querying profile:", err)
		return nil, err
	}
	return &profile, nil
}
