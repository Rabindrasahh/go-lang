package migrations

import (
	"database/sql"
	"log"
)

func CreateUserTypeTable(db *sql.DB) {
	query := `
        CREATE TABLE IF NOT EXISTS user_type (
            id SERIAL PRIMARY KEY,
            type VARCHAR(50) NOT NULL,
            description VARCHAR(255) NOT NULL
        )
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating usertype table: %v", err)
	}

	log.Println("Usertype table created successfully")
}
