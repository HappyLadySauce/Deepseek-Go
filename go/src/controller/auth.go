package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"Deepseek-Go/global"
	"Deepseek-Go/models"
	"Deepseek-Go/utils/auth"
	"Deepseek-Go/utils/email"
)

// 注册
func Register(c *gin.Context) {
	// 获取请求参数
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名
	if !auth.CheckUsername(user.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名格式错误"})
		return
	}

	// 检查用户是否存在
	if global.DB.Where("username = ?", user.Username).First(&user).RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
		return
	}

	// 验证邮箱格式和可用性
	if err := email.ValidateEmail(user.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查邮箱是否已验证
	if err := email.CheckEmailVerified(user.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 加密密码
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	user.Password = hashedPassword

	// 设置邮箱已验证状态
	user.EmailVerified = true

	// JWT 生成token
	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token生成失败"})
		return
	}

	// 创建用户
	if err := global.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}

	// 返回token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// 登录
func Login(c *gin.Context) {
	// 用于接收登录请求的数据
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从数据库获取用户信息
	var user models.User
	if err := global.DB.Where("username = ?", loginRequest.Username).Or("email = ?", loginRequest.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询错误"})
		}
		return
	}

	// 验证密码
	if !auth.CheckPassword(loginRequest.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码错误"})
		return
	}

	// JWT 生成token
	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token生成失败"})
		return
	}

	// 返回token
	c.JSON(http.StatusOK, gin.H{"token": token})
}
