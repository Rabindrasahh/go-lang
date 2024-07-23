package controller

import (
	"fmt"
	"log"
	"net/http"
	"rest-api/auth"
)

func (uc *UserController) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	log.Printf("Received Authorization header: %s", tokenString)

	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := auth.ParseToken(tokenString)
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	log.Printf("Claims: %v", claims)

	userID, ok := claims["sub"].(float64)
	if !ok {
		log.Printf("Error retrieving user ID from claims: %v", claims)
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("Profile data for user ID: " + fmt.Sprintf("%v", userID)))
}
