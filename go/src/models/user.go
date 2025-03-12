package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// MySQL 用户表
// type UserMysql struct {
// 	gorm.Model
// 	Username      string `json:"username" binding:"required" gorm:"unique"`
// 	Password      string `json:"password" binding:"required"`
// 	Email         string `json:"email" binding:"required" gorm:"unique"`
// 	EmailVerified bool   `json:"email_verified" gorm:"default:false"`
// }

// MongoDB 用户表
type UserMongo struct {
	ID				primitive.ObjectID 		`bson:"_id,omitempty"`
	Username      	string 				`bson: "username"`
	PasswordHash    string 				`bson: "password"`
	Email         	string 				`bson: "email"`
	EmailVerified 	bool   				`bson: "email_verified"`
}

