package app

import (
	"fmt"
	"log"
	"manage_system/config"
	"os"
)

func init() {
	logFile, err := os.OpenFile(config.LOG_PATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lmicroseconds | log.Ldate)
}
