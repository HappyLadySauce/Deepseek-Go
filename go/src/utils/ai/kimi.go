package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// KimiModel 表示Moonshot AI的Kimi模型实现
type KimiModel struct {
	apiKey  string
	baseURL string
}

// NewKimiModel 创建一个新的Kimi模型实例
func NewKimiModel(apiKey, baseURL string) *KimiModel {
	return &KimiModel{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// ChatCompletion 实现非流式聊天接口
func (m *KimiModel) ChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error) {
	url := fmt.Sprintf("%s/v1/chat/completions", m.baseURL)

	// 转换请求为JSON格式
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("JSON编码请求失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.apiKey))

	// 发送请求
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析响应
	var response ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &response, nil
}

// StreamChatCompletion 实现流式聊天接口
func (m *KimiModel) StreamChatCompletion(ctx context.Context, request ChatCompletionRequest, callback func(response *ChatCompletionResponse)) error {
	// 确保请求是流式的
	request.Stream = true
	url := fmt.Sprintf("%s/v1/chat/completions", m.baseURL)

	// 转换请求为JSON格式
	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("JSON编码请求失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.apiKey))
	req.Header.Set("Accept", "text/event-stream")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	// 读取SSE流
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("读取流失败: %v", err)
		}

		// 跳过空行
		if line == "\n" || line == "\r\n" {
			continue
		}

		// 解析SSE
		if len(line) >= 6 && line[0:6] == "data: " {
			data := line[6:]

			// 检查流结束
			if data == "[DONE]\n" || data == "[DONE]\r\n" {
				break
			}

			// 解析JSON数据
			var response ChatCompletionResponse
			if err := json.Unmarshal([]byte(data), &response); err != nil {
				return fmt.Errorf("解析SSE数据失败: %v", err)
			}

			// 调用回调函数处理数据
			callback(&response)
		}
	}

	return nil
} 