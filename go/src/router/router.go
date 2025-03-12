package router

import (
	"Deepseek-Go/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 跨域中间件
	router.Use(middlewares.CORSMiddleware())

	api := router.Group("/api/v1")
	auth := api.Group("/auth")
	{
		auth.POST("/login", nil)
	}

	return router
}
