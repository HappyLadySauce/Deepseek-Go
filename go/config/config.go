package config

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	App struct {
		Port int    `mapstructure:"port"`
	}
	MongoDB struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}
	Cors struct {
		AllowOrigins     []string `mapstructure:"allow_origins"`
		AllowCredentials bool     `mapstructure:"allow_credentials"`
		AllowMethods     []string `mapstructure:"allow_methods"`
		AllowHeaders     []string `mapstructure:"allow_headers"`
		MaxAge           int      `mapstructure:"max_age"`
	}
	Email struct {
		Host       string `mapstructure:"host"`
		Port       int    `mapstructure:"port"`
		Username   string `mapstructure:"username"`
		Password   string `mapstructure:"password"`
		From       string `mapstructure:"from"`
		EnableSSL  bool   `mapstructure:"enable_ssl"`
	}
	AI struct {
		// 深度求索 api
		DeepSeek struct {
			APIKey string `mapstructure:"api_key"`
			BaseURL string `mapstructure:"base_url"`
		}
		// 月之暗面 api
		Kimi struct {
			APIKey string `mapstructure:"api_key"`
			BaseURL string `mapstructure:"base_url"`
		}
	}
}

// 配置文件
var Config *config

// 初始化配置文件
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
	// InitMysql()
	// 初始化Redis
	// InitRedis()
	// 初始化MongoDB
	InitMongoDB()
}
