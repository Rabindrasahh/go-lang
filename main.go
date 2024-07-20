package main

import (
	"log"
	"net/http"
	"os"
	"rest-api/controller"
	"rest-api/route"
	"rest-api/service/db"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	db.Init()
	defer db.Close()

	// Create a new UserController instance
	userController := &controller.UserController{DB: db.Conn}

	// Register routes
	route.RegisterRoutes(userController)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is required")
	}

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
