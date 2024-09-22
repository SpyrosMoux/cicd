package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init initializes the database connection and runs migrations.
func Init(databaseFile string, models ...interface{}) {
	var err error
	DB, err = gorm.Open(sqlite.Open(databaseFile), &gorm.Config{})
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
