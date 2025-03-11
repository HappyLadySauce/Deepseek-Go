package middlewares

import (
	"fmt"
	"net/http"

	"Deepseek-Go/global"
	"Deepseek-Go/models"
	"Deepseek-Go/utils/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		// 调试信息
		fmt.Println("接收到的token:", token)

		// 验证JWT令牌
		username, err := auth.ValidateToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌: " + err.Error()})
			c.Abort()
			return
		}

		// 从数据库获取用户信息
		var user models.User
		if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		// 设置用户名和用户ID到上下文
		c.Set("username", username)
		c.Set("userID", user.ID)
		c.Next()
	}
}
