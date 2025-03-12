package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"Deepseek-Go/config"
	"Deepseek-Go/utils/email"
)

// 发送邮箱验证码
func SendEmailVerification(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required"`
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 发送验证码
	if err := email.SendVerificationCode(request.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "验证码发送成功，请在3分钟内验证"})
}

// 验证邮箱
func VerifyEmail(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证邮箱验证码
	if err := email.VerifyEmailCode(request.Email, request.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "邮箱验证成功"})
}

// TestEmailConnection 测试邮件服务器连接
// 这个接口用于测试当前配置的邮件服务器是否可以正常连接
func TestEmailConnection(c *gin.Context) {
	log.Println("开始测试邮件服务器连接...")

	// 输出当前配置信息
	emailCfg := config.Config.Email
	fromAddr := email.ExtractEmailAddress(emailCfg.From)
	log.Printf("邮件配置: 主机=%s, 端口=%d, 用户名=%s, 发件人=%s (提取地址=%s), SSL=%v, 服务器名=%s",
		emailCfg.Host, emailCfg.Port, emailCfg.Username,
		emailCfg.From, fromAddr, emailCfg.EnableSSL, emailCfg.ServerName)

	err := email.TestConnection()
	if err != nil {
		errorMsg := fmt.Sprintf("邮件服务器连接测试失败: %v", err)
		log.Println(errorMsg)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   errorMsg,
			"config": gin.H{
				"host":        emailCfg.Host,
				"port":        emailCfg.Port,
				"username":    emailCfg.Username,
				"from":        emailCfg.From,
				"enable_ssl":  emailCfg.EnableSSL,
				"server_name": emailCfg.ServerName,
			},
		})
		return
	}

	log.Println("邮件服务器连接测试成功")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "邮件服务器连接正常",
		"config": gin.H{
			"host":        emailCfg.Host,
			"port":        emailCfg.Port,
			"username":    emailCfg.Username,
			"from":        emailCfg.From,
			"enable_ssl":  emailCfg.EnableSSL,
			"server_name": emailCfg.ServerName,
		},
	})
}
