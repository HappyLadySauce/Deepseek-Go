package controllers

import (
	"net/http"

	"Deepseek-Go/config"
	"Deepseek-Go/utils/auth"
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// 发送验证邮件
func SendVerificationEmail(c *gin.Context) {
	// 获取请求参数
	type SendVerificationEmailRequest struct {
		Email string `json:"email" binding:"required"`
	}

	// 绑定请求参数
	var sendVerificationEmailRequest SendVerificationEmailRequest
	if err := c.ShouldBindJSON(&sendVerificationEmailRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成验证码
	code := auth.GenerateVerificationCode()

	// 发送验证邮件
	err := auth.SendVerificationEmail(sendVerificationEmailRequest.Email, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录日志
	log.Printf("验证邮件发送成功，邮箱：%s，验证码：%s", sendVerificationEmailRequest.Email, code)

	// 返回验证邮件发送成功
	c.JSON(http.StatusOK, gin.H{"message": "验证邮件发送成功", "email": sendVerificationEmailRequest.Email})
}

// 验证验证码
func VerifyVerificationCode(c *gin.Context) {
	// 获取请求参数
	type VerifyVerificationCodeRequest struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	// 绑定请求参数
	var verifyVerificationCodeRequest VerifyVerificationCodeRequest
	if err := c.ShouldBindJSON(&verifyVerificationCodeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证验证码
	valid, err := auth.VerifyVerificationCode(verifyVerificationCodeRequest.Email, verifyVerificationCodeRequest.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 如果验证码正确，则设置用户状态为已验证
	if valid {
		// 检查用户是否存在
		var user bson.M
		err := config.MongoDB.Database("deepseek").Collection("users").FindOne(
			context.Background(),
			bson.M{"email": verifyVerificationCodeRequest.Email},
		).Decode(&user)

		if err != nil {
			// 如果用户不存在，创建一个新用户记录
			newUser := bson.M{
				"email":          verifyVerificationCodeRequest.Email,
				"email_verified": true,
				"created_at":     time.Now(),
			}

			_, err = config.MongoDB.Database("deepseek").Collection("users").InsertOne(
				context.Background(),
				newUser,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户记录失败"})
				return
			}
		} else {
			// 如果用户已存在，更新其验证状态
			_, err = config.MongoDB.Database("deepseek").Collection("users").UpdateOne(
				context.Background(),
				bson.M{"email": verifyVerificationCodeRequest.Email},
				bson.M{"$set": bson.M{"email_verified": true}},
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户状态失败"})
				return
			}
		}

		// 打印日志
		log.Println("验证验证码成功", verifyVerificationCodeRequest.Email)

		// 返回验证验证码结果
		c.JSON(http.StatusOK, gin.H{"message": "验证验证码成功", "valid": valid})
		return
	}

	// 如果验证码不正确，则返回验证验证码失败
	c.JSON(http.StatusBadRequest, gin.H{"error": "验证验证码失败"})
}

// 注册用户
func Register(c *gin.Context) {
	// 获取请求参数
	type RegisterRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}

	// 绑定请求参数
	var registerRequest RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证邮箱是否已经验证
	var existingUser struct {
		ID            interface{} `bson:"_id"`
		EmailVerified bool        `bson:"email_verified"`
		Username      string      `bson:"username"`
	}

	err := config.MongoDB.Database("deepseek").Collection("users").FindOne(
		context.Background(),
		bson.M{"email": registerRequest.Email},
	).Decode(&existingUser)

	// 如果邮箱不存在，则提示需要先验证邮箱
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请先验证您的邮箱"})
		return
	}

	// 如果邮箱存在但未验证
	if !existingUser.EmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请先验证您的邮箱"})
		return
	}

	// 如果该邮箱对应的用户已经设置了用户名，说明注册已完成
	if existingUser.Username != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该邮箱已经注册过账号"})
		return
	}

	// 检查用户名是否已存在
	var usernameCheck struct {
		Username string `bson:"username"`
	}
	err = config.MongoDB.Database("deepseek").Collection("users").FindOne(
		context.Background(),
		bson.M{"username": registerRequest.Username},
	).Decode(&usernameCheck)

	if err == nil && usernameCheck.Username != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 生成密码哈希
	hashedPassword, err := auth.GeneratePasswordHash(registerRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码哈希失败"})
		return
	}

	// 更新用户记录
	_, err = config.MongoDB.Database("deepseek").Collection("users").UpdateOne(
		context.Background(),
		bson.M{"email": registerRequest.Email},
		bson.M{"$set": bson.M{
			"username":   registerRequest.Username,
			"password":   hashedPassword,
			"updated_at": time.Now(),
		}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	// 生成JWT令牌
	token, err := auth.GenerateJWT(registerRequest.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成JWT令牌失败"})
		return
	}

	// 返回令牌和用户信息
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"token":   token,
	})
}

// 登录用户
func Login(c *gin.Context) {
	// 获取请求参数
	type LoginRequest struct {
		Username string `json:"username" binding:"required"` // 用户名或邮箱
		Password string `json:"password" binding:"required"`
	}

	// 绑定请求参数
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	var user struct {
		ID            interface{} `bson:"_id"`
		Username      string      `bson:"username"`
		Email         string      `bson:"email"`
		Password      string      `bson:"password"`
		EmailVerified bool        `bson:"email_verified"`
	}

	// 尝试通过用户名查找
	err := config.MongoDB.Database("deepseek").Collection("users").FindOne(
		context.Background(),
		bson.M{"username": loginRequest.Username},
	).Decode(&user)

	// 如果通过用户名未找到，尝试通过邮箱查找
	if err != nil {
		err = config.MongoDB.Database("deepseek").Collection("users").FindOne(
			context.Background(),
			bson.M{"email": loginRequest.Username},
		).Decode(&user)

		if err != nil {
			// 如果仍未找到，返回错误
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
	}

	// 检查用户邮箱是否已验证
	if !user.EmailVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先验证您的邮箱"})
		return
	}

	// 验证密码
	if !auth.VerifyPasswordHash(loginRequest.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := auth.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成JWT令牌失败"})
		return
	}

	// 记录登录日志
	log.Printf("用户登录成功，用户名：%s", user.Username)

	// 记录用户最后登录时间
	_, err = config.MongoDB.Database("deepseek").Collection("users").UpdateOne(
		context.Background(),
		bson.M{"username": user.Username},
		bson.M{"$set": bson.M{"last_login": time.Now()}},
	)

	if err != nil {
		log.Printf("更新用户最后登录时间失败：%v", err)
	}

	// 返回登录成功信息
	c.JSON(http.StatusOK, gin.H{
		"message":  "登录成功",
		"token":    token,
	})
}
