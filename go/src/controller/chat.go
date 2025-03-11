package controller

import (
	"Deepseek-Go/models"
	"Deepseek-Go/utils/ai"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 聊天控制器结构体
type ChatController struct {
	DB        *gorm.DB
	AIService *ai.AIService
}

// 聊天请求结构体
type ChatRequest struct {
	SessionID    uint   `json:"session_id"`
	Message      string `json:"message"`
	AIConfigID   uint   `json:"ai_config_id"` // 0表示使用默认配置
	KnowledgeIDs []uint `json:"knowledge_ids"`
}

// AI配置请求结构体
type AIConfigRequest struct {
	ModelName   string  `json:"model_name"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	Provider    string  `json:"provider"`
}

// 聊天响应结构体
type ChatResponse struct {
	ID        uint   `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// NewChatController 创建聊天控制器
func NewChatController(db *gorm.DB) *ChatController {
	return &ChatController{
		DB:        db,
		AIService: ai.NewAIService(db),
	}
}

// Chat 处理聊天请求
func (cc *ChatController) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 获取AI配置
	var aiConfig models.AIConfig
	var err error

	if req.AIConfigID > 0 {
		// 使用指定的配置
		aiConfig, err = cc.getAIConfig(req.AIConfigID, userID.(uint))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 使用默认配置
		config, err := cc.AIService.GetDefaultAIConfig(userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取默认AI配置失败: " + err.Error()})
			return
		}
		aiConfig = *config
	}

	// 调用AI服务处理聊天
	assistantMessage, session, err := cc.AIService.Chat(userID.(uint), req.SessionID, req.Message, aiConfig, req.KnowledgeIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "聊天处理失败: " + err.Error()})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "聊天成功",
		"data": ChatResponse{
			ID:        assistantMessage.ID,
			Role:      assistantMessage.Role,
			Content:   assistantMessage.Content,
			CreatedAt: assistantMessage.CreatedAt.Format(time.RFC3339),
		},
		"session_id": session.ID,
	})
}

// StreamChat 处理流式聊天请求
func (cc *ChatController) StreamChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 获取AI配置
	var aiConfig models.AIConfig
	var err error

	if req.AIConfigID > 0 {
		// 使用指定的配置
		aiConfig, err = cc.getAIConfig(req.AIConfigID, userID.(uint))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 使用默认配置
		config, err := cc.AIService.GetDefaultAIConfig(userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取默认AI配置失败: " + err.Error()})
			return
		}
		aiConfig = *config
	}

	// 设置SSE响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// 确保数据立即发送
	c.Writer.Flush()

	// 处理流式回复的回调函数
	callback := func(response *ai.ChatCompletionResponse) {
		if len(response.Choices) > 0 {
			chunk := response.Choices[0].Message.Content

			// 发送数据到客户端
			data, _ := json.Marshal(gin.H{
				"id":      response.ID,
				"content": chunk,
				"done":    false,
			})
			c.Writer.Write([]byte("data: " + string(data) + "\n\n"))
			c.Writer.Flush()
		}
	}

	// 调用AI服务处理流式聊天
	_, session, err := cc.AIService.StreamChat(userID.(uint), req.SessionID, req.Message, aiConfig, req.KnowledgeIDs, c.Writer, callback)
	if err != nil {
		// 发送错误信息
		data, _ := json.Marshal(gin.H{
			"error": "AI服务调用失败: " + err.Error(),
			"done":  true,
		})
		c.Writer.Write([]byte("data: " + string(data) + "\n\n"))
		c.Writer.Flush()
		return
	}

	// 发送完成消息
	data, _ := json.Marshal(gin.H{
		"id":         "done",
		"content":    "",
		"done":       true,
		"session_id": session.ID,
	})
	c.Writer.Write([]byte("data: " + string(data) + "\n\n"))
	c.Writer.Flush()
}

// GetSessions 获取用户的所有聊天会话
func (cc *ChatController) GetSessions(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 调用服务获取会话列表
	sessions, count, err := cc.AIService.GetSessions(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取会话列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取会话列表成功",
		"data": gin.H{
			"total":    count,
			"page":     page,
			"pageSize": pageSize,
			"sessions": sessions,
		},
	})
}

// GetSessionMessages 获取特定会话的消息历史
func (cc *ChatController) GetSessionMessages(c *gin.Context) {
	// 获取会话ID
	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的会话ID"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 调用服务获取消息历史
	messages, count, err := cc.AIService.GetSessionMessages(uint(sessionID), userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 格式化返回数据
	var formattedMessages []ChatResponse
	for _, msg := range messages {
		formattedMessages = append(formattedMessages, ChatResponse{
			ID:        msg.ID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取消息历史成功",
		"data": gin.H{
			"total":    count,
			"page":     page,
			"pageSize": pageSize,
			"messages": formattedMessages,
		},
	})
}

// UpdateSession 更新会话信息（如标题）
func (cc *ChatController) UpdateSession(c *gin.Context) {
	// 获取会话ID
	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的会话ID"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 解析请求体
	var updateData struct {
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 调用服务更新会话
	session, err := cc.AIService.UpdateSession(uint(sessionID), userID.(uint), updateData.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新会话成功",
		"data":    session,
	})
}

// DeleteSession 删除会话
func (cc *ChatController) DeleteSession(c *gin.Context) {
	// 获取会话ID
	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的会话ID"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 调用服务删除会话
	if err := cc.AIService.DeleteSession(uint(sessionID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除会话成功",
	})
}

// 辅助方法

// getAIConfig 获取AI配置并验证所有权
func (cc *ChatController) getAIConfig(configID, userID uint) (models.AIConfig, error) {
	config, err := cc.AIService.GetAIConfig(configID, userID)
	if err != nil {
		return models.AIConfig{}, err
	}
	return *config, nil
}
