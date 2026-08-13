package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig 默认 CORS 配置
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           3600,
	}
}

// CORSMiddleware CORS 中间件
func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 非跨域请求（无 Origin 头）直接放行
		if origin == "" {
			c.Next()
			return
		}

		// 检查来源是否允许
		allowed := false
		for _, o := range config.AllowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(403, gin.H{
				"code":    "CORS_ERROR",
				"message": "Origin not allowed",
			})
			return
		}

		// 设置 CORS 响应头
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
		c.Header("Access-Control-Allow-Credentials", strconv.FormatBool(config.AllowCredentials))
		c.Header("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatusJSON(204, nil)
			return
		}

		c.Next()
	}
}
