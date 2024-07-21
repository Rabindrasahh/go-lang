package main

import (
	"log"
	"net/http"
	"os"
	"rest-api/controller"
	"rest-api/route"
	"rest-api/service"
	"rest-api/service/db"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize the database
	db.Init()
	defer db.Close()

	// Create a new router
	r := mux.NewRouter()

	// Register routes
	userController := controller.UserController{DB: db.Conn}
	route.RegisterRoutes(r, &userController)
	// Sending Email Test
	// Sample email data
	emailData := struct {
		Name string
	}{
		Name: "John Doe",
	}

	// Send a welcome email
	err = service.SendHTMLEmail("rabindra.nipunasewa@gmail.com", "Welcome to Our Service", "welcome.html", emailData)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully!")
	// Get the port from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is required")
	}

	// Start the server
	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
