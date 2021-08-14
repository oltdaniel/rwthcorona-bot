package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/oltdaniel/rwthcorona-bot/utils"
	"github.com/oltdaniel/rwthcorona-bot/workers"
)

type WebhookController Controller

func (w *WebhookController) PostMessage(c *gin.Context) {
	var update utils.TelegramWebhookUpdate
	err := c.BindJSON(&update)
	if err != nil {
		c.String(400, "something is wrong with that")
		return
	}
	c.AbortWithStatus(200)
	workers.TelegramResponderQueue <- update
}
