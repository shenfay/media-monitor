package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics 指标集合
type Metrics struct {
	// HTTP 请求指标
	HTTPRequestTotal     *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPRequestsInFlight prometheus.Gauge

	// 数据库操作指标
	DBQueryTotal         prometheus.Counter
	DBQueryDuration      *prometheus.HistogramVec
	DBConnections        prometheus.Gauge
	DBConnectionsMax     prometheus.Gauge
	DBConnectionWaitTime *prometheus.HistogramVec
	DBSlowQueriesTotal   prometheus.Counter

	// Redis 操作指标
	RedisCommandTotal      prometheus.Counter
	RedisCommandDuration   *prometheus.HistogramVec
	RedisConnectionsActive prometheus.Gauge
	RedisMemoryUsed        prometheus.Gauge
	RedisHitRate           prometheus.Gauge

	// Asynq 消息队列指标
	AsynqQueueSize      *prometheus.GaugeVec
	AsynqTasksProcessed *prometheus.CounterVec
	AsynqTaskDuration   *prometheus.HistogramVec
	AsynqWorkersActive  *prometheus.GaugeVec
	AsynqRetryTotal     *prometheus.CounterVec

	// 认证相关指标
	AuthAttemptsTotal   *prometheus.CounterVec
	AuthSuccessTotal    *prometheus.CounterVec
	AuthFailureTotal    *prometheus.CounterVec
	ActiveUsers         prometheus.Gauge
	TokenRefreshesTotal prometheus.Counter

	// 业务指标
	UserRegistrationsTotal prometheus.Counter
	EmailsSentTotal        prometheus.Counter
}

// NewMetrics 创建指标集合
func NewMetrics(registry prometheus.Registerer) *Metrics {
	m := &Metrics{
		// HTTP 请求指标
		HTTPRequestTotal: promauto.With(registry).NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"status"},
		),
		HTTPRequestDuration: promauto.With(registry).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestsInFlight: promauto.With(registry).NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Number of HTTP requests currently being processed",
			},
		),

		// 数据库操作指标
		DBQueryTotal: promauto.With(registry).NewCounter(
			prometheus.CounterOpts{
				Name: "db_queries_total",
				Help: "Total number of database queries",
			},
		),
		DBQueryDuration: promauto.With(registry).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "table"},
		),
		DBConnections: promauto.With(registry).NewGauge(
			prometheus.GaugeOpts{
				Name: "db_connections_open",
				Help: "Current number of open database connections",
			},
		),
		DBConnectionsMax: promauto.With(registry).NewGauge(
			prometheus.GaugeOpts{
				Name: "db_connections_max",
				Help: "Maximum number of database connections",
			},
		),
		DBConnectionWaitTime: promauto.With(registry).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_connection_wait_time_seconds",
				Help:    "Time spent waiting for a database connection",
				Buckets: prometheus.DefBuckets,
			},
			[]string{},
		),
		DBSlowQueriesTotal: promauto.With(registry).NewCounter(
			prometheus.CounterOpts{
				Name: "db_slow_queries_total",
				Help: "Total number of slow database queries (>1s)",
			},
		),

		// Redis 操作指标
		RedisCommandTotal: promauto.With(registry).NewCounter(
			prometheus.CounterOpts{
				Name: "redis_commands_total",
				Help: "Total number of Redis commands",
			},
		),
		RedisCommandDuration: promauto.With(registry).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "redis_command_duration_seconds",
				Help:    "Redis command duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"command"},
		),
		RedisConnectionsActive: promauto.With(registry).NewGauge(
			prometheus.GaugeOpts{
				Name: "redis_connections_active",
				Help: "Current number of active Redis connections",
			},
		),
		RedisMemoryUsed: promauto.With(registry).NewGauge(
			prometheus.GaugeOpts{
				Name: "redis_memory_used_bytes",
				Help: "Current Redis memory usage in bytes",
			},
		),
		RedisHitRate: promauto.With(registry).NewGauge(
			prometheus.GaugeOpts{
				Name: "redis_hit_rate",
				Help: "Redis cache hit rate (0-1)",
			},
		),

		// Asynq 消息队列指标
		AsynqQueueSize: promauto.With(registry).NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "asynq_queue_size",
				Help: "Current number of tasks in queue",
			},
			[]string{"queue"},
		),
		AsynqTasksProcessed: promauto.With(registry).NewCounterVec(
			prometheus.CounterOpts{
				Name: "asynq_tasks_processed_total",
				Help: "Total number of processed tasks",
			},
			[]string{"queue", "status"}, // status: success, failed
		),
		AsynqTaskDuration: promauto.With(registry).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "asynq_task_duration_seconds",
				Help:    "Task processing duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"queue", "type"},
		),
		AsynqWorkersActive: promauto.With(registry).NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "asynq_workers_active",
				Help: "Current number of active workers",
			},
			[]string{"queue"},
		),
		AsynqRetryTotal: promauto.With(registry).NewCounterVec(
			prometheus.CounterOpts{
				Name: "asynq_retries_total",
				Help: "Total number of task retries",
			},
			[]string{"queue", "type"},
		),

		// 认证相关指标
		AuthAttemptsTotal: promauto.With(registry).NewCounterVec(
			prometheus.CounterOpts{
				Name: "auth_attempts_total",
				Help: "Total number of authentication attempts",
			},
			[]string{"type", "status"},
		),
		AuthSuccessTotal: promauto.With(registry).NewCounterVec(
			prometheus.CounterOpts{
				Name: "auth_success_total",
				Help: "Total number of successful authentications",
			},
			[]string{"type"},
		),
		AuthFailureTotal: promauto.With(registry).NewCounterVec(
			prometheus.CounterOpts{
				Name: "auth_failure_total",
				Help: "Total number of failed authentications",
			},
			[]string{"type", "reason"},
		),
		ActiveUsers: promauto.With(registry).NewGauge(
			prometheus.GaugeOpts{
				Name: "active_users",
				Help: "Current number of active users",
			},
		),
		TokenRefreshesTotal: promauto.With(registry).NewCounter(
			prometheus.CounterOpts{
				Name: "token_refreshes_total",
				Help: "Total number of token refresh operations",
			},
		),

		// 业务指标
		UserRegistrationsTotal: promauto.With(registry).NewCounter(
			prometheus.CounterOpts{
				Name: "user_registrations_total",
				Help: "Total number of user registrations",
			},
		),
		EmailsSentTotal: promauto.With(registry).NewCounter(
			prometheus.CounterOpts{
				Name: "emails_sent_total",
				Help: "Total number of emails sent",
			},
		),
	}

	return m
}

