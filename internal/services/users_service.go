package services

import "spyrosmoux/api/internal/models"

type UsersService interface {
	Create(user models.User) *models.User
	Update(user models.User)
	Delete(userId string) error
	FindById(userId string) *models.User
	FindAll() *[]models.User
}
