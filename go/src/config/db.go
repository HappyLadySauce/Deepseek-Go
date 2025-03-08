package config

import (
	"Deepseek-Go/global"
	"log"
	"net/url"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DSN() string {
	return Config.Database.User + ":" + Config.Database.Password + "@tcp(" + Config.Database.Host + ":" + strconv.Itoa(Config.Database.Port) + ")/" + Config.Database.Name + "?charset=" + Config.Database.Charset + "&parseTime=" + strconv.FormatBool(Config.Database.ParseTime) + "&loc=" + url.QueryEscape(Config.Database.Loc)
}

// 连接数据库
func InitDB() {
	db, err := gorm.Open(mysql.Open(DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库连接失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(Config.Database.SetMaxIdleConns)
	sqlDB.SetMaxOpenConns(Config.Database.SetMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(Config.Database.SetConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(Config.Database.SetConnMaxIdleTime) * time.Second)

	global.DB = db
}
