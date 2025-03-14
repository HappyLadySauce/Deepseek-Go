package middlewares

import (
	"net/http"
	"SmartDecision/config"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
)

// 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", config.Config.Cors.AllowOrigins[0]) // 允许的前端域名
		c.Writer.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(config.Config.Cors.AllowCredentials)) // 允许携带凭证
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(config.Config.Cors.AllowMethods, ", ")) // 允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(config.Config.Cors.AllowHeaders, ", ")) // 允许的请求头
		// 预检请求缓存24小时(24小时 = 86400秒)
		c.Writer.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.Config.Cors.MaxAge))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

