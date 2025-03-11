package controller

import (
	"Deepseek-Go/global"
	"Deepseek-Go/utils/ai"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 知识库控制器
type KnowledgeController struct {
	DB        *gorm.DB
	AIService *ai.AIService
}

// 文件上传响应
type UploadResponse struct {
	ID        uint      `json:"id"`
	FileName  string    `json:"file_name"`
	FileSize  int64     `json:"file_size"`
	FileType  string    `json:"file_type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// 构造函数
func NewKnowledgeController(db *gorm.DB) *KnowledgeController {
	return &KnowledgeController{
		DB:        db,
		AIService: ai.NewAIService(db),
	}
}

// 上传知识库文件
func (kc *KnowledgeController) UploadFile(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取上传文件失败: " + err.Error()})
		return
	}
	defer file.Close()

	// 校验文件类型
	fileExt := strings.ToLower(strings.TrimSpace(header.Filename))
	fileExt = strings.ToLower(strings.TrimSpace(c.Request.FormValue("fileExt")))
	if !isAllowedFileType(fileExt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型，仅支持pdf, docx, doc, txt, md文件"})
		return
	}

	// 校验文件大小
	if header.Size > global.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过10MB"})
		return
	}

	// 保存文件到磁盘
	filePath, fileType, err := kc.AIService.SaveKnowledgeFile(file, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}

	// 保存文件信息到数据库
	knowledgeFile, err := kc.AIService.UploadKnowledgeFile(userID.(uint), header.Filename, header.Size, fileType, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文件上传成功",
		"data": UploadResponse{
			ID:        knowledgeFile.ID,
			FileName:  knowledgeFile.FileName,
			FileSize:  knowledgeFile.FileSize,
			FileType:  knowledgeFile.FileType,
			Status:    knowledgeFile.Status,
			CreatedAt: knowledgeFile.CreatedAt,
		},
	})
}

// 获取知识库文件列表
func (kc *KnowledgeController) GetFiles(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 调用服务获取文件列表
	files, count, err := kc.AIService.GetKnowledgeFiles(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取文件列表成功",
		"data": gin.H{
			"total":    count,
			"page":     page,
			"pageSize": pageSize,
			"files":    files,
		},
	})
}

// 获取知识库文件详情
func (kc *KnowledgeController) GetFile(c *gin.Context) {
	// 获取文件ID
	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件ID"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 调用服务获取文件详情
	file, vectorCount, err := kc.AIService.GetKnowledgeFile(uint(fileID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构造响应
	response := gin.H{
		"id":           file.ID,
		"file_name":    file.FileName,
		"file_size":    file.FileSize,
		"file_type":    file.FileType,
		"status":       file.Status,
		"processed_at": file.ProcessedAt,
		"created_at":   file.CreatedAt,
		"vector_count": vectorCount,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取文件成功",
		"data":    response,
	})
}

// 删除知识库文件
func (kc *KnowledgeController) DeleteFile(c *gin.Context) {
	// 获取文件ID
	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件ID"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的访问"})
		return
	}

	// 调用服务删除文件
	err = kc.AIService.DeleteKnowledgeFile(uint(fileID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除文件成功",
	})
}

// 辅助函数：检查文件类型是否允许
func isAllowedFileType(fileExt string) bool {
	return global.AllowedFileTypes[strings.ToLower(fileExt)]
}
