package controller

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"rest-api/model"
	"rest-api/service"
	"time"
)

// UserController handles user-related requests.
type UserController struct {
	DB *sql.DB
}

// CreateUserHandler handles POST requests for creating a new user and sending verification email.
func (uc *UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateUserHandler called")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Generate email verification token and set default values
	token, err := generateVerificationToken()
	if err != nil {
		log.Printf("Error generating verification token: %v", err)
		http.Error(w, "Error generating verification token", http.StatusInternalServerError)
		return
	}
	user.EmailVerificationToken = token
	user.IsEmailVerified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Insert user into the database
	createdUser, err := model.CreateUser(uc.DB, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Send verification email
	verificationURL := "http://localhost:8000/verify?token=" + user.EmailVerificationToken
	if err := service.SendVerificationEmail(user.Email, "Email Verification", verificationURL); err != nil {
		log.Printf("Error sending verification email: %v", err)
		http.Error(w, "Error sending verification email", http.StatusInternalServerError)
		return
	}

	// Respond with the created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		log.Printf("Error encoding user data: %v", err)
		http.Error(w, "Error encoding user data", http.StatusInternalServerError)
	}
}

// VerifyEmailHandler handles GET requests for email verification.
func (uc *UserController) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Verification token is required", http.StatusBadRequest)
		return
	}

	// Verify user email using the token
	user, err := model.VerifyUserEmail(uc.DB, token)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		http.Error(w, "Invalid or expired verification token", http.StatusBadRequest)
		return
	}

	if user.IsEmailVerified {
		http.Error(w, "Email is already verified", http.StatusBadRequest)
		return
	}

	user.IsEmailVerified = true
	if err := model.UpdateUser(uc.DB, user); err != nil {
		log.Printf("Error updating user: %v", err)
		http.Error(w, "Error verifying email", http.StatusInternalServerError)
		return
	}

	// Success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email verified successfully"))
}

// generateVerificationToken generates a secure random token for email verification.
func generateVerificationToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
