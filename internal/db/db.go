package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"spyrosmoux/api/internal/models"
)

func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=api password=api dbname=api port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Project{}, &models.Repository{}, &models.Pipeline{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
