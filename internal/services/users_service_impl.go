package services

import (
	"spyrosmoux/api/internal/models"
	"spyrosmoux/api/internal/repositories"
)

type UsersServiceImpl struct {
	UsersRepository repositories.UsersRepository
}

func NewUsersServiceImpl(usersRepository repositories.UsersRepository) UsersService {
	return &UsersServiceImpl{UsersRepository: usersRepository}
}

func (u UsersServiceImpl) Create(user models.User) *models.User {
	newUser, err := u.UsersRepository.Save(&user)
	if err != nil {
		return nil
	}

	return newUser
}

func (u UsersServiceImpl) Update(user models.User) {
	//TODO(spyrosmoux) implement me
	panic("implement me")
}

func (u UsersServiceImpl) Delete(userId string) error {
	err := u.UsersRepository.Delete(userId)
	if err != nil {
		return err
	}

	return nil
}

func (u UsersServiceImpl) FindById(userId string) *models.User {
	user, err := u.UsersRepository.FindById(userId)
	if err != nil {
		return nil
	}

	return user
}

func (u UsersServiceImpl) FindAll() *[]models.User {
	users, err := u.UsersRepository.FindAll()
	if err != nil {
		return nil
	}

	return users
}
