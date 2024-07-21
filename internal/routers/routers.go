package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"gorm.io/gorm"
	"net/http"
	"spyrosmoux/api/internal/handlers"
	"spyrosmoux/api/internal/repositories"
	"spyrosmoux/api/internal/services"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	/* SuperTokens Routers */

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders: append([]string{"content-type"},
			supertokens.GetAllCORSHeaders()...),
		AllowCredentials: true,
	}))

	// Adding the SuperTokens middleware
	router.Use(func(c *gin.Context) {
		supertokens.Middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				c.Next()
			})).ServeHTTP(c.Writer, c.Request)
		// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
		c.Abort()
	})

	/* API Routers */

	// Webhook
	router.POST("/webhook", handlers.HandleWebhook)

	// Projects
	projectsRepository := repositories.NewProjectsRepositoryImpl(db)
	projectsService := services.NewProjectsServiceImpl(projectsRepository)
	projectsHandler := handlers.NewProjectsHandler(projectsService)
	router.POST("/projects", projectsHandler.AddProject)
	router.GET("/projects", projectsHandler.FindAll)
	router.GET("/projects/:id", projectsHandler.FindProjectById)
	router.DELETE("/projects/:id", projectsHandler.Delete)

	// Users
	usersRepository := repositories.NewUsersRepositoryImpl(db)
	usersService := services.NewUsersServiceImpl(usersRepository)
	usersHandler := handlers.NewUsersHandler(usersService)
	router.POST("/users", usersHandler.CreateUser)
	router.GET("/users", usersHandler.FindAllUsers)
	router.GET("/users/:id", usersHandler.FindUserById)
	router.DELETE("/users/:id", usersHandler.DeleteUser)

	return router
}

// This is a function that wraps the supertokens verification function
// to work the gin
func verifySession(options *sessmodels.VerifySessionOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		session.VerifySession(options, func(rw http.ResponseWriter, r *http.Request) {
			c.Request = c.Request.WithContext(r.Context())
			c.Next()
		})(c.Writer, c.Request)
		// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
		c.Abort()
	}
}
