package seeders

import (
	"database/sql"
	"log"
)

func RunSeeders(db *sql.DB) {
	log.Println("Starting database seeding...")
	//\\== Comment if not require seeders ==//\\
	// SeedUsers(db)

	log.Println("Database seeding completed.")
}
