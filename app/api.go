package app

import (
	"manage_system/config"
	"manage_system/dbconn"
	"net/http"

	"github.com/gin-gonic/gin"
)

func routeApi(e *gin.Engine) {
	route_group := e.Group("/api")
	route_group.GET("/ping/", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	route_group.POST("/login/", func(c *gin.Context) {
		session, err := store.Get(c.Request, config.SESSION_NAME)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。请联系管理。【" + srcLoc() + "】",
			})
			return
		}

		if session.Values["user"] != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "您已登录，请刷新界面。",
			})
			return
		}

		name := c.PostForm("username")
		passwd := c.PostForm("password")

		ok, id, grant := dbconn.Login(name, passwd)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "登录失败。用户名或密码错误。",
			})
			return
		}
		session.Values["user"] = UserInfo{id, name, grant}
		session.Save(c.Request, c.Writer)
		c.JSON(http.StatusOK, gin.H{
			"res": true,
			"msg": "登录成功。",
		})
	})

}
