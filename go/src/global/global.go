package global

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	// 数据库连接
	DB *gorm.DB
	// Redis连接
	RedisDB *redis.Client
)
