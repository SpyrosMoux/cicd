package user

type UsersServiceImpl struct {
	UsersRepository UsersRepository
}

func NewUsersServiceImpl(usersRepository UsersRepository) UsersService {
	return &UsersServiceImpl{UsersRepository: usersRepository}
}

func (u UsersServiceImpl) Create(user User) *User {
	newUser, err := u.UsersRepository.Save(&user)
	if err != nil {
		return nil
	}

	return newUser
}

func (u UsersServiceImpl) Update(user User) {
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

func (u UsersServiceImpl) FindById(userId string) *User {
	user, err := u.UsersRepository.FindById(userId)
	if err != nil {
		return nil
	}

	return user
}

func (u UsersServiceImpl) FindAll() *[]User {
	users, err := u.UsersRepository.FindAll()
	if err != nil {
		return nil
	}

	return users
}
