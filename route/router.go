package route

import (
	"net/http"
	"rest-api/controller"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, uc *controller.UserController) {
	r.HandleFunc("/users", uc.CreateUserHandler).Methods("POST")

	r.HandleFunc("/verify", uc.VerifyEmailHandler).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
}
