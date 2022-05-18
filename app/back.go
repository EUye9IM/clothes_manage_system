package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func routeBack(e *gin.Engine) {
	route_group := e.Group("/backend")
	route_group.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "back/login.tmpl", nil)
	})
	route_group.GET("/index", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
