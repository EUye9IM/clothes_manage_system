package app

import (
	"encoding/gob"
	"fmt"
	"manage_system/config"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	ID    int
	Name  string
	Grant int
}

func init() {
	gob.Register(UserInfo{})
}

func srcLoc() string {
	_, str, line, _ := runtime.Caller(1)
	return str + ":" + fmt.Sprint(line)
}

func routeBack(e *gin.Engine) {
	route_group := e.Group("/backend")
	//login
	route_group.GET("/login/", func(c *gin.Context) {
		session, err := store.Get(c.Request, config.SESSION_NAME)
		if err != nil {
			c.String(http.StatusOK, "错误。请联系管理。【"+srcLoc()+"】")
		}
		if session.Values["user"] != nil {
			c.Redirect(http.StatusFound, "/backend/")
			return
		}
		c.HTML(http.StatusOK, "back/login.tmpl", nil)
	})

	route_group.GET("/", func(c *gin.Context) {
		session, err := store.Get(c.Request, config.SESSION_NAME)
		if err != nil {
			c.String(http.StatusOK, "错误。请联系管理。【"+srcLoc()+"】")
		}
		if session.Values["user"] == nil {
			c.Redirect(http.StatusFound, "/backend/login/")
			return
		}
		c.HTML(http.StatusOK, "back/index.tmpl", gin.H{
			"user": session.Values["user"].(UserInfo),
		})
	})
	//logout
	route_group.GET("/logout/", func(c *gin.Context) {
		session, err := store.Get(c.Request, config.SESSION_NAME)
		if err != nil {
			c.String(http.StatusOK, "错误。请联系管理。【"+srcLoc()+"】")
		}
		session.Options.MaxAge = -1
		session.Save(c.Request, c.Writer)

		c.Redirect(http.StatusFound, "/backend/")
	})
}
