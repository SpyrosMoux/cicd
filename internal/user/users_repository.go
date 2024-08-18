package user

type UsersRepository interface {
	Save(user *User) (*User, error)
	Update(user *User)
	Delete(userId string) error
	FindById(userId string) (*User, error)
	FindAll() (*[]User, error)
}
