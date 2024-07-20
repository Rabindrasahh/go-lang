package db

import (
	"database/sql"
	"log"
	"os"
	"rest-api/database/migrations"
	"rest-api/database/seeders"

	_ "github.com/lib/pq"
)

var Conn *sql.DB

func Init() {
	logDir := "var/log"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0777)
		if err != nil {
			log.Fatalf("Error creating log directory %s: %v", logDir, err)
		}
		log.Println("Log directory created")
	}

	logFilePath := "var/log/go.log"
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	log.SetOutput(logFile)
	databaseURL := "user=postgres password=root dbname=test sslmode=disable"

	// Initialize the database connection
	Conn, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	err = Conn.Ping()
	if err != nil {
		log.Fatalf("Unable to verify connection: %v", err)
	}

	log.Println("Successfully connected to the database")

	// Run migrations and seeders
	migrations.RunMigrations(Conn)
	seeders.RunSeeders(Conn)
}

func Close() {
	if Conn != nil {
		if err := Conn.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}
