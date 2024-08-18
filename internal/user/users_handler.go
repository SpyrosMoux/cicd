package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UsersHandler struct {
	usersService UsersService
}

func NewUsersHandler(usersService UsersService) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
	}
}

func (usersHandler *UsersHandler) CreateUser(c *gin.Context) {
	var user User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser := usersHandler.usersService.Create(user)
	if createdUser == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": createdUser})
}

func (usersHandler *UsersHandler) UpdateUser(c *gin.Context) {
	// TODO(spyrosmoux) implement this
	panic("implement me")
}

func (usersHandler *UsersHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := usersHandler.usersService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "success")
}

func (usersHandler *UsersHandler) FindUserById(c *gin.Context) {
	id := c.Param("id")
	user := usersHandler.usersService.FindById(id)

	if user == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (usersHandler *UsersHandler) FindAllUsers(c *gin.Context) {
	users := usersHandler.usersService.FindAll()
	if users == nil {
		c.JSON(http.StatusNotFound, nil)
	}

	c.JSON(http.StatusOK, users)
}
