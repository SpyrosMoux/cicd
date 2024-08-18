package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUsersRouter(db *gorm.DB, router *gin.Engine) {
	usersRepository := NewUsersRepositoryImpl(db)
	usersService := NewUsersServiceImpl(usersRepository)
	usersHandler := NewUsersHandler(usersService)

	router.POST("/users", usersHandler.CreateUser)
	router.GET("/users", usersHandler.FindAllUsers)
	router.GET("/users/:id", usersHandler.FindUserById)
	router.DELETE("/users/:id", usersHandler.DeleteUser)
}
