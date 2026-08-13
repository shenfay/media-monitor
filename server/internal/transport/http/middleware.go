package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-react-admin/internal/infra/config"
	"github.com/shenfay/go-react-admin/internal/transport/http/middleware"
	"github.com/shenfay/go-react-admin/pkg/metrics"
)

// Middlewares 注册全局中间件
func Middlewares(engine *gin.Engine, m *metrics.Metrics, corsCfg config.CORSConfig) {
	// CORS 中间件
	engine.Use(middleware.CORSMiddleware(middleware.CORSConfig{
		AllowedOrigins:   corsCfg.AllowedOrigins,
		AllowedMethods:   corsCfg.AllowedMethods,
		AllowedHeaders:   corsCfg.AllowedHeaders,
		AllowCredentials: corsCfg.AllowCredentials,
		MaxAge:           corsCfg.MaxAge,
	}))

	// Prometheus 监控中间件
	engine.Use(middleware.PrometheusMiddleware(m))

	// Recovery 中间件(必须)
	engine.Use(gin.Recovery())
}
