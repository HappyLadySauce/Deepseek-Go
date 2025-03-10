package email

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"Deepseek-Go/global"
	"Deepseek-Go/models"
)

const (
	// 验证码过期时间（分钟）
	VerificationCodeExpiration = 3
)

var (
	ErrEmailInvalid     = errors.New("邮箱格式错误")
	ErrEmailNotAllowed  = errors.New("邮箱不符合要求，QQ邮箱需要5-11位数字，163邮箱需要6-18位字符，Gmail邮箱用户名最多30个字符")
	ErrEmailUsed        = errors.New("邮箱已被使用")
	ErrEmailNotVerified = errors.New("请先验证邮箱")
	ErrCodeExpired      = errors.New("验证码已过期，请重新发送")
	ErrCodeInvalid      = errors.New("验证码不正确")
	ErrCodeNotFound     = errors.New("未找到验证记录，请先发送验证码")
	ErrDBOperation      = errors.New("数据库操作失败")
	ErrSendVerification = errors.New("验证码发送失败")
)

// 验证邮箱格式和可用性
func ValidateEmail(email string) error {
	log.Printf("验证邮箱格式和可用性: %s", email)

	// 检查邮箱格式
	if !CheckEmailFormat(email) {
		log.Printf("邮箱格式错误: %s", email)
		return ErrEmailInvalid
	}

	// 检查邮箱是否允许
	if !AllowEmailFormat(email) {
		log.Printf("邮箱不符合要求: %s", email)
		return ErrEmailNotAllowed
	}

	// 检查邮箱是否已被其他用户使用
	var user models.User
	if global.DB.Where("email = ?", email).First(&user).RowsAffected > 0 {
		if user.EmailVerified {
			log.Printf("邮箱已被使用: %s", email)
			return ErrEmailUsed
		}
	}

	log.Printf("邮箱验证通过: %s", email)
	return nil
}

// 发送验证码
func SendVerificationCode(email string) error {
	log.Printf("开始发送验证码流程: %s", email)

	// 验证邮箱
	if err := ValidateEmail(email); err != nil {
		return err
	}

	// 生成验证码
	code := GenerateVerificationCode()
	log.Printf("为邮箱 %s 生成验证码: %s", email, code)

	// 更新或创建验证记录
	var verification models.EmailVerification
	result := global.DB.Where("email = ?", email).First(&verification)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("数据库查询错误: %v", result.Error)
		return ErrDBOperation
	}

	// 设置过期时间
	expiredAt := time.Now().Add(VerificationCodeExpiration * time.Minute)

	// 如果记录不存在，创建新记录
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		verification = models.EmailVerification{
			Email:     email,
			Code:      code,
			ExpiredAt: expiredAt,
		}
		if err := global.DB.Create(&verification).Error; err != nil {
			log.Printf("验证码创建失败: %v", err)
			return ErrDBOperation
		}
		log.Printf("为邮箱 %s 创建新的验证记录", email)
	} else {
		// 更新现有记录
		verification.Code = code
		verification.ExpiredAt = expiredAt
		verification.IsVerified = false
		if err := global.DB.Save(&verification).Error; err != nil {
			log.Printf("验证码更新失败: %v", err)
			return ErrDBOperation
		}
		log.Printf("为邮箱 %s 更新验证记录", email)
	}

	// 发送验证邮件
	log.Printf("开始发送验证码到邮箱 %s", email)
	if err := SendVerificationEmail(email, code); err != nil {
		log.Printf("验证码发送失败: %v", err)
		return fmt.Errorf("%w: %v", ErrSendVerification, err)
	}

	log.Printf("验证码发送成功: %s -> %s", email, code)
	return nil
}

// 验证邮箱验证码
func VerifyEmailCode(email, code string) error {
	log.Printf("开始验证邮箱验证码: %s, 验证码: %s", email, code)

	// 查询验证记录
	var verification models.EmailVerification
	if err := global.DB.Where("email = ?", email).First(&verification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("未找到验证记录: %s", email)
			return ErrCodeNotFound
		}
		log.Printf("数据库查询错误: %v", err)
		return ErrDBOperation
	}

	log.Printf("找到验证记录: 邮箱=%s, 验证码=%s, 过期时间=%v, 是否已验证=%v",
		verification.Email, verification.Code, verification.ExpiredAt, verification.IsVerified)

	// 检查验证码是否过期
	if time.Now().After(verification.ExpiredAt) {
		log.Printf("验证码已过期: %s, 过期时间: %v", email, verification.ExpiredAt)
		return ErrCodeExpired
	}

	// 检查验证码是否正确
	if verification.Code != code {
		log.Printf("验证码不正确: 期望=%s, 实际=%s", verification.Code, code)
		return ErrCodeInvalid
	}

	// 更新验证状态
	verification.IsVerified = true
	if err := global.DB.Save(&verification).Error; err != nil {
		log.Printf("验证状态更新失败: %v", err)
		return ErrDBOperation
	}

	// 如果有关联用户，更新用户的邮箱验证状态
	var user models.User
	if result := global.DB.Where("email = ?", email).First(&user); result.Error == nil {
		user.EmailVerified = true
		if err := global.DB.Save(&user).Error; err != nil {
			log.Printf("用户邮箱验证状态更新失败: %v", err)
			return ErrDBOperation
		}
		log.Printf("更新用户 %s 的邮箱验证状态", user.Username)
	}

	log.Printf("邮箱验证成功: %s", email)
	return nil
}

// 检查邮箱是否已验证
func CheckEmailVerified(email string) error {
	log.Printf("检查邮箱是否已验证: %s", email)

	var verification models.EmailVerification
	if err := global.DB.Where("email = ? AND is_verified = ?", email, true).First(&verification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("邮箱未验证: %s", email)
			return ErrEmailNotVerified
		}
		log.Printf("数据库查询错误: %v", err)
		return ErrDBOperation
	}

	log.Printf("邮箱已验证: %s", email)
	return nil
}
