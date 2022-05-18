package app

import (
	"manage_system/dbconn"
	"net/http"

	"github.com/gin-gonic/gin"
)

func routeFront(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/index")
	})
	e.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "front/index.tmpl", nil)
	})
	e.GET("/q", func(c *gin.Context) {
		id := c.Query("item_id")
		if id == "" {
			c.Redirect(http.StatusMovedPermanently, "/")
		}
		c.Redirect(http.StatusMovedPermanently, "/q/"+id)
	})
	e.GET("/q/:id", func(c *gin.Context) {
		id := c.Param("id")
		data, err := dbconn.GetItemInfomation(id)
		errstr := ""
		if err != nil {
			errstr = err.Error()
		}
		c.HTML(http.StatusOK, "front/search.tmpl", gin.H{
			"id":   id,
			"data": data,
			"err":  errstr,
		})
	})
	e.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test/test.tmpl", nil)
	})
}

func routeApi(e *gin.Engine) {
	route_group := e.Group("/api")
	route_group.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}
