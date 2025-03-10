package models

import (
	"time"

	"gorm.io/gorm"
)

// EmailVerification 邮箱验证模型
type EmailVerification struct {
	gorm.Model
	Email      string    `json:"email" binding:"required" gorm:"index"`
	Code       string    `json:"code" binding:"required"`
	ExpiredAt  time.Time `json:"expired_at"`
	IsVerified bool      `json:"is_verified" gorm:"default:false"`
}
