package global

import (
	"gorm.io/gorm"
	"github.com/go-redis/redis/v8"
)

var (
	DB *gorm.DB
	RedisDB *redis.Client
)