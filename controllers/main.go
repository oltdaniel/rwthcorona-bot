package controllers

import "github.com/gin-gonic/gin"

type Controller uint8

type MainController Controller

func (m *MainController) GetIndex(c *gin.Context) {
	c.String(200, "pong")
}
