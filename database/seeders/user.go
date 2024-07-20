package seeders

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v4"
)

func SeedUsers(db *sql.DB) {
	query := `INSERT INTO users (name, email, class, password) VALUES ($1, $2, $3, $4)`

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 10; i++ {
		_, err := db.Exec(query, faker.Name(), faker.Email(), randomClass(rng), faker.Password())
		if err != nil {
			log.Fatalf("Error seeding user: %v", err)
		}
	}

	log.Println("User seeding completed")
}

func randomClass(rng *rand.Rand) string {
	classes := []string{"Class 1", "Class 2", "Class 3", "Class 4", "Class 5", "Class 6", "Class 7", "Class 8"}
	return classes[rng.Intn(len(classes))]
}
