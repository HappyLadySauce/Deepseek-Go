package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required" gorm:"unique"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required" gorm:"unique"`
}