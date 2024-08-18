package user

import (
	"gorm.io/gorm"
)

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func NewUsersRepositoryImpl(Db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

func (u UsersRepositoryImpl) Save(user *User) (*User, error) {
	result := u.Db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u UsersRepositoryImpl) Update(user *User) {
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

func (u UsersRepositoryImpl) FindById(userId string) (*User, error) {
	user := &User{}
	result := u.Db.Find(&user, "id = ?", userId)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u UsersRepositoryImpl) FindAll() (*[]User, error) {
	users := &[]User{}
	result := u.Db.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
