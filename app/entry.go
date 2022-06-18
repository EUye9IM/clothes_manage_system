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
	engine *gin.Engine
	store  = sessions.NewCookieStore([]byte(config.SESSION_SEC))
)

func Run() {
	if err := dbconn.Connect(config.DB_IP, config.DB_PORT, config.DB_DATABASE, config.DB_USER, config.DB_PASSWORD); err != nil {
		log.Fatalln("Load failed.", err)
	}
	defer dbconn.Close()

	gin.SetMode(gin.ReleaseMode)
	engine = gin.Default()
	engine.SetTrustedProxies(nil)

	engine.Static("/static", "./resource/static")
	engine.StaticFile("/favicon.ico", "./resource/favicon.ico")

	engine.Delims("[[[", "]]]")
	engine.LoadHTMLGlob("./resource/templates/**/*")

	routeFront(engine)
	routeBack(engine)
	routeApi(engine)
	log.Println("Load success.")
	engine.Run(fmt.Sprintf(":%v", config.APP_PORT))
}
