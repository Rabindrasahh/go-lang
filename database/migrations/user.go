package migrations

import (
	"database/sql"
	"log"
)

func CreateUserTable(db *sql.DB) {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
            class VARCHAR(50) NOT NULL,
            password VARCHAR(100) NOT NULL
        )
    `

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	log.Println("Users table created successfully")
}
