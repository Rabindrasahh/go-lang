package controller

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"rest-api/auth"
	"rest-api/helper"
	"rest-api/model"
	"rest-api/service"
	"time"
)

type UserController struct {
	DB *sql.DB
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (uc *UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := Response{
			Success: false,
			Message: "Invalid request method",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	var user model.User
	log.Printf("Request Body: %s", r.Body)

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response := Response{
			Success: false,
			Message: "Invalid request payload",
		}
		log.Printf("Error decoding request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	// Generate verification token
	token, err := generateVerificationToken()
	if err != nil {
		response := Response{
			Success: false,
			Message: "Error generating verification token",
		}
		log.Printf("Error generating verification token: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}
	user.EmailVerificationToken = token
	user.IsEmailVerified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	createdUser, err := model.CreateUser(uc.DB, user)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Error creating user",
		}
		log.Printf("Error creating user: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	verificationURL := "http://localhost:8000/verify?token=" + user.EmailVerificationToken
	if err := service.SendVerificationEmail(user.Email, "Email Verification", verificationURL); err != nil {
		response := Response{
			Success: false,
			Message: "Error sending verification email",
		}
		log.Printf("Error sending verification email: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	response := Response{
		Success: true,
		Message: "User created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		response.Success = false
		response.Message = "Error encoding user data"
		log.Printf("Error encoding user data: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

func (uc *UserController) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		response := Response{
			Success: false,
			Message: "Verification token is required",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	user, err := model.VerifyUserEmail(uc.DB, token)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid or expired verification token",
		}
		log.Printf("Error retrieving user: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	if user.IsEmailVerified {
		response := Response{
			Success: false,
			Message: "Email is already verified",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	user.IsEmailVerified = true
	if err := model.UpdateUser(uc.DB, user); err != nil {
		response := Response{
			Success: false,
			Message: "Error verifying email",
		}
		log.Printf("Error updating user: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	response := Response{
		Success: true,
		Message: "Email verified successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
func (uc *UserController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Retrieve the user from the database by email
	user, err := model.GetUserByEmail(uc.DB, req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check if the user exists and the password is correct
	if user.ID == 0 || !helper.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check if the email is verified
	if !user.IsEmailVerified {
		http.Error(w, "Email not verified", http.StatusForbidden)
		return
	}

	// Generate a JWT token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Send the response with the token
	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
func (uc *UserController) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email           string `json:"email"`
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Retrieve the user from the database by email
	user, err := model.GetUserByEmail(uc.DB, req.Email)
	log.Printf("Fetched user: %+v", user)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify current password
	if !helper.CheckPasswordHash(req.CurrentPassword, user.Password) {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	// Validate new password (implement your own validation logic here)
	if len(req.NewPassword) < 8 {
		http.Error(w, "New password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Hash the new password
	hashedPassword, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "Error hashing new password", http.StatusInternalServerError)
		return
	}

	// Update password in the database
	if err := model.UpdateUserPassword(uc.DB, user.ID, hashedPassword); err != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}
	err = service.SendPasswordChangeNotification(user.Email, user.Name)
	if err != nil {
		http.Error(w, "Error sending email notification", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password changed successfully"))
}

func generateVerificationToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
