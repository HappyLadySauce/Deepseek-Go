package config

import (
	"Deepseek-Go/global"
	"Deepseek-Go/models"
	"log"
	"net/url"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DSN() string {
	return Config.MySQL.User + ":" + Config.MySQL.Password + "@tcp(" + Config.MySQL.Host + ":" + strconv.Itoa(Config.MySQL.Port) + ")/" + Config.MySQL.Name + "?charset=" + Config.MySQL.Charset + "&parseTime=" + strconv.FormatBool(Config.MySQL.ParseTime) + "&loc=" + url.QueryEscape(Config.MySQL.Loc)
}

// 连接数据库
func InitMysql() {
	db, err := gorm.Open(mysql.Open(DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库连接失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(Config.MySQL.SetMaxIdleConns)
	sqlDB.SetMaxOpenConns(Config.MySQL.SetMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(Config.MySQL.SetConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(Config.MySQL.SetConnMaxIdleTime) * time.Second)

	global.DB = db

	// 数据库迁移
	migrateMysql()
}

// 数据库迁移
func migrateMysql() {
	err := global.DB.AutoMigrate(
		&models.User{},
		&models.EmailVerification{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移成功")
}
