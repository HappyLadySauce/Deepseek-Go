package controller

import (
	"Deepseek-Go/utils/ai"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AI配置控制器
type AIConfigController struct {
	DB        *gorm.DB
	AIService *ai.AIService
}

// AI配置请求
type AIConfigCreateRequest struct {
	ModelName   string  `json:"model_name" binding:"required"`
	Temperature float64 `json:"temperature" binding:"required,min=0,max=1"`
	MaxTokens   int     `json:"max_tokens" binding:"required,min=1,max=4096"`
	Provider    string  `json:"provider" binding:"required,oneof=deepseek kimi"`
	IsDefault   bool    `json:"is_default"`
}

// 构造函数
func NewAIConfigController(db *gorm.DB) *AIConfigController {
	return &AIConfigController{
		DB:        db,
		AIService: ai.NewAIService(db),
	}
}

// CreateConfig 创建AI配置
func (ac *AIConfigController) CreateConfig(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 解析请求体
	var req AIConfigCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数: " + err.Error()})
		return
	}

	// 调用服务创建配置
	config, err := ac.AIService.CreateAIConfig(
		userID.(uint),
		req.ModelName,
		req.Temperature,
		req.MaxTokens,
		req.Provider,
		req.IsDefault,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建AI配置成功",
		"data":    config,
	})
}

// GetConfigs 获取用户的所有AI配置
func (ac *AIConfigController) GetConfigs(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 调用服务获取配置列表
	configs, err := ac.AIService.GetAIConfigs(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取配置列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取AI配置成功",
		"data":    configs,
	})
}

// GetConfig 获取单个AI配置
func (ac *AIConfigController) GetConfig(c *gin.Context) {
	// 获取配置ID
	configID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 调用服务获取配置
	config, err := ac.AIService.GetAIConfig(uint(configID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取配置成功",
		"data":    config,
	})
}

// UpdateConfig 更新AI配置
func (ac *AIConfigController) UpdateConfig(c *gin.Context) {
	// 获取配置ID
	configID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 解析请求体
	var req AIConfigCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数: " + err.Error()})
		return
	}

	// 调用服务更新配置
	config, err := ac.AIService.UpdateAIConfig(
		uint(configID),
		userID.(uint),
		req.ModelName,
		req.Temperature,
		req.MaxTokens,
		req.Provider,
		req.IsDefault,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新配置成功",
		"data":    config,
	})
}

// DeleteConfig 删除AI配置
func (ac *AIConfigController) DeleteConfig(c *gin.Context) {
	// 获取配置ID
	configID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 调用服务删除配置
	if err := ac.AIService.DeleteAIConfig(uint(configID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除配置成功",
	})
}

// GetDefaultConfig 获取用户默认AI配置
func (ac *AIConfigController) GetDefaultConfig(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 调用服务获取默认配置
	config, err := ac.AIService.GetDefaultAIConfig(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取默认配置失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取默认配置成功",
		"data":    config,
	})
}

// GetAvailableModels 获取可用的AI模型列表
func (ac *AIConfigController) GetAvailableModels(c *gin.Context) {
	// DeepSeek模型列表
	deepseekModels := []gin.H{
		{"name": "deepseek-chat", "provider": "deepseek", "description": "基础模型"},
		{"name": "deepseek-reasoner", "provider": "deepseek", "description": "深度思考模型"},
	}

	// Kimi模型列表
	kimiModels := []gin.H{
		{"name": "moonshot-v1-8k", "provider": "kimi", "description": "基础模型，支持8K上下文"},
		{"name": "moonshot-v1-32k", "provider": "kimi", "description": "基础模型，支持32K上下文"},
		{"name": "moonshot-v1-128k", "provider": "kimi", "description": "基础模型，支持128K上下文"},
		{"name": "moonshot-v1-auto", "provider": "kimi", "description": "自动选择模型，根据上下文长度"},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取可用模型列表成功",
		"data": gin.H{
			"deepseek": deepseekModels,
			"kimi":     kimiModels,
		},
	})
}
