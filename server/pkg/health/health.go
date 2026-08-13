package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HealthStatus 健康状态
type HealthStatus string

const (
	StatusOK      HealthStatus = "ok"
	StatusError   HealthStatus = "error"
	StatusWarning HealthStatus = "warning"
)

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status      HealthStatus    `json:"status"`
	Timestamp   string          `json:"timestamp"` // RFC3339 格式
	Version     string          `json:"version,omitempty"`
	Environment string          `json:"environment,omitempty"`
	Checks      ComponentChecks `json:"checks"`
}

// ComponentChecks 组件检查
type ComponentChecks struct {
	Database *DatabaseHealth `json:"database,omitempty"`
	Redis    *RedisHealth    `json:"redis,omitempty"`
	Asynq    *AsynqHealth    `json:"asynq,omitempty"`
}

// DatabaseHealth 数据库健康
type DatabaseHealth struct {
	Status       HealthStatus `json:"status"`
	ResponseTime string       `json:"response_time_ms,omitempty"`
	Error        string       `json:"error,omitempty"`
}

// RedisHealth Redis 健康
type RedisHealth struct {
	Status       HealthStatus `json:"status"`
	ResponseTime string       `json:"response_time_ms,omitempty"`
	PingResult   string       `json:"ping_result,omitempty"`
	Error        string       `json:"error,omitempty"`
}

// AsynqHealth Asynq 消息队列健康
type AsynqHealth struct {
	Status    HealthStatus `json:"status"`
	QueueSize int64        `json:"queue_size,omitempty"`
	Error     string       `json:"error,omitempty"`
}

// Handler 健康检查 Handler
type Handler struct {
	version     string
	environment string
	db          *gorm.DB
	redis       *redis.Client
	asynq       *asynq.Inspector
}

// NewHandler 创建健康检查 Handler
func NewHandler(db *gorm.DB, redisClient *redis.Client, version, env string) *Handler {
	return &Handler{
		version:     version,
		environment: env,
		db:          db,
		redis:       redisClient,
		asynq:       nil, // 默认不检查 Asynq
	}
}

// NewHandlerWithAsynq 创建包含 Asynq 检查的健康检查 Handler
func NewHandlerWithAsynq(db *gorm.DB, redisClient *redis.Client, asynqInspector *asynq.Inspector, version, env string) *Handler {
	return &Handler{
		version:     version,
		environment: env,
		db:          db,
		redis:       redisClient,
		asynq:       asynqInspector,
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/health", h.HandleHealth)
	router.GET("/health/live", h.HandleLive)   // Liveness probe
	router.GET("/health/ready", h.HandleReady) // Readiness probe
}

// HandleHealth 完整健康检查
func (h *Handler) HandleHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	response := &HealthResponse{
		Status:      StatusOK,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Version:     h.version,
		Environment: h.environment,
		Checks:      ComponentChecks{},
	}

	// 检查数据库
	response.Checks.Database = h.checkDatabase(ctx)
	if response.Checks.Database.Status == StatusError {
		response.Status = StatusError
	}

	// 检查 Redis
	response.Checks.Redis = h.checkRedis(ctx)
	if response.Checks.Redis.Status == StatusError {
		response.Status = StatusError
	}

	// 检查 Asynq（如果配置了）
	if h.asynq != nil {
		response.Checks.Asynq = h.checkAsynq(ctx)
		if response.Checks.Asynq.Status == StatusError {
			response.Status = StatusError
		}
	}

	c.JSON(http.StatusOK, response)
}

// HandleLive Liveness 探针（只检查服务是否存活）
func (h *Handler) HandleLive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   h.version,
	})
}

// HandleReady Readiness 探针（检查是否准备好接收流量）
func (h *Handler) HandleReady(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// 检查数据库连接
	dbErr := h.checkDBConnection(ctx)

	// 检查 Redis 连接
	_, redisErr := h.checkRedisConnection(ctx)

	if dbErr != nil || redisErr != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "not_ready",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"database": func() string {
				if dbErr != nil {
					return "unhealthy"
				}
				return "healthy"
			}(),
			"redis": func() string {
				if redisErr != nil {
					return "unhealthy"
				}
				return "healthy"
			}(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   h.version,
	})
}

// checkDatabase 检查数据库健康
func (h *Handler) checkDatabase(ctx context.Context) *DatabaseHealth {
	start := time.Now()

	err := h.checkDBConnection(ctx)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		return &DatabaseHealth{
			Status: StatusError,
			Error:  err.Error(),
		}
	}

	health := &DatabaseHealth{
		Status:       StatusOK,
		ResponseTime: formatDuration(duration),
	}

	// 如果响应时间超过阈值，标记为 warning
	if duration > 1000 {
		health.Status = StatusWarning
	}

	return health
}

// checkDBConnection 检查数据库连接
func (h *Handler) checkDBConnection(ctx context.Context) error {
	// 使用 SQL 查询验证连接
	var result int
	if err := h.db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error; err != nil {
		return err
	}

	// 检查连接池状态
	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}

	stats := sqlDB.Stats()
	if stats.OpenConnections >= stats.MaxOpenConnections {
		// 连接池已满，返回 warning
		return nil
	}

	return nil
}

// checkRedis 检查 Redis 健康
func (h *Handler) checkRedis(ctx context.Context) *RedisHealth {
	start := time.Now()

	result, err := h.checkRedisConnection(ctx)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		return &RedisHealth{
			Status: StatusError,
			Error:  err.Error(),
		}
	}

	health := &RedisHealth{
		Status:       StatusOK,
		ResponseTime: formatDuration(duration),
		PingResult:   result,
	}

	// 如果响应时间超过阈值，标记为 warning
	if duration > 100 {
		health.Status = StatusWarning
	}

	return health
}

// checkRedisConnection 检查 Redis 连接
func (h *Handler) checkRedisConnection(ctx context.Context) (string, error) {
	// 使用 Ping 命令验证连接
	result, err := h.redis.Ping(ctx).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

// checkAsynq 检查 Asynq 消息队列健康
func (h *Handler) checkAsynq(ctx context.Context) *AsynqHealth {
	if h.asynq == nil {
		return nil
	}

	// 获取队列信息
	info, err := h.asynq.GetQueueInfo("default")
	if err != nil {
		return &AsynqHealth{
			Status: StatusError,
			Error:  fmt.Sprintf("failed to get queue info: %v", err),
		}
	}

	health := &AsynqHealth{
		Status:    StatusOK,
		QueueSize: int64(info.Size), // 待处理任务数
	}

	// 如果队列积压超过阈值，标记为 warning
	if health.QueueSize > 1000 {
		health.Status = StatusWarning
	}

	return health
}

// formatDuration 格式化持续时间
func formatDuration(ms int64) string {
	if ms < 10 {
		return "<10ms"
	} else if ms < 100 {
		return "<100ms"
	} else if ms < 1000 {
		return ">100ms"
	}
	return ">1s"
}
