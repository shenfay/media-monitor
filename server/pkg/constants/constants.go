package constants

// 系统常量
const (
	// 项目名称
	ProjectName = "Go DDD Scaffold"

	// API 版本
	APIVersion = "v1"
	APIPrefix  = "/api/" + APIVersion

	// 默认分页参数
	DefaultPageSize   = 20
	MaxPageSize       = 100
	DefaultPageNumber = 1

	// 时间格式
	TimeLayoutRFC3339 = "2006-01-02T15:04:05Z07:00"
	TimeLayoutDate    = "2006-01-02"
	TimeLayoutTime    = "15:04:05"
)

// Redis Key 前缀
const (
	RedisKeyPrefix            = "kiqi:"
	RedisKeyRefreshToken      = RedisKeyPrefix + "refresh_token:"
	RedisKeyAccessDevice      = RedisKeyPrefix + "access_device:" // access_token → device_token_id 映射
	RedisKeyUserSession       = RedisKeyPrefix + "user_session:"
	RedisKeyLoginAttempts     = RedisKeyPrefix + "login_attempts:"
	RedisKeyEmailVerification = RedisKeyPrefix + "email_verification:"
	RedisKeyPasswordReset     = RedisKeyPrefix + "password_reset:"
)

// EventName 领域事件名称类型
type EventName string

// QueueName 消息队列名称类型
type QueueName string

// AsynqTaskType Asynq 任务类型
type AsynqTaskType string

// 领域事件名称常量
const (
	EventOperationLog EventName = "operation.log"   // 统一操作日志事件
	EventSendEmail    EventName = "auth:send_email" // 发送邮件事件
)

// Asynq 任务类型
const (
	AsynqTaskSendVerificationEmail  = "auth:send_verification_email"
	AsynqTaskSendPasswordResetEmail = "auth:send_password_reset_email"
	AsynqTaskSendWelcomeEmail       = "auth:send_welcome_email"
	AsynqTaskLogUserRegistration    = "auth:log_user_registration"
	AsynqTaskLogLoginAttempt        = "auth:log_login_attempt"
	AsynqTaskCleanupExpiredTokens   = "auth:cleanup_expired_tokens"
	AsynqTaskOperationLog           = "log:operation"     // 统一操作日志任务类型
	AsynqTaskNotification           = "notification:send" // 消息通知任务类型
)

// 队列名称
const (
	QueueCritical     QueueName = "critical"
	QueueDefault      QueueName = "default"
	QueueLow          QueueName = "low"
	QueueLogs         QueueName = "logs"         // 操作日志专用队列
	QueueNotification QueueName = "notification" // 消息通知专用队列
)
