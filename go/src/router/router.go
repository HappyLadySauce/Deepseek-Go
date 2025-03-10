package router

import (
	"Deepseek-Go/controller"
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
		auth.POST("/login", controller.Login)
		auth.POST("/register", controller.Register)

		// 邮箱验证相关接口
		auth.POST("/send-verification", controller.SendEmailVerification)
		auth.POST("/verify-email", controller.VerifyEmail)
		auth.GET("/test-email-connection", controller.TestEmailConnection)
	}
	api.Use(middlewares.AuthMiddleware())
	{

	}

	return router
}
