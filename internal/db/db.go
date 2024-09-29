package db

import (
	"gorm.io/driver/postgres"
	"log"

	"gorm.io/gorm"
)

var DB *gorm.DB

// Init initializes the database connection and runs migrations.
func Init(dsn string, models ...interface{}) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Run migrations
	err = DB.AutoMigrate(models...)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database connection established and migrations ran successfully.")
}
