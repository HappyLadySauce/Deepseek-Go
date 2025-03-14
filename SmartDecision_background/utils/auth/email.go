package auth

import (
	"SmartDecision/config"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gomail.v2"
)

// 发送邮件
func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Config.Email.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(config.Config.Email.Host, config.Config.Email.Port, config.Config.Email.Username, config.Config.Email.Password)

	if config.Config.Email.EnableSSL {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: config.Config.Email.EnableSSL}
	}

	return d.DialAndSend(m)
}

// 生成验证码
func GenerateVerificationCode() string {
	code := ""
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	return code
}

// 发送验证邮件
func SendVerificationEmail(to, code string) error {
	subject := "Deepseek-Go 验证邮件"
	body := fmt.Sprintf("您的验证码是：%s，请在3分钟内完成验证。", code)

	// 将验证码更新到 MongoDB 中
	config.MongoDB.Database("deepseek").Collection("verification_codes").UpdateOne(context.Background(), bson.M{"email": to}, bson.M{"$set": bson.M{"code": code, "created_at": time.Now()}}, options.Update().SetUpsert(true))

	return SendEmail(to, subject, body)
}

// 验证验证码
func VerifyVerificationCode(email, code string) (bool, error) {
	// 从 MongoDB 中获取验证码
	var verificationCode struct {
		Code string `bson:"code"`
		CreatedAt time.Time `bson:"created_at"`
	}
	
	err := config.MongoDB.Database("deepseek").Collection("verification_codes").FindOne(context.Background(), bson.M{"email": email}).Decode(&verificationCode)
	if err != nil {
		return false, errors.New("验证码不存在")
	}
	
	if verificationCode.Code != code {
		return false, errors.New("验证码错误")
	}
	
	if time.Since(verificationCode.CreatedAt) > 3*time.Minute {
		return false, errors.New("验证码已过期")
	}
	
	return true, nil
}

