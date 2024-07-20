package migrations

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) {
	log.Println("Starting database migrations...")

	CreateUserTable(db)

	// Add calls to other table migrations here
	// e.g., CreateOtherTable(db)

	log.Println("Database migrations completed.")
}
