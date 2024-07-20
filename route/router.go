package route

import (
	"log"
	"rest-api/controller"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, userController *controller.UserController) {
	log.Fatal(("Hello form the router"))

	router.HandleFunc("/user", userController.GetUserHandler).Methods("GET")
}
