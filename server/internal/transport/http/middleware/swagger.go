package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwaggerConfig Swagger UI 路由配置
type SwaggerConfig struct {
	BasePath                 string // Swagger 路由基础路径
	URL                      string // Swagger JSON 文档 URL
	DefaultModelsExpandDepth int    // 默认模型展开深度,-1 表示不展开
}

// DefaultSwaggerConfig 默认 Swagger UI 配置
func DefaultSwaggerConfig() SwaggerConfig {
	return SwaggerConfig{
		BasePath:                 "/swagger",
		URL:                      "/swagger/doc.json",
		DefaultModelsExpandDepth: -1,
	}
}

// RegisterSwagger 注册 Swagger UI 路由(仅开发环境)
// 注意：这不是中间件，而是路由注册辅助函数
func RegisterSwagger(engine *gin.Engine, config SwaggerConfig) {
	// 仅在 Debug 模式下注册 Swagger 路由
	if gin.Mode() != gin.DebugMode {
		return
	}

	log.Println("Swagger UI available at http://localhost:8080" + config.BasePath + "/index.html")

	// 注册 Swagger 路由组
	swaggerGroup := engine.Group(config.BasePath)
	{
		swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL(config.URL),
			ginSwagger.DefaultModelsExpandDepth(config.DefaultModelsExpandDepth),
		))
	}
}
