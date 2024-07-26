package migrations

import (
	"database/sql"
	"log"
)

func CreateUserTable(db *sql.DB) {
	// Drop the existing trigger if it exists
	_, err := db.Exec(`
		DROP TRIGGER IF EXISTS update_users_modtime ON users;
	`)
	if err != nil {
		log.Fatalf("Error dropping existing trigger: %v", err)
	}

	// Drop the existing function if it exists
	_, err = db.Exec(`
		DROP FUNCTION IF EXISTS update_users_modtime;
	`)
	if err != nil {
		log.Fatalf("Error dropping existing function: %v", err)
	}

	// Create the users table
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			user_type_id INT NOT NULL,
			password VARCHAR(100) NOT NULL,
			email_verification_token VARCHAR(255) NOT NULL DEFAULT '',
			is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_type_id) REFERENCES user_type(id)
		)
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	// Create the trigger function
	_, err = db.Exec(`
		CREATE OR REPLACE FUNCTION update_users_modtime() RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = NOW();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		log.Fatalf("Error creating trigger function: %v", err)
	}

	// Create the trigger
	_, err = db.Exec(`
		CREATE TRIGGER update_users_modtime
		BEFORE UPDATE ON users
		FOR EACH ROW
		EXECUTE FUNCTION update_users_modtime();
	`)
	if err != nil {
		log.Fatalf("Error creating trigger: %v", err)
	}

	log.Println("Users table and trigger created successfully")
}
