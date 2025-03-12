package global

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	// 数据库连接
	DB *gorm.DB
	// Redis连接
	RedisDB *redis.Client
	// MongoDB连接
	MongoDB *mongo.Client
)
