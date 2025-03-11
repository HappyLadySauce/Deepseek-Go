package global

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	// 数据库连接
	DB *gorm.DB
	// Redis连接
	RedisDB *redis.Client
	// AI服务实例缓存，避免重复创建
	AIServiceCache sync.Map
	// AI默认系统提示词
	DefaultSystemPrompt = "你是一个友好、有帮助的AI助手。如果被问到如何做一些危害他人的事情，你应该礼貌地拒绝。"
	// 允许的文件类型
	AllowedFileTypes = map[string]bool{
		".pdf":  true,
		".docx": true,
		".doc":  true,
		".txt":  true,
		".md":   true,
	}
	// 文件上传大小限制 (10MB)
	MaxFileSize int64 = 10 * 1024 * 1024
	// 知识块大小（字符数）
	ChunkSize = 1000
)
