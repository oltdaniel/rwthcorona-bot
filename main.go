package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/oltdaniel/rwth-coronabot/controllers"
	"github.com/oltdaniel/rwth-coronabot/utils"
	"github.com/oltdaniel/rwth-coronabot/workers"
)

func main() {
	// initially set webhook
	err := utils.TelegramSetWebhook(fmt.Sprintf("https://%v:8443/message", os.Getenv("HOSTNAME")), os.Getenv("CERTIFICATE_FILE"))
	if err != nil {
		log.Fatal(err)
	}
	// start background workers
	go workers.TelegramResponder()
	go workers.DataCrawler()
	go workers.DataConverter()
	// start main server
	startApi()
}

func startApi() {
	s := gin.Default()

	webhookController := new(controllers.WebhookController)

	s.POST("/message", webhookController.PostMessage)

	s.RunTLS(":8443", os.Getenv("CERTIFICATE_FILE"), os.Getenv("KEY_FILE"))
}
