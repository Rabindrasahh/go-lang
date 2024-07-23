package route

import (
	"net/http"
	"rest-api/controller"
	"rest-api/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, uc *controller.UserController) {
	// Public routes
	r.HandleFunc("/users", uc.CreateUserHandler).Methods("POST")
	r.HandleFunc("/verify", uc.VerifyEmailHandler).Methods("GET")
	r.HandleFunc("/login", uc.LoginHandler).Methods("POST")

	// Protected routes
	protectedRoutes := r.PathPrefix("/protected").Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware)
	protectedRoutes.HandleFunc("/profile", uc.ProfileHandler).Methods("GET")

	// NotFoundHandler for unmatched routes
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}
