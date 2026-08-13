package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用程序配置
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Asynq     AsynqConfig     `mapstructure:"asynq"`
	Logger    LoggerConfig    `mapstructure:"logger"`
	CORS      CORSConfig      `mapstructure:"cors"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Email     EmailConfig     `mapstructure:"email"`
	WebSocket WebSocketConfig `mapstructure:"websocket"`
}

// ServerConfig HTTP 服务器配置
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"` // debug, release, test
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig 数据库连接配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Name            string        `mapstructure:"name"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// DSN 返回数据库连接字符串
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// RedisConfig Redis 连接配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig JWT 令牌配置
type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	AccessExpire  time.Duration `mapstructure:"access_expire"`
	RefreshExpire time.Duration `mapstructure:"refresh_expire"`
	Issuer        string        `mapstructure:"issuer"`
}

// AsynqConfig Asynq 配置
type AsynqConfig struct {
	Addr        string         `mapstructure:"addr"`
	Concurrency int            `mapstructure:"concurrency"`
	Queues      map[string]int `mapstructure:"queues"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"` // json, console
	OutputPath string `mapstructure:"output_path"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	MaxDevicesPerUser       int           `mapstructure:"max_devices_per_user"`
	AutoRevokeOldest        bool          `mapstructure:"auto_revoke_oldest"`
	EmailVerificationExpire time.Duration `mapstructure:"email_verification_expire"`
	PasswordResetExpire     time.Duration `mapstructure:"password_reset_expire"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	From                     string `mapstructure:"from"`
	VerificationURLTemplate  string `mapstructure:"verification_url_template"`
	PasswordResetURLTemplate string `mapstructure:"password_reset_url_template"`
	SMTPHost                 string `mapstructure:"smtp_host"`
	SMTPPort                 int    `mapstructure:"smtp_port"`
	SMTPUsername             string `mapstructure:"smtp_username"`
	SMTPPassword             string `mapstructure:"smtp_password"`
}

// CORSConfig 跨域配置
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// WebSocketConfig WebSocket 实时推送配置
type WebSocketConfig struct {
	Enabled bool `mapstructure:"enabled"` // 是否启用 WebSocket 端点
}

// Load 加载配置
func Load(env string) (*Config, error) {
	// 设置配置文件路径
	viper.SetConfigFile(fmt.Sprintf("configs/%s.yaml", env))

	// 允许通过环境变量覆盖配置
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaults 设置默认值
func setDefaults() {
	// Server
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 30*time.Second)
	viper.SetDefault("server.write_timeout", 30*time.Second)
	viper.SetDefault("server.idle_timeout", 60*time.Second)

	// Database
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "kiqi")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", 5*time.Minute)

	// Redis
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	// JWT
	viper.SetDefault("jwt.secret", "your-jwt-secret-key-change-in-production")
	viper.SetDefault("jwt.access_expire", 30*time.Minute)
	viper.SetDefault("jwt.refresh_expire", 7*24*time.Hour)
	viper.SetDefault("jwt.issuer", "kiqi")

	// Asynq
	viper.SetDefault("asynq.addr", "localhost:6379")
	viper.SetDefault("asynq.concurrency", 10)
	viper.SetDefault("asynq.queues", map[string]int{
		"critical":     6,
		"default":      3,
		"logs":         4,
		"notification": 4,
		"low":          1,
	})

	// Logger
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("logger.format", "console")
	viper.SetDefault("logger.output_path", "stdout")

	// Auth
	viper.SetDefault("auth.max_devices_per_user", 5)
	viper.SetDefault("auth.auto_revoke_oldest", false)
	viper.SetDefault("auth.email_verification_expire", 24*time.Hour)
	viper.SetDefault("auth.password_reset_expire", 1*time.Hour)

	// Email
	viper.SetDefault("email.from", "noreply@example.com")
	viper.SetDefault("email.verification_url_template", "http://localhost:3000/verify-email?token=%s&user_id=%s")
	viper.SetDefault("email.password_reset_url_template", "http://localhost:3000/reset-password?token=%s&user_id=%s")
	viper.SetDefault("email.smtp_host", "")
	viper.SetDefault("email.smtp_port", 587)
	viper.SetDefault("email.smtp_username", "")
	viper.SetDefault("email.smtp_password", "")

	// CORS
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:5173"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Authorization", "Content-Type", "X-Requested-With"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", 3600)

	// WebSocket
	viper.SetDefault("websocket.enabled", true)
}
