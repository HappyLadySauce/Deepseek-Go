package main

import (
	"fmt"
	"Deepseek-Go/config"
	"Deepseek-Go/router"
	// "Deepseek-Go/utils/deepseek"
)

func main() {
	config.InitConfig()

	router := router.InitRouter()
	router.Run(fmt.Sprintf(":%d", config.Config.App.Port))
}
