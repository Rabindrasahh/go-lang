package model

import (
	"database/sql"
	"log"
	"rest-api/helper"
	"time"
)

type User struct {
	ID                     int       `json:"id"`
	Name                   string    `json:"name"`
	Email                  string    `json:"email"`
	UserTypeID             int       `json:"user_type_id"`
	Password               string    `json:"password"`
	EmailVerificationToken string    `json:"email_verification_token"`
	IsEmailVerified        bool      `json:"is_email_verified"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

func CreateUser(db *sql.DB, user User) (User, error) {

	hashedPassword, err := helper.HashPassword(user.Password)

	if err != nil {
		return User{}, err
	}

	query := `
        INSERT INTO users (name, email, user_type_id, password, email_verification_token, is_email_verified, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, name, email, user_type_id, password, email_verification_token, is_email_verified, created_at, updated_at
    `

	var createdUser User
	err = db.QueryRow(query, user.Name, user.Email, user.UserTypeID, hashedPassword, user.EmailVerificationToken, user.IsEmailVerified, user.CreatedAt, user.UpdatedAt).
		Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.UserTypeID, &createdUser.Password, &createdUser.EmailVerificationToken, &createdUser.IsEmailVerified, &createdUser.CreatedAt, &createdUser.UpdatedAt)

	if err != nil {
		return User{}, err
	}

	log.Printf("User created successfully: %+v", createdUser)
	return createdUser, nil
}

func VerifyUserEmail(db *sql.DB, token string) (User, error) {
	query := `
        SELECT id, name, email, user_type_id, password, email_verification_token, is_email_verified, created_at, updated_at
        FROM users
        WHERE email_verification_token = $1
    `

	var user User
	err := db.QueryRow(query, token).Scan(&user.ID, &user.Name, &user.Email, &user.UserTypeID, &user.Password, &user.EmailVerificationToken, &user.IsEmailVerified, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, err
	}

	return user, nil
}

func GetUserByEmail(db *sql.DB, email string) (User, error) {
	query := `
        SELECT id, name, email, user_type_id, password, email_verification_token, is_email_verified, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	var user User
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.UserTypeID, &user.Password, &user.EmailVerificationToken, &user.IsEmailVerified, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, err
	}

	return user, nil
}

func UpdateUserPassword(db *sql.DB, userID int, hashedPassword string) error {
	log.Printf("Attempting to update password for user ID: %d", userID)

	// Use PostgreSQL parameter placeholders $1, $2, etc.
	query := "UPDATE users SET password = $1 WHERE id = $2"
	_, err := db.Exec(query, hashedPassword, userID)
	if err != nil {
		log.Printf("Error updating password for user ID %d: %v", userID, err)
		return err
	}

	log.Printf("Password updated successfully for user ID %d", userID)

	return nil
}

func UpdateUser(db *sql.DB, user User) error {
	query := `
        UPDATE users
        SET name = $1, email = $2, user_type_id = $3, password = $4, email_verification_token = $5, is_email_verified = $6, updated_at = $7
        WHERE id = $8
    `

	_, err := db.Exec(query, user.Name, user.Email, user.UserTypeID, user.Password, user.EmailVerificationToken, user.IsEmailVerified, user.UpdatedAt, user.ID)
	return err
}
