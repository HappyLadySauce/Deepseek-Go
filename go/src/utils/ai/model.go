package ai

import (
	"Deepseek-Go/config"
	"context"
	"errors"
)

// ChatMessage 定义聊天消息的结构
type ChatMessage struct {
	Role    string `json:"role"`    // 消息角色：user, assistant, system
	Content string `json:"content"` // 消息内容
}

// ChatCompletionRequest 定义聊天请求参数
type ChatCompletionRequest struct {
	Model        string        `json:"model"`                   // 模型名称
	Messages     []ChatMessage `json:"messages"`                // 消息历史
	Temperature  float64       `json:"temperature,omitempty"`   // 温度参数，控制随机性
	MaxTokens    int           `json:"max_tokens,omitempty"`    // 最大token数
	Stream       bool          `json:"stream,omitempty"`        // 是否使用流式输出
	KnowledgeIDs []uint        `json:"knowledge_ids,omitempty"` // 知识库ID列表
}

// ChatCompletionResponse 定义聊天响应结构
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// AIModel 定义AI模型接口
type AIModel interface {
	ChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error)
	StreamChatCompletion(ctx context.Context, request ChatCompletionRequest, callback func(response *ChatCompletionResponse)) error
}

// GetAIModel 根据提供商获取对应的AI模型实例
func GetAIModel(provider string) (AIModel, error) {
	switch provider {
	case "deepseek":
		return NewDeepSeekModel(config.Config.AI.DeepSeek.APIKey, config.Config.AI.DeepSeek.BaseURL), nil
	case "kimi":
		return NewKimiModel(config.Config.AI.Kimi.APIKey, config.Config.AI.Kimi.BaseURL), nil
	default:
		return nil, errors.New("不支持的AI提供商")
	}
}
