package ai

import (
	"Deepseek-Go/global"
	"Deepseek-Go/models"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AIService 提供AI服务的结构体
type AIService struct {
	DB *gorm.DB
}

// NewAIService 创建新的AI服务实例
func NewAIService(db *gorm.DB) *AIService {
	return &AIService{
		DB: db,
	}
}

// 聊天相关服务 ---------------------------------------------------------

// Chat 处理普通聊天请求
func (s *AIService) Chat(userID uint, sessionID uint, message string, aiConfig models.AIConfig, knowledgeIDs []uint) (*models.ChatMessage, *models.ChatSession, error) {
	// 获取或创建会话
	session, err := s.getOrCreateSession(userID, sessionID, message)
	if err != nil {
		return nil, nil, fmt.Errorf("会话处理失败: %v", err)
	}

	// 获取历史消息
	messages, err := s.getSessionMessages(session.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("获取历史消息失败: %v", err)
	}

	// 构建AI请求消息
	aiMessages := s.buildAIMessages(messages, message, knowledgeIDs, userID)

	// 保存用户消息
	userMessage := models.ChatMessage{
		SessionID: session.ID,
		Role:      "user",
		Content:   message,
		CreatedAt: time.Now(),
	}
	if err := s.DB.Create(&userMessage).Error; err != nil {
		return nil, nil, fmt.Errorf("保存用户消息失败: %v", err)
	}

	// 更新会话最后消息
	session.LastMessage = message
	s.DB.Save(&session)

	// 获取AI模型
	aiModel, err := GetAIModel(aiConfig.Provider)
	if err != nil {
		return nil, nil, fmt.Errorf("获取AI模型失败: %v", err)
	}

	// 调用AI服务
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := aiModel.ChatCompletion(ctx, ChatCompletionRequest{
		Model:       aiConfig.ModelName,
		Messages:    aiMessages,
		Temperature: aiConfig.Temperature,
		MaxTokens:   aiConfig.MaxTokens,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("AI服务调用失败: %v", err)
	}

	// 提取AI回复
	if len(response.Choices) == 0 {
		return nil, nil, fmt.Errorf("AI返回了空回复")
	}

	aiReply := response.Choices[0].Message.Content

	// 保存AI回复到数据库
	assistantMessage := models.ChatMessage{
		SessionID: session.ID,
		Role:      "assistant",
		Content:   aiReply,
		CreatedAt: time.Now(),
	}
	if err := s.DB.Create(&assistantMessage).Error; err != nil {
		return nil, nil, fmt.Errorf("保存AI回复失败: %v", err)
	}

	// 更新会话最后消息
	session.LastMessage = aiReply
	s.DB.Save(&session)

	return &assistantMessage, session, nil
}

// StreamChat 处理流式聊天
func (s *AIService) StreamChat(userID uint, sessionID uint, message string, aiConfig models.AIConfig, knowledgeIDs []uint, writer gin.ResponseWriter, callback func(response *ChatCompletionResponse)) (string, *models.ChatSession, error) {
	// 获取或创建会话
	session, err := s.getOrCreateSession(userID, sessionID, message)
	if err != nil {
		return "", nil, fmt.Errorf("会话处理失败: %v", err)
	}

	// 获取历史消息
	messages, err := s.getSessionMessages(session.ID)
	if err != nil {
		return "", nil, fmt.Errorf("获取历史消息失败: %v", err)
	}

	// 构建AI请求消息
	aiMessages := s.buildAIMessages(messages, message, knowledgeIDs, userID)

	// 保存用户消息
	userMessage := models.ChatMessage{
		SessionID: session.ID,
		Role:      "user",
		Content:   message,
		CreatedAt: time.Now(),
	}
	if err := s.DB.Create(&userMessage).Error; err != nil {
		return "", nil, fmt.Errorf("保存用户消息失败: %v", err)
	}

	// 更新会话最后消息
	session.LastMessage = message
	s.DB.Save(&session)

	// 获取AI模型
	aiModel, err := GetAIModel(aiConfig.Provider)
	if err != nil {
		return "", nil, fmt.Errorf("获取AI模型失败: %v", err)
	}

	// 用于收集完整回复的缓冲区
	var fullReply string

	// 调用AI服务（流式）
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// 处理流式回复的回调函数
	streamCallback := func(response *ChatCompletionResponse) {
		if len(response.Choices) > 0 {
			chunk := response.Choices[0].Message.Content
			fullReply += chunk
		}
		callback(response)
	}

	err = aiModel.StreamChatCompletion(ctx, ChatCompletionRequest{
		Model:       aiConfig.ModelName,
		Messages:    aiMessages,
		Temperature: aiConfig.Temperature,
		MaxTokens:   aiConfig.MaxTokens,
		Stream:      true,
	}, streamCallback)

	if err != nil {
		return "", nil, fmt.Errorf("AI服务调用失败: %v", err)
	}

	// 保存完整回复到数据库
	assistantMessage := models.ChatMessage{
		SessionID: session.ID,
		Role:      "assistant",
		Content:   fullReply,
		CreatedAt: time.Now(),
	}
	if err := s.DB.Create(&assistantMessage).Error; err != nil {
		return "", nil, fmt.Errorf("保存AI回复失败: %v", err)
	}

	// 更新会话最后消息
	session.LastMessage = fullReply
	s.DB.Save(&session)

	return fullReply, session, nil
}

// 会话管理服务 ---------------------------------------------------------

// GetSessions 获取用户的所有会话
func (s *AIService) GetSessions(userID uint, page, pageSize int) ([]models.ChatSession, int64, error) {
	var sessions []models.ChatSession
	var count int64

	// 获取总数
	if err := s.DB.Model(&models.ChatSession{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := s.DB.Where("user_id = ?", userID).Order("updated_at desc").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, count, nil
}

// GetSessionMessages 获取会话的所有消息
func (s *AIService) GetSessionMessages(sessionID, userID uint, page, pageSize int) ([]models.ChatMessage, int64, error) {
	// 验证会话存在性和所有权
	var session models.ChatSession
	if err := s.DB.First(&session, sessionID).Error; err != nil {
		return nil, 0, fmt.Errorf("会话不存在")
	}

	if session.UserID != userID {
		return nil, 0, fmt.Errorf("无权访问此会话")
	}

	var messages []models.ChatMessage
	var count int64

	// 获取总数
	if err := s.DB.Model(&models.ChatMessage{}).Where("session_id = ?", sessionID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := s.DB.Where("session_id = ?", sessionID).Order("created_at asc").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, count, nil
}

// UpdateSession 更新会话信息
func (s *AIService) UpdateSession(sessionID, userID uint, title string) (*models.ChatSession, error) {
	// 验证会话存在性和所有权
	var session models.ChatSession
	if err := s.DB.First(&session, sessionID).Error; err != nil {
		return nil, fmt.Errorf("会话不存在")
	}

	if session.UserID != userID {
		return nil, fmt.Errorf("无权修改此会话")
	}

	// 更新会话标题
	session.Title = title
	if err := s.DB.Save(&session).Error; err != nil {
		return nil, fmt.Errorf("更新会话失败: %v", err)
	}

	return &session, nil
}

// DeleteSession 删除会话
func (s *AIService) DeleteSession(sessionID, userID uint) error {
	// 验证会话存在性和所有权
	var session models.ChatSession
	if err := s.DB.First(&session, sessionID).Error; err != nil {
		return fmt.Errorf("会话不存在")
	}

	if session.UserID != userID {
		return fmt.Errorf("无权删除此会话")
	}

	// 开启事务
	tx := s.DB.Begin()

	// 删除会话中的所有消息
	if err := tx.Where("session_id = ?", sessionID).Delete(&models.ChatMessage{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除会话消息失败: %v", err)
	}

	// 删除会话本身
	if err := tx.Delete(&session).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除会话失败: %v", err)
	}

	// 提交事务
	tx.Commit()
	return nil
}

// 知识库相关服务 ---------------------------------------------------------

// UploadKnowledgeFile 上传知识库文件
func (s *AIService) UploadKnowledgeFile(userID uint, fileName string, fileSize int64, fileType string, filePath string) (*models.KnowledgeFile, error) {
	// 保存文件信息到数据库
	knowledgeFile := models.KnowledgeFile{
		UserID:   userID,
		FileName: fileName,
		FilePath: filePath,
		FileSize: fileSize,
		FileType: fileType,
		Status:   "pending", // 初始状态为待处理
	}

	if err := s.DB.Create(&knowledgeFile).Error; err != nil {
		return nil, fmt.Errorf("保存文件记录失败: %v", err)
	}

	// 启动异步处理
	go s.ProcessKnowledgeFile(knowledgeFile)

	return &knowledgeFile, nil
}

// ProcessKnowledgeFile 处理知识库文件
func (s *AIService) ProcessKnowledgeFile(file models.KnowledgeFile) {
	// 更新状态为处理中
	s.DB.Model(&file).Update("status", "processing")

	// 模拟处理时间
	time.Sleep(2 * time.Second)

	// 读取文件内容
	content, err := os.ReadFile(file.FilePath)
	if err != nil {
		s.DB.Model(&file).Update("status", "failed")
		return
	}

	// 文本分块
	text := string(content)
	chunks := s.ChunkText(text, global.ChunkSize)

	// 保存文本块到向量存储
	for _, chunk := range chunks {
		vectorStore := models.KnowledgeVectorStore{
			FileID:    file.ID,
			Text:      chunk,
			Embedding: []byte{}, // 实际应该调用嵌入模型
			Metadata:  `{"source": "` + file.FileName + `"}`,
		}

		if err := s.DB.Create(&vectorStore).Error; err != nil {
			s.DB.Model(&file).Update("status", "failed")
			return
		}
	}

	// 更新状态为完成
	now := time.Now()
	s.DB.Model(&file).Updates(map[string]interface{}{
		"status":       "completed",
		"processed_at": &now,
	})
}

// GetKnowledgeFiles 获取知识库文件列表
func (s *AIService) GetKnowledgeFiles(userID uint, page, pageSize int) ([]models.KnowledgeFile, int64, error) {
	var files []models.KnowledgeFile
	var count int64

	// 获取总数
	if err := s.DB.Model(&models.KnowledgeFile{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := s.DB.Where("user_id = ?", userID).Order("created_at desc").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, count, nil
}

// GetKnowledgeFile 获取单个知识库文件详情
func (s *AIService) GetKnowledgeFile(fileID, userID uint) (*models.KnowledgeFile, int64, error) {
	// 获取文件信息
	var file models.KnowledgeFile
	if err := s.DB.First(&file, fileID).Error; err != nil {
		return nil, 0, fmt.Errorf("文件不存在")
	}

	// 验证文件所有权
	if file.UserID != userID {
		return nil, 0, fmt.Errorf("无权访问此文件")
	}

	// 获取文件的向量存储数量
	var vectorCount int64
	s.DB.Model(&models.KnowledgeVectorStore{}).Where("file_id = ?", fileID).Count(&vectorCount)

	return &file, vectorCount, nil
}

// DeleteKnowledgeFile 删除知识库文件
func (s *AIService) DeleteKnowledgeFile(fileID, userID uint) error {
	// 获取文件信息
	var file models.KnowledgeFile
	if err := s.DB.First(&file, fileID).Error; err != nil {
		return fmt.Errorf("文件不存在")
	}

	// 验证文件所有权
	if file.UserID != userID {
		return fmt.Errorf("无权删除此文件")
	}

	// 开启事务
	tx := s.DB.Begin()

	// 删除向量存储
	if err := tx.Where("file_id = ?", fileID).Delete(&models.KnowledgeVectorStore{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除向量存储失败: %v", err)
	}

	// 删除文件记录
	if err := tx.Delete(&file).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除文件记录失败: %v", err)
	}

	// 提交事务
	tx.Commit()

	// 删除物理文件
	os.Remove(file.FilePath)

	return nil
}

// SaveKnowledgeFile 保存上传的知识库文件到磁盘
func (s *AIService) SaveKnowledgeFile(file io.Reader, originalFileName string) (string, string, error) {
	// 创建上传目录
	uploadDir := "./uploads/knowledge"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", "", fmt.Errorf("创建上传目录失败: %v", err)
	}

	// 生成唯一文件名
	fileExt := filepath.Ext(originalFileName)
	newFileName := uuid.New().String() + fileExt
	filePath := filepath.Join(uploadDir, newFileName)

	// 保存文件
	out, err := os.Create(filePath)
	if err != nil {
		return "", "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	// 复制文件内容
	if _, err = io.Copy(out, file); err != nil {
		return "", "", fmt.Errorf("保存文件失败: %v", err)
	}

	return filePath, strings.TrimPrefix(fileExt, "."), nil
}

// AI配置相关服务 ---------------------------------------------------------

// GetDefaultAIConfig 获取用户默认AI配置
func (s *AIService) GetDefaultAIConfig(userID uint) (*models.AIConfig, error) {
	var config models.AIConfig
	err := s.DB.Where("user_id = ? AND is_default = ?", userID, true).First(&config).Error

	if err != nil {
		// 如果没有默认配置，则创建一个
		newConfig := models.AIConfig{
			UserID:      userID,
			ModelName:   "deepseek-v1-8k",
			Temperature: 0.7,
			MaxTokens:   2048,
			Provider:    "deepseek",
			IsDefault:   true,
		}

		if err := s.DB.Create(&newConfig).Error; err != nil {
			return nil, fmt.Errorf("创建默认配置失败: %v", err)
		}

		return &newConfig, nil
	}

	return &config, nil
}

// GetAIConfig 获取单个AI配置
func (s *AIService) GetAIConfig(configID, userID uint) (*models.AIConfig, error) {
	var config models.AIConfig
	if err := s.DB.First(&config, configID).Error; err != nil {
		return nil, fmt.Errorf("配置不存在")
	}

	if config.UserID != userID {
		return nil, fmt.Errorf("无权访问此配置")
	}

	return &config, nil
}

// GetAIConfigs 获取所有AI配置
func (s *AIService) GetAIConfigs(userID uint) ([]models.AIConfig, error) {
	var configs []models.AIConfig
	if err := s.DB.Where("user_id = ?", userID).Find(&configs).Error; err != nil {
		return nil, err
	}

	// 如果用户没有任何配置，则创建默认配置
	if len(configs) == 0 {
		// 为DeepSeek创建默认配置
		deepseekConfig := models.AIConfig{
			UserID:      userID,
			ModelName:   "deepseek-v1-8k",
			Temperature: 0.7,
			MaxTokens:   2048,
			Provider:    "deepseek",
			IsDefault:   true,
		}

		// 为Kimi创建默认配置
		kimiConfig := models.AIConfig{
			UserID:      userID,
			ModelName:   "moonshot-v1-8k",
			Temperature: 0.7,
			MaxTokens:   2048,
			Provider:    "kimi",
			IsDefault:   false,
		}

		// 保存默认配置
		s.DB.Create(&deepseekConfig)
		s.DB.Create(&kimiConfig)

		// 添加到返回结果
		configs = append(configs, deepseekConfig, kimiConfig)
	}

	return configs, nil
}

// CreateAIConfig 创建AI配置
func (s *AIService) CreateAIConfig(userID uint, modelName string, temperature float64, maxTokens int, provider string, isDefault bool) (*models.AIConfig, error) {
	// 如果设置为默认，则将其他配置设为非默认
	if isDefault {
		if err := s.DB.Model(&models.AIConfig{}).Where("user_id = ?", userID).
			Update("is_default", false).Error; err != nil {
			return nil, fmt.Errorf("更新默认配置状态失败: %v", err)
		}
	}

	// 创建新配置
	config := models.AIConfig{
		UserID:      userID,
		ModelName:   modelName,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		Provider:    provider,
		IsDefault:   isDefault,
	}

	if err := s.DB.Create(&config).Error; err != nil {
		return nil, fmt.Errorf("创建配置失败: %v", err)
	}

	return &config, nil
}

// UpdateAIConfig 更新AI配置
func (s *AIService) UpdateAIConfig(configID, userID uint, modelName string, temperature float64, maxTokens int, provider string, isDefault bool) (*models.AIConfig, error) {
	// 获取配置
	var config models.AIConfig
	if err := s.DB.First(&config, configID).Error; err != nil {
		return nil, fmt.Errorf("配置不存在")
	}

	// 验证配置所有权
	if config.UserID != userID {
		return nil, fmt.Errorf("无权修改此配置")
	}

	// 如果将配置设置为默认，则将其他配置设为非默认
	if isDefault && !config.IsDefault {
		if err := s.DB.Model(&models.AIConfig{}).Where("user_id = ?", userID).
			Update("is_default", false).Error; err != nil {
			return nil, fmt.Errorf("更新默认配置状态失败: %v", err)
		}
	}

	// 更新配置
	config.ModelName = modelName
	config.Temperature = temperature
	config.MaxTokens = maxTokens
	config.Provider = provider
	config.IsDefault = isDefault

	if err := s.DB.Save(&config).Error; err != nil {
		return nil, fmt.Errorf("更新配置失败: %v", err)
	}

	return &config, nil
}

// DeleteAIConfig 删除AI配置
func (s *AIService) DeleteAIConfig(configID, userID uint) error {
	// 获取配置
	var config models.AIConfig
	if err := s.DB.First(&config, configID).Error; err != nil {
		return fmt.Errorf("配置不存在")
	}

	// 验证配置所有权
	if config.UserID != userID {
		return fmt.Errorf("无权删除此配置")
	}

	// 禁止删除默认配置
	if config.IsDefault {
		return fmt.Errorf("不能删除默认配置，请先将其他配置设为默认")
	}

	// 删除配置
	if err := s.DB.Delete(&config).Error; err != nil {
		return fmt.Errorf("删除配置失败: %v", err)
	}

	return nil
}

// 辅助方法 ---------------------------------------------------------

// getOrCreateSession 获取或创建会话
func (s *AIService) getOrCreateSession(userID, sessionID uint, message string) (*models.ChatSession, error) {
	var session models.ChatSession

	if sessionID > 0 {
		// 尝试获取现有会话
		if err := s.DB.First(&session, sessionID).Error; err != nil {
			return nil, fmt.Errorf("会话不存在")
		}

		// 验证所有权
		if session.UserID != userID {
			return nil, fmt.Errorf("无权访问此会话")
		}
	} else {
		// 创建新会话
		title := message
		if len([]rune(message)) > 30 {
			title = string([]rune(message)[:30])
		}

		session = models.ChatSession{
			UserID:      userID,
			Title:       title,
			LastMessage: message,
		}
		if err := s.DB.Create(&session).Error; err != nil {
			return nil, fmt.Errorf("创建会话失败: %v", err)
		}
	}

	return &session, nil
}

// getSessionMessages 获取会话的历史消息
func (s *AIService) getSessionMessages(sessionID uint) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	if err := s.DB.Where("session_id = ?", sessionID).Order("created_at asc").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// buildAIMessages 构建AI请求消息列表
func (s *AIService) buildAIMessages(messages []models.ChatMessage, newMessage string, knowledgeIDs []uint, userID uint) []ChatMessage {
	aiMessages := []ChatMessage{}

	// 添加系统消息
	aiMessages = append(aiMessages, ChatMessage{
		Role:    "system",
		Content: global.DefaultSystemPrompt,
	})

	// 添加知识库内容到系统提示（如果有）
	if len(knowledgeIDs) > 0 {
		knowledgeContent := s.getKnowledgeContent(knowledgeIDs, userID)
		if knowledgeContent != "" {
			aiMessages[0].Content += "\n\n以下是一些你可以参考的知识：\n" + knowledgeContent
		}
	}

	// 添加历史消息
	for _, msg := range messages {
		aiMessages = append(aiMessages, ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 添加用户最新消息
	aiMessages = append(aiMessages, ChatMessage{
		Role:    "user",
		Content: newMessage,
	})

	return aiMessages
}

// getKnowledgeContent 获取知识库内容
func (s *AIService) getKnowledgeContent(knowledgeIDs []uint, userID uint) string {
	var knowledgeContent string

	// 获取用户所有可用知识库文件
	var knowledgeFiles []models.KnowledgeFile
	if err := s.DB.Where("id IN ? AND user_id = ? AND status = ?", knowledgeIDs, userID, "completed").Find(&knowledgeFiles).Error; err != nil {
		return ""
	}

	// 对于每个知识库文件，获取其向量存储内容
	for _, file := range knowledgeFiles {
		var vectors []models.KnowledgeVectorStore
		if err := s.DB.Where("file_id = ?", file.ID).Find(&vectors).Error; err != nil {
			continue
		}

		// 向系统消息添加知识内容
		for _, vector := range vectors {
			knowledgeContent += vector.Text + "\n"
		}
	}

	return knowledgeContent
}

// ChunkText 将文本分块
func (s *AIService) ChunkText(text string, chunkSize int) []string {
	var chunks []string
	textRunes := []rune(text)
	textLen := len(textRunes)

	for i := 0; i < textLen; i += chunkSize {
		end := i + chunkSize
		if end > textLen {
			end = textLen
		}
		chunks = append(chunks, string(textRunes[i:end]))
	}
	return chunks
}
