// models/models.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 用户集合
type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Username      string             `bson:"username" index:"unique"`
	PasswordHash  string             `bson:"password_hash"`
	Email         string             `bson:"email" index:"unique"`
	EmailVerified bool               `bson:"email_verified"`
	Sex           string             `bson:"sex"`
	Avatar        string             `bson:"avatar"`
	Role          string             `bson:"role"`
	Introduction  string             `bson:"introduction"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

// 验证码集合
type VerificationCode struct {
	Email string `bson:"email"`
	Code  string `bson:"code"`
	CreatedAt time.Time `bson:"created_at"`
}
