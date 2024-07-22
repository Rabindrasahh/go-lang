package model

import (
	"database/sql"
	"log"
	"rest-api/helper" // Import the helper package for hashing
	"time"
)

type User struct {
	ID                     int       `json:"id"`
	Name                   string    `json:"name"`
	Email                  string    `json:"email"`
	Class                  string    `json:"class"`
	Password               string    `json:"-"`
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
        INSERT INTO users (name, email, class, password, email_verification_token, is_email_verified, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, name, email, class, password, email_verification_token, is_email_verified, created_at, updated_at
    `

	var createdUser User
	err = db.QueryRow(query, user.Name, user.Email, user.Class, hashedPassword, user.EmailVerificationToken, user.IsEmailVerified, user.CreatedAt, user.UpdatedAt).
		Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.Class, &createdUser.Password, &createdUser.EmailVerificationToken, &createdUser.IsEmailVerified, &createdUser.CreatedAt, &createdUser.UpdatedAt)

	if err != nil {
		return User{}, err
	}

	// Log the created user including password (ensure this is done only in secure contexts)
	log.Printf("User created successfully: %+v", createdUser)

	return createdUser, nil
}

func VerifyUserEmail(db *sql.DB, token string) (User, error) {
	query := `
        SELECT id, name, email, class, password, email_verification_token, is_email_verified, created_at, updated_at
        FROM users
        WHERE email_verification_token = $1
    `

	var user User
	err := db.QueryRow(query, token).Scan(&user.ID, &user.Name, &user.Email, &user.Class, &user.Password, &user.EmailVerificationToken, &user.IsEmailVerified, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, err
	}

	return user, nil
}

func UpdateUser(db *sql.DB, user User) error {
	query := `
        UPDATE users
        SET name = $1, email = $2, class = $3, password = $4, email_verification_token = $5, is_email_verified = $6, updated_at = $7
        WHERE id = $8
    `

	_, err := db.Exec(query, user.Name, user.Email, user.Class, user.Password, user.EmailVerificationToken, user.IsEmailVerified, user.UpdatedAt, user.ID)
	return err
}
