package repositories

import (
	"spyrosmoux/api/internal/models"
)

type UsersRepository interface {
	Save(user *models.User) (*models.User, error)
	Update(user *models.User)
	Delete(userId string) error
	FindById(userId string) (*models.User, error)
	FindAll() (*[]models.User, error)
}
