package repositories

import (
	"spyrosmoux/api/internal/models"

	"gorm.io/gorm"
)

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func NewUsersRepositoryImpl(Db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

func (u UsersRepositoryImpl) Save(user *models.User) (*models.User, error) {
	result := u.Db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u UsersRepositoryImpl) Update(user *models.User) {
	// TODO(spyrosmoux) implement me
	panic("implement me")
}

func (u UsersRepositoryImpl) Delete(userId string) error {
	user, err := u.FindById(userId)
	if err != nil {
		return err
	}

	result := u.Db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u UsersRepositoryImpl) FindById(userId string) (*models.User, error) {
	user := &models.User{}
	result := u.Db.Find(&user, "id = ?", userId)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u UsersRepositoryImpl) FindAll() (*[]models.User, error) {
	users := &[]models.User{}
	result := u.Db.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
