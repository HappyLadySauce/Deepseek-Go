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
	MySQL struct {
		Host               string `mapstructure:"host"`
		Port               int    `mapstructure:"port"`
		User               string `mapstructure:"user"`
		Password           string `mapstructure:"password"`
		Name               string `mapstructure:"name"`
		Charset            string `mapstructure:"charset"`
		ParseTime          bool   `mapstructure:"parseTime"`
		Loc                string `mapstructure:"loc"`
		SetMaxIdleConns    int    `mapstructure:"SetMaxIdleConns"`
		SetMaxOpenConns    int    `mapstructure:"SetMaxOpenConns"`
		SetConnMaxLifetime int    `mapstructure:"SetConnMaxLifetime"`
		SetConnMaxIdleTime int    `mapstructure:"SetConnMaxIdleTime"`
	}
	MongoDB struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}
	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
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
		ServerName string `mapstructure:"server_name"`
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
