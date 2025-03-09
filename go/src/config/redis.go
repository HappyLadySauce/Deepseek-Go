package config

import (
	"Deepseek-Go/global"
	"context"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func InitRedis() {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Host + ":" + strconv.Itoa(Config.Redis.Port),
		Password: Config.Redis.Password,
		DB:       Config.Redis.DB,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}

	global.RedisDB = RedisClient
}
