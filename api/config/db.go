package config

import (
	"fmt"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var DB *gorm.DB

// Init initializes the database connection and runs migrations.
func Init(dsn string, models ...interface{}) error {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("unable to connect to database, %w", err)
	}

	// Run migrations
	err = DB.AutoMigrate(models...)
	if err != nil {
		return err
	}
	return nil
}
