package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/oltdaniel/rwth-coronabot/utils"
	"github.com/oltdaniel/rwth-coronabot/workers"
)

type WebhookController Controller

func (w *WebhookController) PostMessage(c *gin.Context) {
	var update utils.TelegramWebhookUpdate
	err := c.BindJSON(&update)
	if err != nil {
		c.String(400, "something is wrong with that")
	}
	workers.TelegramResponderQueue <- &update
	c.String(200, "all ok")
}
