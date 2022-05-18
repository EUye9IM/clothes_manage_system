package main

import (
	"fmt"
	"manage_system/app"
)

func main() {
	fmt.Println("http://localhost:80/api/ping")
	app.Run()
}
