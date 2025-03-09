package router

import (
	"Deepseek-Go/controller"
	"Deepseek-Go/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/register", controller.Register)
	}

	api := router.Group("/api/v1")
	api.Use(middlewares.AuthMiddleware())
	{
		
	}

	return router
}
