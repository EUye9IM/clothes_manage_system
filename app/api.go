package app

import (
	"manage_system/config"
	"manage_system/dbconn"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func routeApi(e *gin.Engine) {
	route_group := e.Group("/api")
	// test server alive
	route_group.GET("/ping/", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	/*
		login
		post
		username
		password
	*/
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
	/*
		change password
		post
		old-password
		new-password
	*/
	route_group.POST("/changepw/", func(c *gin.Context) {
		session, err := store.Get(c.Request, config.SESSION_NAME)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。请联系管理。【" + srcLoc() + "】",
			})
			return
		}

		if session.Values["user"] == nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "您未登录，请刷新界面。",
			})
			return
		}
		oldpw := c.PostForm("old-password")
		newpw := c.PostForm("new-password")

		ok, _, _ := dbconn.Login(session.Values["user"].(UserInfo).Name, oldpw)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "密码错误。",
			})
			return
		}

		ret, err := dbconn.SetUserPassword(strconv.Itoa(session.Values["user"].(UserInfo).ID), newpw)
		if ret != 1 {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。请联系管理。【" + srcLoc() + "】",
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
			})
			return
		}
		session.Values["user"] = nil
		session.Save(c.Request, c.Writer)
		c.JSON(http.StatusOK, gin.H{
			"res": true,
			"msg": "修改成功。请重新登录。",
		})
		return
	})
}