// IncHTTPRequests 增加 HTTP 请求计数
func (m *Metrics) IncHTTPRequests(status string) {
	m.HTTPRequestTotal.WithLabelValues(status).Inc()
}

// ObserveHTTPDuration 记录 HTTP 请求耗时
func (m *Metrics) ObserveHTTPDuration(method, path, status string, duration float64) {
	m.HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration)
}

// IncDBQuery 增加数据库查询计数
func (m *Metrics) IncDBQuery(operation, table string) {
	m.DBQueryTotal.Inc()
}

// ObserveDBQueryDuration 记录数据库查询耗时
func (m *Metrics) ObserveDBQueryDuration(operation, table string, duration float64) {
	m.DBQueryDuration.WithLabelValues(operation, table).Observe(duration)
	// 检测慢查询（>1s）
	if duration > 1.0 {
		m.DBSlowQueriesTotal.Inc()
	}
}

// UpdateDBConnections 更新数据库连接池状态
func (m *Metrics) UpdateDBConnections(open, max int) {
	m.DBConnections.Set(float64(open))
	m.DBConnectionsMax.Set(float64(max))
}

// ObserveDBConnectionWaitTime 记录等待连接的时间
func (m *Metrics) ObserveDBConnectionWaitTime(duration float64) {
	m.DBConnectionWaitTime.WithLabelValues().Observe(duration)
}

// IncRedisCommand 增加 Redis 命令计数
func (m *Metrics) IncRedisCommand(command string) {
	m.RedisCommandTotal.Inc()
}

// ObserveRedisCommandDuration 记录 Redis 命令耗时
func (m *Metrics) ObserveRedisCommandDuration(command string, duration float64) {
	m.RedisCommandDuration.WithLabelValues(command).Observe(duration)
}

// UpdateRedisConnections 更新 Redis 连接数
func (m *Metrics) UpdateRedisConnections(active int) {
	m.RedisConnectionsActive.Set(float64(active))
}

// UpdateRedisMemory 更新 Redis 内存使用
func (m *Metrics) UpdateRedisMemory(bytes int64) {
	m.RedisMemoryUsed.Set(float64(bytes))
}

// UpdateRedisHitRate 更新 Redis 命中率
func (m *Metrics) UpdateRedisHitRate(rate float64) {
	m.RedisHitRate.Set(rate)
}

// UpdateAsynqQueueSize 更新队列大小
func (m *Metrics) UpdateAsynqQueueSize(queue string, size int) {
	m.AsynqQueueSize.WithLabelValues(queue).Set(float64(size))
}

// IncAsynqTaskProcessed 增加任务处理计数
func (m *Metrics) IncAsynqTaskProcessed(queue, status string) {
	m.AsynqTasksProcessed.WithLabelValues(queue, status).Inc()
}

// ObserveAsynqTaskDuration 记录任务处理耗时
func (m *Metrics) ObserveAsynqTaskDuration(queue, taskType string, duration float64) {
	m.AsynqTaskDuration.WithLabelValues(queue, taskType).Observe(duration)
}

// UpdateAsynqWorkersActive 更新活跃 Worker 数
func (m *Metrics) UpdateAsynqWorkersActive(queue string, count int) {
	m.AsynqWorkersActive.WithLabelValues(queue).Set(float64(count))
}

// IncAsynqRetry 增加重试计数
func (m *Metrics) IncAsynqRetry(queue, taskType string) {
	m.AsynqRetryTotal.WithLabelValues(queue, taskType).Inc()
}

// IncAuthAttempt 记录认证尝试
func (m *Metrics) IncAuthAttempt(authType, status string) {
	m.AuthAttemptsTotal.WithLabelValues(authType, status).Inc()
}

// IncAuthSuccess 记录认证成功
func (m *Metrics) IncAuthSuccess(authType string) {
	m.AuthSuccessTotal.WithLabelValues(authType).Inc()
}

// IncAuthFailure 记录认证失败
func (m *Metrics) IncAuthFailure(authType, reason string) {
	m.AuthFailureTotal.WithLabelValues(authType, reason).Inc()
}

// IncUserRegistration 增加用户注册计数
func (m *Metrics) IncUserRegistration() {
	m.UserRegistrationsTotal.Inc()
}

// IncEmailSent 增加邮件发送计数
func (m *Metrics) IncEmailSent() {
	m.EmailsSentTotal.Inc()
}

// SetActiveUsers 设置活跃用户数
func (m *Metrics) SetActiveUsers(count int) {
	m.ActiveUsers.Set(float64(count))
}

// IncTokenRefresh 增加 Token 刷新计数
func (m *Metrics) IncTokenRefresh() {
	m.TokenRefreshesTotal.Inc()
}
