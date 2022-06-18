package app

import (
	"encoding/json"
	"errors"
	"log"
	"manage_system/config"
	"manage_system/dbconn"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

const MAX_LIMIT = 100

func tbToJson(tb *dbconn.Table) (data []gin.H, err error) {
	data = nil
	if tb == nil {
		err = errors.New("bad table")
		return
	}
	col := len(tb.Header)
	for _, i := range tb.Content {
		if len(i) < col {
			err = errors.New("bad table")
			return
		}
		row := make(gin.H)
		for k, v := range tb.Header {
			row[v] = i[k]
		}
		data = append(data, row)
	}
	return
}

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
		c.JSON(http.StatusOK, gin.H{
			"code":  -1,
			"msg":   "您无权限。",
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
	// 用户管理 //////////////////////////////////////////////////////////////////////////////////////////////
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
			tb, err = dbconn.Select("user", nil, nil, []string{"u_id as id", "u_name as name", "u_grant as `grant`"})
		} else {
			tb, err = dbconn.Select("user", []string{"u_name LIKE"}, []string{search_name},
				[]string{"u_id as id", "u_name as name", "u_grant as `grant`"})
		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		data, err := tbToJson(&tb)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】",
				"count": 0,
				"data":  []gin.H{},
			})
			return
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
		if err != nil {
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"res": false,
					"msg": "错误。用户已存在。",
				})
			}
			return
		}
		log.Println("user id:" + strconv.Itoa(uinfo.ID) + " add_user id:" + strconv.Itoa(int(ret)) + " name:" + name + " grant:" + strconv.Itoa(grant))
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
		ret, err := dbconn.Delete("", "user", []string{"u_id ="}, []string{uid})
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

	// 品类管理 //////////////////////////////////////////////////////////////////////////////////////////////
	/*
		select_pattern
		GET
			search_key
		RET JSON
			res
			msg
			count
			data JSON LIST
				id
				name
				brand
				price
	*/

	route_group.GET("/select_pattern/", func(c *gin.Context) {
		ok, _, _ := checkgrant(c, -1)
		if !ok {
			return
		}
		var tb dbconn.Table
		var err error
		search_key := c.Query("search_key")
		if search_key == "" {
			tb, err = dbconn.Select("pattern", nil, nil,
				[]string{"pt_id as id", "pt_name as name", "pt_brand as brand", "pt_price as price"})
		} else {
			tb, err = dbconn.Select("pattern", []string{"pt_name LIKE", "true OR pt_brand LIKE"},
				[]string{search_key, search_key},
				[]string{"pt_id as id", "pt_name as name", "pt_brand as brand", "pt_price as price"})
		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}

		data, err := tbToJson(&tb)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】",
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"msg":   "成功。",
			"count": len(data),
			"data":  data,
		})
	})
	/*
		select_pattern
		POST
			name
			brand
			price
		RET JSON
			code int 0 true -1 false
			msg string
			count int
			data list-json
				id
	*/
	route_group.POST("/add_pattern/", func(c *gin.Context) {
		ok, uinfo, _ := checkgrant(c, dbconn.GRANT_PRODUCT_ADD)
		if !ok {
			return
		}
		var err error
		name := c.PostForm("name")
		brand := c.PostForm("brand")
		price := c.PostForm("price")

		if name == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "品名不能为空",
			})
			return
		}
		if brand == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "品牌不能为空",
			})
			return
		}
		if price == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "价格不能为空",
			})
			return
		}
		ret, err := dbconn.Insert("pattern", []string{"pt_name", "pt_brand", "pt_price"}, []string{name, brand, price})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。价格不合法。",
			})
			return
		}
		log.Println("user id:" + strconv.Itoa(uinfo.ID) + " add_pattern id:" + strconv.Itoa(int(ret)) + " name:" + name + " brand:" + brand + " price:" + price)
		c.JSON(http.StatusOK, gin.H{
			"res": true,
			"msg": "添加成功。",
		})
	})
	/*
		del_pattern
		post
			id
		ret json
			res
			msg
	*/
	route_group.POST("/del_pattern/", func(c *gin.Context) {
		ok, uinfo, _ := checkgrant(c, dbconn.GRANT_DEL)
		if !ok {
			return
		}
		var err error
		id := c.PostForm("id")

		if id == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "id不能为空",
			})
			return
		}
		_, err = dbconn.Delete("item", "item INNER JOIN product ON it_pd_id = pd_id INNER JOIN pattern ON pd_pt_id = pt_id",
			[]string{"pt_id ="}, []string{id})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
			})
			return
		}
		_, err = dbconn.Delete("product", "product INNER JOIN pattern ON pd_pt_id = pt_id",
			[]string{"pt_id ="}, []string{id})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
			})
			return
		}
		_, err = dbconn.Delete("", "pattern",
			[]string{"pt_id ="}, []string{id})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
			})
			return
		}
		log.Println("user id:" + strconv.Itoa(uinfo.ID) + " remove_pattern id:" + id)
		c.JSON(http.StatusOK, gin.H{
			"res": true,
			"msg": "删除成功。",
		})
	})

	// 品类-规格
	/**	select_product
	GET
		ptid ptid
	RET
		code int 0 true -1 false
		msg string
		count int
		data LIST JSON
			id
			SKU
			color
			size
	*/
	route_group.GET("/select_product/", func(c *gin.Context) {
		ok, _, _ := checkgrant(c, -1)
		if !ok {
			return
		}
		var err error
		ptid := c.Query("ptid")

		if ptid == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "ptid不能为空",
			})
			return
		}
		tb, err := dbconn.Select("product", []string{"pd_pt_id ="}, []string{ptid},
			[]string{"pd_id as id", "pd_SKU as SKU", "pd_color as color", "pd_size as size"})

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}

		data, err := tbToJson(&tb)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】",
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"msg":   "成功。",
			"count": len(data),
			"data":  data,
		})
	})
	/**	add_product
	POST
		id ptid
		SKU
		color
		size
	RET JSON
			res
			msg
	*/
	route_group.POST("/add_product/", func(c *gin.Context) {
		ok, uinfo, _ := checkgrant(c, dbconn.GRANT_PRODUCT_ADD)
		if !ok {
			return
		}
		var err error
		ptid := c.PostForm("id")
		SKU := c.PostForm("SKU")
		color := c.PostForm("color")
		size := c.PostForm("size")

		if ptid == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "品类 id 不能为空",
			})
			return
		}

		ret, err := dbconn.Insert("product", []string{"pd_pt_id", "pd_SKU", "pd_color", "pd_size"}, []string{ptid, SKU, color, size})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。" + err.Error(),
			})
			return
		}
		log.Println("user id:" + strconv.Itoa(uinfo.ID) + " add_pattern id:" + strconv.Itoa(int(ret)) + " SKU:" + SKU + " color:" + color + " size:" + size)
		c.JSON(http.StatusOK, gin.H{
			"res": true,
			"msg": "添加成功。",
		})
	})
	/**	del_product
	POST
		id id
	RET JSON
		res
		msg
		count int
		data LIST JSON
			id
			SKU
			color
			size
	*/
	route_group.POST("/del_product/", func(c *gin.Context) {
		ok, uinfo, _ := checkgrant(c, dbconn.GRANT_DEL)
		if !ok {
			return
		}
		var err error
		id := c.PostForm("id")
		if id == "" {
			c.JSON(http.StatusOK, gin.H{
				"res":   false,
				"msg":   "款式 id 不能为空",
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}

		tb, err := dbconn.Select("product", []string{"pd_id ="}, []string{id}, []string{"pd_pt_id as ptid"})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res":   false,
				"msg":   "错误。" + err.Error(),
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		if len(tb.Content) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"res":   false,
				"msg":   "错误。款式 " + id + " 不存在.",
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		ptid := tb.Content[0][0]

		_, err = dbconn.Delete("item", "item INNER JOIN product ON it_pd_id = pd.pd_id",
			[]string{"pt_id ="}, []string{id})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res":   false,
				"msg":   "错误。" + err.Error(),
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		_, err = dbconn.Delete("", "product",
			[]string{"pt_id ="}, []string{id})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res":   false,
				"msg":   "错误。" + err.Error(),
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		log.Println("user id:" + strconv.Itoa(uinfo.ID) + " del_pattern id:" + id)

		tb, err = dbconn.Select("product", []string{"pd_pt_id ="}, []string{ptid},
			[]string{"pd_id as id", "pd_SKU as SKU", "pd_color as color", "pd_size as size"})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  0,
				"msg":   "删除成功，但错误。请联系管理。【" + srcLoc() + "】",
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		data, err := tbToJson(&tb)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  0,
				"msg":   "删除成功，但错误。请联系管理。【" + srcLoc() + "】",
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"res":   true,
			"msg":   "删除成功。",
			"count": len(data),
			"data":  data,
		})
	})

	// 商品管理
	/**	add_batch
	POST
		data JSON
			pdid: count
	RET JSON
			res
			msg
	*/
	route_group.POST("/add_batch/", func(c *gin.Context) {
		ok, uinfo, _ := checkgrant(c, dbconn.GRANT_ITEM_ADD)
		if !ok {
			return
		}
		//	var err error
		raw_data := c.PostForm("data")

		if raw_data == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "data 不能为空",
			})
			return
		}

		m := make(map[string]int)

		err := json.Unmarshal([]byte(raw_data), &m)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "data 不是合法json",
			})
			return
		}
		if len(m) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "data 不能为空",
			})
			return
		}

		ret, err := dbconn.AddItem(uinfo.ID, m)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。" + err.Error(),
			})
			return
		}
		datastr := ""
		for k, v := range m {
			datastr += " " + k + ":" + strconv.Itoa(v)
		}
		log.Println("user id:" + strconv.Itoa(uinfo.ID) + " add_pattern id:" + strconv.Itoa(int(ret)) + " data:" + datastr)
		c.JSON(http.StatusOK, gin.H{
			"res": true,
			"msg": "添加成功。",
		})
	})
	/**	select_batch
	GET
	RET
		code int 0 true -1 false
		msg string
		count int
		data LIST JSON
			id
			uid
			time
			count
	*/
	route_group.GET("/select_batch/", func(c *gin.Context) {
		ok, _, _ := checkgrant(c, -1)
		if !ok {
			return
		}
		var tb dbconn.Table
		var err error
		tb, err = dbconn.SelectEX("batch b LEFT JOIN item i ON b.bt_id = i.it_bt_id", nil, nil,
			[]string{"bt_id as id", "bt_u_id as uid", "bt_time as time", "count(*) as count"}, " GROUP BY bt_id", nil)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}

		data, err := tbToJson(&tb)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】",
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"msg":   "成功。",
			"count": len(data),
			"data":  data,
		})
	})

	/**	del_batch
	POST
		id id
	RET JSON
		res
		msg
	*/
	route_group.POST("/del_batch/", func(c *gin.Context) {
		ok, uinfo, _ := checkgrant(c, dbconn.GRANT_DEL)
		if !ok {
			return
		}
		var err error
		id := c.PostForm("id")
		if id == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "批次 id 不能为空",
			})
			return
		}

		_, _ = dbconn.Delete("", "item", []string{"it_bt_id ="}, []string{id})
		ret, err := dbconn.Delete("", "batch", []string{"bt_id ="}, []string{id})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。" + err.Error(),
			})
			return
		}
		if ret == 0 {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "错误。款式 " + id + " 不存在.",
			})
			return
		}

		log.Println("user id:" + strconv.Itoa(uinfo.ID) + " del_batch id:" + id)

		c.JSON(http.StatusOK, gin.H{
			"res": true,
			"msg": "删除成功。",
		})
	})

	/**	info_batch
	GET
		id id
	RET
		code int 0 true -1 false
		msg string
		count int
		data LIST JSON
			pd_id
			count
	*/
	route_group.GET("/info_batch/", func(c *gin.Context) {
		ok, _, _ := checkgrant(c, -1)
		if !ok {
			return
		}
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusOK, gin.H{
				"res": false,
				"msg": "批次 id 不能为空",
			})
			return
		}
		var tb dbconn.Table
		var err error
		tb, err = dbconn.SelectEX("item JOIN product ON pd_id = it_pd_id JOIN pattern on pt_id = pd_pt_id", []string{"it_bt_id ="}, []string{id},
			[]string{"pd_id", "pt_brand", "pt_name", "pd_color", "pd_size", "count(*) as count"}, " GROUP BY it_pd_id", nil)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}

		data, err := tbToJson(&tb)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】",
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"msg":   "成功。",
			"count": len(data),
			"data":  data,
		})
	})

	/**	select_item
	GET
		page
		limit
		id
		SKU
		batch
	RET
		code int 0 true -1 false
		msg string
		count int
		data LIST JSON
			id
			uid
			time
			count
	*/
	route_group.GET("/select_item/", func(c *gin.Context) {
		var err error
		ok, _, _ := checkgrant(c, -1)
		if !ok {
			return
		}
		page, _ := strconv.Atoi(c.Query("page"))
		limit, err := strconv.Atoi(c.Query("limit"))
		id := c.Query("id")
		SKU := c.Query("SKU")
		batch := c.Query("batch")

		if err != nil || limit > MAX_LIMIT || page <= 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "参数 不合法",
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		tx, err := dbconn.BeginTx()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		defer dbconn.RollbackTx(tx)

		sk := make([]string, 0, 3)
		sv := make([]string, 0, 3)

		if batch != "" {
			sk = append(sk, "bt_id =")
			sv = append(sv, batch)
		}
		if id != "" {
			sk = append(sk, "it_id like")
			sv = append(sv, "%"+id+"%")
		}
		if SKU != "" {
			sk = append(sk, "pd_SKU =")
			sv = append(sv, SKU)
		}

		var tb dbconn.Table
		tb, err = dbconn.TxSelectEX(tx, "item JOIN product ON it_pd_id = pd_id JOIN pattern ON pt_id = pd_pt_id JOIN batch ON bt_id = it_bt_id",
			sk, sv,
			[]string{"sql_calc_found_rows bt_id", "it_id", "pt_name", "pt_brand", "pt_price", "pd_size", "pd_SKU", "pd_color"}, "ORDER BY it_id LIMIT ?, ?",
			[]string{strconv.Itoa((page - 1) * limit), strconv.Itoa(limit)})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}

		data, err := tbToJson(&tb)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】",
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		tb, err = dbconn.TxSelectEX(tx, "", nil, nil, []string{"FOUND_ROWS()"}, "", nil)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  -1,
				"msg":   "错误。请联系管理。【" + srcLoc() + "】" + err.Error(),
				"count": 0,
				"data":  []gin.H{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"msg":   "成功。",
			"count": tb.Content[0][0],
			"data":  data,
		})
	})
}
