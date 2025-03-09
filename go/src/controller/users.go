package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Deepseek-Go/global"
	"Deepseek-Go/models"
)

// 获取用户信息
func GetUser(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var user models.User
	global.DB.Where("username = ?", username).First(&user)
	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"email":    user.Email,
	})
}
