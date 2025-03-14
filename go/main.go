package main

import (
	"fmt"
	
	"Deepseek-Go/config"
	"Deepseek-Go/router"
)

func main() {
	config.InitConfig()

	router := router.InitRouter()
	router.Run(fmt.Sprintf(":%d", config.Config.App.Port))
}
