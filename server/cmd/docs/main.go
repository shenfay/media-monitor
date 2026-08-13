package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Go DDD Scaffold API
// @version         1.0
// @description     生产级 DDD 脚手架项目的 API 文档，包含用户认证、事件驱动等核心功能。
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.github.com/shenfay/go-react-admin
// @contact.email  support@example.com

// @license.name   MIT
// @license.url    https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 使用 JWT Token，格式：Bearer {token}

func main() {
	r := gin.Default()

	// 注册 Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Swagger UI available at http://localhost:8080/swagger/index.html")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
