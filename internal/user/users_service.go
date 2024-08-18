package user

type UsersService interface {
	Create(user User) *User
	Update(user User)
	Delete(userId string) error
	FindById(userId string) *User
	FindAll() *[]User
}
