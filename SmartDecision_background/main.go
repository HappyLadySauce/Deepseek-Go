package main

import (
	"fmt"
	
	"SmartDecision/config"
	"SmartDecision/router"
)

func main() {
	config.InitConfig()

	router := router.InitRouter()
	router.Run(fmt.Sprintf(":%d", config.Config.App.Port))
}
