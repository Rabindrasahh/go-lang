package middleware

import (
	"log"
	"net/http"
	"rest-api/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			log.Println("No token provided")
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		_, err := auth.ParseToken(tokenString)
		if err != nil {
			log.Printf("Error parsing token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Println("Token is valid")

		// Optionally, you can store user info in the context
		// ctx := context.WithValue(r.Context(), "userID", claims["sub"])
		// r = r.WithContext(ctx)

		// Pass control to the next handler
		next.ServeHTTP(w, r)
	})
}
