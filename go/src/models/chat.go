package models

import (
	"time"

	"gorm.io/gorm"
)

// ChatSession 聊天会话模型
type ChatSession struct {
	gorm.Model
	UserID      uint   `json:"user_id" gorm:"index"` // 用户ID
	Title       string `json:"title"`                // 会话标题
	LastMessage string `json:"last_message"`         // 最后一条消息内容
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	gorm.Model
	SessionID uint      `json:"session_id" gorm:"index"`  // 所属会话ID
	Role      string    `json:"role"`                     // 消息角色：user 或 assistant
	Content   string    `json:"content" gorm:"type:text"` // 消息内容
	CreatedAt time.Time `json:"created_at"`               // 创建时间
}

// KnowledgeFile 知识库文件模型
type KnowledgeFile struct {
	gorm.Model
	UserID      uint       `json:"user_id" gorm:"index"` // 上传用户ID
	FileName    string     `json:"file_name"`            // 文件名称
	FilePath    string     `json:"file_path"`            // 文件路径
	FileSize    int64      `json:"file_size"`            // 文件大小(bytes)
	FileType    string     `json:"file_type"`            // 文件类型(如pdf, docx, txt)
	ProcessedAt *time.Time `json:"processed_at"`         // 处理完成时间
	Status      string     `json:"status"`               // 处理状态: pending, processing, completed, failed
}

// KnowledgeVectorStore 知识库向量存储模型
type KnowledgeVectorStore struct {
	gorm.Model
	FileID    uint   `json:"file_id" gorm:"index"`       // 关联的文件ID
	Text      string `json:"text" gorm:"type:text"`      // 文本片段
	Embedding []byte `json:"embedding" gorm:"type:blob"` // 向量嵌入
	Metadata  string `json:"metadata" gorm:"type:json"`  // 元数据JSON
}

// AIConfig AI配置模型
type AIConfig struct {
	gorm.Model
	UserID      uint    `json:"user_id" gorm:"index"` // 用户ID
	ModelName   string  `json:"model_name"`           // 模型名称
	Temperature float64 `json:"temperature"`          // 温度参数
	MaxTokens   int     `json:"max_tokens"`           // 最大Token数
	Provider    string  `json:"provider"`             // 提供商 (deepseek, kimi)
	IsDefault   bool    `json:"is_default"`           // 是否为默认配置
}
