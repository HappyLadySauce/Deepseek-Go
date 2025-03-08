package config

import (
	"log"
	"github.com/spf13/viper"
)

type config struct {
	App struct {
		Name string `mapstructure:"name"`
		Port int    `mapstructure:"port"`
	}
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		Charset  string `mapstructure:"charset"`
		ParseTime bool   `mapstructure:"parseTime"`
		Loc       string `mapstructure:"loc"`
		SetMaxIdleConns    int    `mapstructure:"SetMaxIdleConns"`
		SetMaxOpenConns    int    `mapstructure:"SetMaxOpenConns"`
		SetConnMaxLifetime int    `mapstructure:"SetConnMaxLifetime"`
		SetConnMaxIdleTime  int    `mapstructure:"SetConnMaxIdleTime"`
	}
	DeepSeek struct {
		APIKey  string `mapstructure:"api_key"`
		BaseURL string `mapstructure:"base_url"`
		Model   string `mapstructure:"model"`
		Stream  bool   `mapstructure:"stream"`
	}
}

var Config *config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	Config = &config{}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	if Config.App.Port == 0 {
		Config.App.Port = 14020
	}

	// 初始化数据库
	InitDB()
}
