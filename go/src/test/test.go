package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DeepSeekConfig 存储DeepSeek API配置
type DeepSeekConfig struct {
	APIKey  string
	BaseURL string
}

// Message 表示聊天消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 表示聊天请求
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// ChatResponse 表示聊天响应
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// NewDeepSeekConfig 创建新的DeepSeek配置
func NewDeepSeekConfig(apiKey string) *DeepSeekConfig {
	return &DeepSeekConfig{
		APIKey:  apiKey,
		BaseURL: "https://api.deepseek.com",
	}
}

// Chat 发送聊天请求到DeepSeek API
func (c *DeepSeekConfig) Chat(messages []Message) (*ChatResponse, error) {
	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)

	reqBody := ChatRequest{
		Model:    "deepseek-chat",
		Messages: messages,
		Stream:   false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	return &chatResp, nil
}
