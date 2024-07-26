package seeders

import (
	"database/sql"
	"log"
)

func SeedUserTypeTable(db *sql.DB) {
	// Define the user types to be seeded
	userTypes := []struct {
		Type        string
		Description string
	}{
		{"customer", "A regular customer"},
		{"admin", "Administrative user with extended privileges"},
		{"super_admin", "User with full access rights"},
	}

	// Insert each user type into the user_type table
	for _, userType := range userTypes {
		_, err := db.Exec(`
			INSERT INTO user_type (type, description) 
			VALUES ($1, $2);
		`, userType.Type, userType.Description)
		if err != nil {
			log.Fatalf("Error seeding user_type table: %v", err)
		}
	}

	log.Println("User_type table seeded successfully")
}
