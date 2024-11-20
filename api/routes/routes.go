package routes

import (
	"github.com/spyrosmoux/cicd/api/config"
	"github.com/spyrosmoux/cicd/api/gh"
	"github.com/spyrosmoux/cicd/api/pipelineruns"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	/* SuperTokens Routers */

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:63342"},
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

	pipelineRunsRepo := pipelineruns.NewRepository(config.DB)
	pipelineRunsSvc := pipelineruns.NewService(pipelineRunsRepo)
	pipelineRunsHandler := pipelineruns.NewHandler(pipelineRunsSvc)
	pipelineruns.Routes(router, pipelineRunsHandler)

	ghService := gh.NewService()
	ghHandler := gh.NewHandler(ghService)
	gh.Routes(router, ghHandler)

	return router
}

// This is a function that wraps the supertokens verification function to work the gin
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
