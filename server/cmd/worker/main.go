package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	crawlapp "github.com/shenfay/go-react-admin/internal/app/crawl"
	"github.com/shenfay/go-react-admin/internal/infra/config"
	"github.com/shenfay/go-react-admin/internal/infra/mail"
	"github.com/shenfay/go-react-admin/internal/infra/messaging"
	"github.com/shenfay/go-react-admin/internal/infra/repository"
	workerhandlers "github.com/shenfay/go-react-admin/internal/transport/worker/handlers"
	"github.com/shenfay/go-react-admin/pkg/constants"
	"github.com/shenfay/go-react-admin/pkg/logger"
)

func main() {
	// 1. 加载配置
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	cfg, err := config.Load(env)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. 初始化日志系统
	if err := logger.Init(cfg.Logger.Level); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Asynq Worker...")

	// 3. 初始化 Redis 客户端和数据库
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis: ", err)
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	logger.Info("Redis connection established")

	// 初始化数据库
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database: ", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Info("Database connection established")

	// 4. 初始化仓储
	operationLogRepo := repository.NewOperationLogRepository(db)

	// 抓取模块（Go↔Python Stream 集成）：调度器入队 Redis Stream，Python 消费执行
	crawlSourceRepo := repository.NewCrawlSourceRepository(db)
	crawlArticleRepo := repository.NewCrawlArticleRepository(db)
	crawlTaskRepo := repository.NewCrawlTaskRunRepository(db)
	crawlSvc := crawlapp.NewService(crawlSourceRepo, crawlArticleRepo, crawlTaskRepo, db, redisClient, cfg.Scraper)

	// 5. 创建处理器
	operationLogHandler := workerhandlers.NewOperationLogHandler(operationLogRepo)

	// 邮件发送处理器
	emailSender := mail.NewNoopSender()
	emailHandler := workerhandlers.NewSendEmailHandler(emailSender)

	// 6. 注册 Asynq 任务处理器
	mux := asynq.NewServeMux()

	// 抓取定时调度：cron 到点 → 入队 Redis Stream（crawl:task:dispatch）
	mux.HandleFunc(crawlapp.SchedulerTaskType, crawlapp.ScheduleHandler(crawlSvc))
	crawlapp.StartScheduler(asynq.RedisClientOpt{
		Addr:     cfg.Asynq.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}, crawlSvc)

	// Stream 消费者：消费 Python Worker 回传的文章和事件
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	go crawlSvc.StartArticleConsumer(workerCtx)
	go crawlSvc.StartEventConsumer(workerCtx)
	logger.Info("Stream consumers started (article + event)")

	// 从事件注册表获取所有路由到 logs 队列的事件类型（单一真相来源）
	for _, eventName := range messaging.LogEventTypes() {
		mux.HandleFunc(string(eventName), operationLogHandler.ProcessTask)
	}

	// AsynqTaskOperationLog 通过 InProcessBus → Bridge 入队，由于事件名与任务名不同，单独注册路由
	mux.HandleFunc(string(constants.AsynqTaskOperationLog), operationLogHandler.ProcessTask)

	// 发送邮件事件（auth:send_email 通过 Bridge 入队，事件名与任务名相同）
	mux.HandleFunc(string(constants.EventSendEmail), emailHandler.ProcessTask)

	// 7. 创建 Asynq 服务器
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Asynq.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		},
		asynq.Config{
			Concurrency:    cfg.Asynq.Concurrency,
			Queues:         cfg.Asynq.Queues,
			StrictPriority: true,
		},
	)

	logger.Info("Asynq server created with concurrency=", cfg.Asynq.Concurrency)

	// 8. 启动 Worker
	go func() {
		logger.Info("Starting Asynq Worker processor...")
		if err := srv.Run(mux); err != nil {
			logger.Error("Failed to run Asynq server: ", err)
			log.Fatalf("Failed to run Asynq server: %v", err)
		}
	}()

	// 9. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")

	// 10. 优雅关闭
	srv.Shutdown()
	logger.Info("Worker stopped gracefully")
}
