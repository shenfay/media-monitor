package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-react-admin/internal/transport/http/response"
	"github.com/shenfay/go-react-admin/pkg/logger"
	"github.com/shenfay/go-react-admin/pkg/utils"
	"go.uber.org/zap"
)

const (
	// TraceIDHeader HTTP Header
	TraceIDHeader = "X-Trace-ID"
)

// ContextKey 自定义 Context 键类型
type ContextKey string

const (
	// TraceIDContextKey trace_id 在 context 中的键
	TraceIDContextKey ContextKey = "trace_id"
)

// TraceID 链路追踪 ID 生成与传递中间件
// 为每个请求生成或传递 trace_id，用于日志关联和分布式追踪
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先使用客户端传入的 trace_id（用于跨服务链路追踪）
		traceID := c.GetHeader(TraceIDHeader)
		if traceID == "" {
			// 生成新的 trace_id（ULID 格式）
			traceID = utils.GenerateID()
		}

		// 存入 Gin Context
		c.Set(response.TraceIDKey, traceID)

		// 响应头中也包含 trace_id
		c.Header(TraceIDHeader, traceID)

		// 注入到 Request Context（供 Service 层使用）
		ctx := context.WithValue(c.Request.Context(), TraceIDContextKey, traceID)
		c.Request = c.Request.WithContext(ctx)

		// 记录请求开始时间
		startTime := time.Now()

		// 记录请求开始日志（包含 trace_id）
		logger.Info("Request started",
			"trace_id", traceID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
		)

		c.Next()

		// 计算延迟
		latency := time.Since(startTime)

		// 记录请求结束日志（包含 trace_id 和响应状态）
		logger.Info("Request completed",
			"trace_id", traceID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", latency,
		)
	}
}

// GetTraceID 从 Gin Context 获取 trace_id
// Deprecated: 使用 response.GetTraceID 替代
func GetTraceID(c *gin.Context) string {
	return response.GetTraceID(c)
}

// GetTraceIDFromContext 从标准 Context 获取 trace_id（供 Service 层使用）
func GetTraceIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(TraceIDContextKey).(string); ok {
		return id
	}
	return ""
}

// Logger 获取带 trace_id 的 logger（供 Handler/Service 层使用）
func Logger(c *gin.Context) *zap.SugaredLogger {
	traceID := GetTraceID(c)
	if traceID != "" {
		return logger.Logger.With("trace_id", traceID)
	}
	return logger.Logger
}
