package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/webhooks/v6/github"
	"log"
	"spyrosmoux/api/internal/helpers"
)

func HandleWebhook(c *gin.Context) {
	ghSecret := helpers.LoadEnvVariable("GH_SECRET")
	log.Println(ghSecret)
	hook, _ := github.New(github.Options.Secret(ghSecret))

	payload, err := hook.Parse(c.Request, github.PushEvent)
	if err != nil {
		log.Panicf("Error parsing webhook payload: %v", err)
	}

	log.Printf("Received webhook payload: %+v", payload)
}
