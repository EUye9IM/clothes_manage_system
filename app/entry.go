package app

import (
	"fmt"
	"log"

	"manage_system/config"
	"manage_system/dbconn"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var (
	engin *gin.Engine
	store = sessions.NewCookieStore([]byte(config.SESSION_SEC))
)

func Run() {
	if err := dbconn.Connect(config.DB_IP, config.DB_PORT, config.DB_DATABASE, config.DB_USER, config.DB_PASSWORD); err != nil {
		log.Fatalln("Load failed.", err)
	}
	defer dbconn.Close()

	gin.SetMode(gin.ReleaseMode)
	engin = gin.Default()

	engin.Static("/static", "./resource/static")
	engin.StaticFile("/favicon.ico", "./resource/favicon.ico")

	engin.Delims("[[[", "]]]")
	engin.LoadHTMLGlob("./resource/templates/**/*")

	routeFront(engin)
	routeBack(engin)
	routeApi(engin)
	log.Println("Load success.")
	engin.Run(fmt.Sprintf(":%v", config.APP_PORT))
}
