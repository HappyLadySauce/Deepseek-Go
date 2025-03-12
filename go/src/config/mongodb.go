package config

import (
	"log"
	"strconv"

	"Deepseek-Go/global"

	"go.mongodb.org/mongo-driver/mongo"
)
func InitMongoDB() {
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(Config.MongoDB.Host + ":" + strconv.Itoa(Config.MongoDB.Port)))
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}

	err = db.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}

	global.MongoDBClient = db
}
