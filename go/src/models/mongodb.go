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
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}