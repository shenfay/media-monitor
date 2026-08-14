package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 应用程序配置
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	Asynq       AsynqConfig       `mapstructure:"asynq"`
	Logger      LoggerConfig      `mapstructure:"logger"`
	CORS        CORSConfig        `mapstructure:"cors"`
	Auth        AuthConfig        `mapstructure:"auth"`
	Email       EmailConfig       `mapstructure:"email"`
	WebSocket   WebSocketConfig   `mapstructure:"websocket"`
	Encryption  EncryptionConfig  `mapstructure:"encryption"`
	RateLimit   RateLimitConfig   `mapstructure:"rate_limit"`
	Scraper     ScraperConfig     `mapstructure:"scraper"`
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

// EncryptionConfig 加密配置
type EncryptionConfig struct {
	Key string `mapstructure:"key"` // AES-256-GCM 加密密钥（至少 32 字符）
}

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	Enabled  bool              `mapstructure:"enabled"`
	General  RateLimitItemConfig `mapstructure:"general"`
	Login    RateLimitItemConfig `mapstructure:"login"`
	Register RateLimitItemConfig `mapstructure:"register"`
}

// RateLimitItemConfig 速率限制项配置
type RateLimitItemConfig struct {
	Rate  int `mapstructure:"rate"`
	Burst int `mapstructure:"burst"`
}

// ScraperConfig Python 抓取服务对接配置（Go↔Python Stream 集成缝）
type ScraperConfig struct {
	Stream        string `mapstructure:"stream"`         // Redis Stream 任务分发队列名，默认 crawl:task:dispatch
	ArticleStream string `mapstructure:"article_stream"` // Redis Stream 文章回传队列名，默认 crawl:article:ingest
	EventStream   string `mapstructure:"event_stream"`   // Redis Stream 事件回传队列名，默认 crawl:task:event
	ConsumerGroup string `mapstructure:"consumer_group"` // Stream 消费者组名，默认 crawl:go
}

// Load 加载配置
// 配置优先级：环境变量 > .env 文件 > YAML 配置文件 > 默认值
func Load(env string) (*Config, error) {
	// 1. 加载 .env 文件（如果存在），将变量注入进程环境变量
	// 优先加载 configs/.env，其次加载项目根目录的 .env
	// 注意：godotenv.Load 不会覆盖已存在的环境变量
	for _, envFile := range []string{"configs/.env", ".env"} {
		if err := godotenv.Load(envFile); err == nil {
			log.Printf("Loaded .env file: %s", envFile)
			break
		}
	}

	// 2. 读取 YAML 配置文件
	viper.SetConfigFile(fmt.Sprintf("configs/%s.yaml", env))
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 3. 用环境变量覆盖配置（.env 加载后已存在于 os.Environ 中）
	overrideFromEnv(&cfg)

	return &cfg, nil
}

// overrideFromEnv 使用环境变量覆盖配置值
// 环境变量命名规则：APP_<SECTION>_<KEY>，如 APP_DATABASE_HOST
func overrideFromEnv(cfg *Config) {
	// Server
	if v := os.Getenv("APP_SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("APP_SERVER_MODE"); v != "" {
		cfg.Server.Mode = v
	}

	// Database
	if v := os.Getenv("APP_DATABASE_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("APP_DATABASE_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Database.Port = port
		}
	}
	if v := os.Getenv("APP_DATABASE_NAME"); v != "" {
		cfg.Database.Name = v
	}
	if v := os.Getenv("APP_DATABASE_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("APP_DATABASE_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("APP_DATABASE_SSL_MODE"); v != "" {
		cfg.Database.SSLMode = v
	}

	// Redis
	if v := os.Getenv("APP_REDIS_ADDR"); v != "" {
		cfg.Redis.Addr = v
	}
	if v := os.Getenv("APP_REDIS_PASSWORD"); v != "" {
		cfg.Redis.Password = v
	}
	if v := os.Getenv("APP_REDIS_DB"); v != "" {
		if db, err := strconv.Atoi(v); err == nil {
			cfg.Redis.DB = db
		}
	}

	// JWT
	if v := os.Getenv("APP_JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("APP_JWT_ACCESS_EXPIRE"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.JWT.AccessExpire = d
		}
	}
	if v := os.Getenv("APP_JWT_REFRESH_EXPIRE"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.JWT.RefreshExpire = d
		}
	}

	// Logger
	if v := os.Getenv("APP_LOGGING_LEVEL"); v != "" {
		cfg.Logger.Level = v
	}
	if v := os.Getenv("APP_LOGGING_FORMAT"); v != "" {
		cfg.Logger.Format = v
	}

	// Encryption
	if v := os.Getenv("APP_ENCRYPTION_KEY"); v != "" {
		cfg.Encryption.Key = v
	}
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

	// Encryption
	viper.SetDefault("encryption.key", "your-encryption-key-change-in-production")

	// Rate Limit
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.general.rate", 60)
	viper.SetDefault("rate_limit.general.burst", 100)
	viper.SetDefault("rate_limit.login.rate", 5)
	viper.SetDefault("rate_limit.login.burst", 10)

	// Scraper（Python 抓取服务对接）
	viper.SetDefault("scraper.stream", "crawl:task:dispatch")
	viper.SetDefault("scraper.article_stream", "crawl:article:ingest")
	viper.SetDefault("scraper.event_stream", "crawl:task:event")
	viper.SetDefault("scraper.consumer_group", "crawl:go")
	viper.SetDefault("rate_limit.register.rate", 10)
	viper.SetDefault("rate_limit.register.burst", 20)
}
