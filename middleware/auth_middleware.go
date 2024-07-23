package middleware

import (
	"net/http"
	"rest-api/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		_, err := auth.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Optionally, you can store user info in the context
		// ctx := context.WithValue(r.Context(), "userID", claims["sub"])
		// r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
