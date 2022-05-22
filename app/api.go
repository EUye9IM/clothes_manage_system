package app

import (
	"log"
	"manage_system/config"
	"manage_system/dbconn"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func checkgrant(c *gin.Context, grant int) (ok bool, uinfo UserInfo, session *sessions.Session) {
	ok = false
	session, err := store.Get(c.Request, config.SESSION_NAME)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  -1,
			"msg":   "错误。请联系管理。【" + srcLoc() + "】",
			"count": 0,
			"data":  []gin.H{},
		})
		return
	}
	if session.Values["user"] == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  -1,
			"msg":   "您未登录，请刷新界面。",
			"count": 0,
			"data":  []gin.H{},
		})
		return
	}
	uinfo = session.Values["user"].(UserInfo)
	if grant < 0 {
		ok = true
		return
	}
	if uinfo.Grant&grant == 0 {
		log.Println("user id:" + strconv.Itoa(uinfo.ID) + " have no grant " + strconv.Itoa(grant))
		c.JSON(http.StatusOK, gin.H{
			"code":  -1,
			"msg":   "您无权限，系统将记录您的行为。",
			"count": 0,
			"data":  []gin.H{},
		})
		return
	}
	ok = true
	return
}

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
		ok, uinfo, session := checkgrant(c, -1)
		if !ok {
			return
		}

		oldpw := c.PostForm("old-password")
		newpw := c.PostForm("new-password")

		ok, _, _ = dbconn.Login(uinfo.Name, oldpw)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "密码错误。",
			})
			return
		}
		ret, err := dbconn.SetUserPassword(strconv.Itoa(uinfo.ID), newpw)
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
		log.Println("user id:" + strconv.Itoa(uinfo.ID) + " change own password")
		session.Values["user"] = nil
		session.Save(c.Request, c.Writer)
		c.JSON(http.StatusOK, gin.H{
			"res": true,
			"msg": "修改成功。请重新登录。",
		})
	})
	/*
		select user
		GET
		RET JSON
			code int 0 true -1 false
			msg string
			count int
			data list-json id,name,grant

	*/
	route_group.GET("/select_user/", func(c *gin.Context) {
		ok, _, _ := checkgrant(c, dbconn.GRANT_USER)
		if !ok {
			return
		}
		var tb dbconn.Table
		var err error
		search_name := c.Query("search_name")
		if search_name == "" {
			tb, err = dbconn.Select("user", nil, nil, []string{"u_id", "u_name", "u_grant"})
		} else {
			tb, err = dbconn.Select("user", []string{"u_name LIKE"}, []string{search_name}, []string{"u_id", "u_name", "u_grant"})
		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
				"count": 0,
				"data":  []gin.H{},
			})
		}
		data := []gin.H{}
		for _, i := range tb.Content {
			if len(i) < 3 {
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"code":  -1,
						"msg":   "错误。请联系管理。【" + srcLoc() + "】",
						"count": 0,
						"data":  []gin.H{},
					})
				}
			}
			data = append(data, gin.H{
				"id":    i[0],
				"name":  i[1],
				"grant": i[2],
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"msg":   "成功。",
			"count": len(data),
			"data":  data,
		})
	})
	/*
		adduser
		post
			name
			password
			grant-user
			grant-product-add
			grant-product-edit
			grant-item-read
			grant-item-add
			grant-item-edit
		ret
			res
			msg
	*/
	route_group.POST("/add_user/", func(c *gin.Context) {
		ok, uinfo, _ := checkgrant(c, dbconn.GRANT_USER)
		if !ok {
			return
		}
		var err error
		name := c.PostForm("name")
		password := c.PostForm("password")
		grant, err := strconv.Atoi(c.PostForm("grant"))

		if name == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "用户名不能为空",
			})
			return
		}
		if password == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "密码不能为空",
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "权限标志不合法",
			})
			return
		}
		ret, err := dbconn.AddUser(name, password, grant)
		if err != nil || ret != 1 {
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"res": false,
					"msg": "错误。用户已存在。",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"res": false,
					"msg": "错误。请联系管理。【" + srcLoc() + "】",
				})
			}
			return
		}
		id, _ := dbconn.GetUserID(name)
		log.Println("user id:" + strconv.Itoa(uinfo.ID) + " add_user id:" + id + " name:" + name)
		c.JSON(http.StatusOK, gin.H{
			"res": true,
			"msg": "添加成功。",
		})
	})
	/*
		set_user
		post
			uid
			password
			grant
		ret json
			res
			msg
	*/
	route_group.POST("/set_user/", func(c *gin.Context) {
		ok, uinfo, _ := checkgrant(c, dbconn.GRANT_USER)
		if !ok {
			return
		}
		var err error
		uid := c.PostForm("uid")
		password := c.PostForm("password")
		grant := c.PostForm("grant")

		if uid == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "id不能为空",
			})
			return
		}
		if uid == strconv.Itoa(uinfo.ID) {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "不能修改自身账号",
			})
			return
		}
		if password != "" {
			ret, err := dbconn.SetUserPassword(uid, password)
			if err != nil || ret != 1 {
				c.JSON(http.StatusOK, gin.H{
					"res": false,
					"msg": "uid 不正确",
				})
				return
			}
		}
		_, err = dbconn.Update("user", []string{"u_grant"}, []string{grant}, []string{"u_id ="}, []string{uid})
		if err != nil {
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"res": false,
					"msg": "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
				})
			}
			return
		}
		log.Println("user id:" + strconv.Itoa(uinfo.ID) + " set_user id:" + uid + " grant:" + grant)
		c.JSON(http.StatusOK, gin.H{
			"res": true,
			"msg": "修改成功。",
		})
	})

	/*
		del_user
		post
			uid
		ret json
			res
			msg
	*/
	route_group.POST("/del_user/", func(c *gin.Context) {
		ok, uinfo, _ := checkgrant(c, dbconn.GRANT_USER)
		if !ok {
			return
		}
		var err error
		uid := c.PostForm("uid")

		if uid == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "id不能为空",
			})
			return
		}
		if uid == strconv.Itoa(uinfo.ID) {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "不能删除自身账号",
			})
			return
		}
		ret, err := dbconn.Delete("user", []string{"u_id ="}, []string{uid})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
			})
			return
		}
		if ret != 1 {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。请联系管理。【" + srcLoc() + "】",
			})
			return
		}
		log.Println("user id:" + strconv.Itoa(uinfo.ID) + " remove_user id:" + uid)
		c.JSON(http.StatusOK, gin.H{
			"res": true,
			"msg": "删除成功。",
		})
	})
}
