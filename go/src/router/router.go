package router

import (
	"Deepseek-Go/controller"
	"Deepseek-Go/global"
	"Deepseek-Go/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// 跨域中间件
	router.Use(middlewares.CORSMiddleware())

	// 初始化控制器
	chatController := controller.NewChatController(global.DB)
	knowledgeController := controller.NewKnowledgeController(global.DB)
	aiConfigController := controller.NewAIConfigController(global.DB)

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

	// 需要认证的接口
	authorized := api.Group("/")
	authorized.Use(middlewares.AuthMiddleware())
	{
		// 聊天相关接口
		chat := authorized.Group("/chat")
		{
			chat.POST("/completions", chatController.Chat)               // 普通聊天
			chat.POST("/stream", chatController.StreamChat)              // 流式聊天
			chat.GET("/sessions", chatController.GetSessions)            // 获取会话列表
			chat.GET("/sessions/:id", chatController.GetSessionMessages) // 获取会话消息
			chat.PUT("/sessions/:id", chatController.UpdateSession)      // 更新会话信息
			chat.DELETE("/sessions/:id", chatController.DeleteSession)   // 删除会话
		}

		// 知识库相关接口
		knowledge := authorized.Group("/knowledge")
		{
			knowledge.POST("/upload", knowledgeController.UploadFile)      // 上传知识库文件
			knowledge.GET("/files", knowledgeController.GetFiles)          // 获取文件列表
			knowledge.GET("/files/:id", knowledgeController.GetFile)       // 获取文件详情
			knowledge.DELETE("/files/:id", knowledgeController.DeleteFile) // 删除文件
		}

		// AI配置相关接口
		aiConfig := authorized.Group("/ai-config")
		{
			aiConfig.POST("/", aiConfigController.CreateConfig)            // 创建配置
			aiConfig.GET("/", aiConfigController.GetConfigs)               // 获取所有配置
			aiConfig.GET("/default", aiConfigController.GetDefaultConfig)  // 获取默认配置
			aiConfig.GET("/models", aiConfigController.GetAvailableModels) // 获取可用模型列表
			aiConfig.GET("/:id", aiConfigController.GetConfig)             // 获取单个配置
			aiConfig.PUT("/:id", aiConfigController.UpdateConfig)          // 更新配置
			aiConfig.DELETE("/:id", aiConfigController.DeleteConfig)       // 删除配置
		}
	}

	return router
}
