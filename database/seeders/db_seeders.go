package seeders

import (
	"database/sql"
	"log"
	"time"
)

func RunSeeders(db *sql.DB) {
	log.Println("Starting database seeding...")

	// Start measuring time
	startTime := time.Now()

	//\\== Comment if not require seeders ==//\\
	// SeedUsers(db)

	// Calculate elapsed time
	elapsedTime := time.Since(startTime)

	log.Printf("Database seeding completed in %s.", elapsedTime)
}
