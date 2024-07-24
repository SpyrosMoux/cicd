package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"spyrosmoux/api/internal/handlers"
	"spyrosmoux/api/internal/repositories"
	"spyrosmoux/api/internal/services"
)

func SetupUsersRouter(db *gorm.DB, router *gin.Engine) {
	usersRepository := repositories.NewUsersRepositoryImpl(db)
	usersService := services.NewUsersServiceImpl(usersRepository)
	usersHandler := handlers.NewUsersHandler(usersService)

	router.POST("/users", usersHandler.CreateUser)
	router.GET("/users", usersHandler.FindAllUsers)
	router.GET("/users/:id", usersHandler.FindUserById)
	router.DELETE("/users/:id", usersHandler.DeleteUser)
}
