package config

import (
	"context"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Client

func InitMongoDB() {
	// 连接MongoDB
	// mongo.Connect() 函数用于连接 MongoDB 数据库。它接受一个上下文（context）和一个选项（options）。
	// 上下文用于管理连接的生命周期，而选项则用于配置连接参数。
	// 连接字符串 格式：mongodb://用户名:密码@host:port/数据库名
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://" + Config.MongoDB.Host + ":" + strconv.Itoa(Config.MongoDB.Port)))
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}

	// 检查连接
	// Ping() 方法用于检查 MongoDB 数据库的连接状态。它接受一个上下文（context）和一个选项（options）。
	err = db.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	
	// 设置全局变量
	// 将连接对象赋值给全局变量
	MongoDB = db
}