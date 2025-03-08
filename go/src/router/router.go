package router

import (
	"Deepseek-Go/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/register", controller.Register)
	}

	return router
}
