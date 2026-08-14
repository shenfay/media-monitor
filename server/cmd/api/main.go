package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/prometheus/client_golang/prometheus"

	// 导入生成的 Swagger 文档
	_ "github.com/shenfay/go-react-admin/api/swagger"

	"github.com/shenfay/go-react-admin/internal/app/admin"
	"github.com/shenfay/go-react-admin/internal/app/authentication"
	crawlapp "github.com/shenfay/go-react-admin/internal/app/crawl"
	crawldomain "github.com/shenfay/go-react-admin/internal/domain/crawl"
	"github.com/shenfay/go-react-admin/internal/app/emailverification"
	notificationapp "github.com/shenfay/go-react-admin/internal/app/notification"
	"github.com/shenfay/go-react-admin/internal/app/passwordreset"
	"github.com/shenfay/go-react-admin/internal/app/port"
	appsetting "github.com/shenfay/go-react-admin/internal/app/setting"
	"github.com/shenfay/go-react-admin/internal/app/shared/operationlog"
	"github.com/shenfay/go-react-admin/internal/app/tokenmanager"
	"github.com/shenfay/go-react-admin/internal/domain/notification"
	"github.com/shenfay/go-react-admin/internal/domain/operation"
	"github.com/shenfay/go-react-admin/internal/domain/rbac"
	"github.com/shenfay/go-react-admin/internal/domain/setting"
	"github.com/shenfay/go-react-admin/internal/domain/user"
	casbinenforcer "github.com/shenfay/go-react-admin/internal/infra/authorize"
	"github.com/shenfay/go-react-admin/internal/infra/bus"
	"github.com/shenfay/go-react-admin/internal/infra/config"
	"github.com/shenfay/go-react-admin/internal/infra/mail"
	"github.com/shenfay/go-react-admin/internal/infra/messaging"
	"github.com/shenfay/go-react-admin/internal/infra/repository"
	"github.com/shenfay/go-react-admin/internal/infra/ws"
	transhttp "github.com/shenfay/go-react-admin/internal/transport/http"
	"github.com/shenfay/go-react-admin/internal/transport/http/handlers"
	"github.com/shenfay/go-react-admin/pkg/health"
	pkglogger "github.com/shenfay/go-react-admin/pkg/logger"
	"github.com/shenfay/go-react-admin/pkg/metrics"
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
	// 1. 加载配置 + 安全检查
	env := loadEnv()
	cfg := mustLoadConfig(env)
	validateJWTSecret(cfg)

	// 2. 初始化日志和指标
	pkglogger.Init(cfg.Logger.Level)
	defer pkglogger.Sync()
	m := metrics.NewMetrics(prometheus.DefaultRegisterer)

	// 3. 初始化基础设施
	infra := initInfrastructure(cfg, m)
	defer infra.close()

	// 4. 初始化仓储
	repos := initRepositories(infra.db)

	// 5. 初始化应用服务
	svcs := initServices(cfg, infra, repos, m)

	// 6. 初始化传输层
	hdls := initHandlers(svcs, repos, infra)

	// 7. 启动 HTTP 服务器
	startServer(&startDeps{
		cfg:          cfg,
		metrics:      m,
		hdls:         hdls,
		tokenService: svcs.tokenService,
		enforcer:     infra.enforcer,
		db:           infra.db,
		redisClient:  infra.redisClient,
		hub:          infra.hub,
	})
}

// --- Provider 结构体 ---

type infraDeps struct {
	db             *gorm.DB
	redisClient    *redis.Client
	asynqClient    *asynq.Client
	bus            *bus.InProcessBus
	bridge         *messaging.DomainToIntegrationBridge
	enforcer       *casbinenforcer.Enforcer
	hub            *ws.Hub
	pushSubscriber *ws.RedisSubscriber
}

func (d *infraDeps) close() {
	if d.pushSubscriber != nil {
		d.pushSubscriber.Stop()
	}
	if sqlDB, _ := d.db.DB(); sqlDB != nil {
		sqlDB.Close()
	}
	d.redisClient.Close()
	d.asynqClient.Close()
}

// repoDeps 仓储层依赖
type repoDeps struct {
	userRepo    user.UserRepository
	roleRepo    rbac.RoleRepository
	menuRepo    rbac.MenuRepository
	operLogRepo operation.LogRepository
	settingRepo setting.Repository
	messageRepo notification.MessageRepository

	// 抓取模块（Go↔Python 集成）
	crawlSourceRepo crawldomain.SourceRepository
	crawlArticleRepo crawldomain.ArticleRepository
	crawlTaskRepo    crawldomain.TaskRunRepository
}

// svcDeps 应用服务层依赖
type svcDeps struct {
	tokenService     authentication.TokenService
	authService      *authentication.Service
	adminService     *admin.Service
	settingSvc       *appsetting.Service
	notificationSvc  *notificationapp.AppService
	emailVerifySvc   *emailverification.Service
	passwordResetSvc *passwordreset.Service

	// 抓取模块（Go↔Python 集成）
	crawlSvc *crawlapp.Service
}

// handlerDeps 传输层依赖
type handlerDeps struct {
	authHandler      *handlers.AuthHandler
	adminHandler     *handlers.AdminHandler
	operLogHdlr      *handlers.OperationLogHandler
	settingHdlr      *handlers.SettingHandler
	notificationHdlr *handlers.NotificationHandler
	wsHandler        *handlers.WSHandler
	crawlHdlr       *handlers.CrawlHandler
}

// --- Provider 函数 ---

// loadEnv 读取运行环境
func loadEnv() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	return env
}

// mustLoadConfig 加载配置，失败则终止
func mustLoadConfig(env string) *config.Config {
	cfg, err := config.Load(env)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return cfg
}

// validateJWTSecret 生产环境 JWT Secret 安全检查
func validateJWTSecret(cfg *config.Config) {
	if cfg.JWT.Secret == "your-jwt-secret-key-change-in-production" {
		log.Fatalf("FATAL: JWT secret is using the default value. Please set a secure JWT_SECRET in your configuration.")
	}
}

// initInfrastructure 初始化基础设施依赖
func initInfrastructure(cfg *config.Config, m *metrics.Metrics) *infraDeps {
	db, err := initDatabase(cfg.Database, m)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	redisClient := initRedis(cfg.Redis)
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.Asynq.Addr})

	bus := bus.NewInProcessBus()
	bridge := messaging.NewBridge(asynqClient)
	bridge.SubscribeTo(bus)

	enforcer, err := casbinenforcer.NewEnforcer(db)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin enforcer: %v", err)
	}
	pkglogger.Info("\u2713 Casbin enforcer initialized")

	deps := &infraDeps{
		db:          db,
		redisClient: redisClient,
		asynqClient: asynqClient,
		bus:         bus,
		bridge:      bridge,
		enforcer:    enforcer,
	}

	// WebSocket 实时推送（根据配置决定是否启用）
	if cfg.WebSocket.Enabled {
		hub := ws.NewHub()
		pushSubscriber := ws.NewRedisSubscriber(redisClient, hub)
		pushSubscriber.Start()
		deps.hub = hub
		deps.pushSubscriber = pushSubscriber
		pkglogger.Info("\u2713 WebSocket real-time push enabled")
	} else {
		pkglogger.Info("WebSocket real-time push is disabled")
	}

	return deps
}

// initRepositories 初始化仓储
func initRepositories(db *gorm.DB) *repoDeps {
	return &repoDeps{
		userRepo:    repository.NewUserRepository(db),
		roleRepo:    repository.NewRoleRepository(db),
		menuRepo:    repository.NewMenuRepository(db),
		operLogRepo: repository.NewOperationLogRepository(db),
		settingRepo: repository.NewSettingRepository(db),
		messageRepo: repository.NewMessageRepository(db),

		// 抓取模块（Go↔Python 集成）
		crawlSourceRepo: repository.NewCrawlSourceRepository(db),
		crawlArticleRepo: repository.NewCrawlArticleRepository(db),
		crawlTaskRepo:    repository.NewCrawlTaskRunRepository(db),
	}
}

// initServices 初始化应用服务
func initServices(cfg *config.Config, infra *infraDeps, repos *repoDeps, m *metrics.Metrics) *svcDeps {
	tokenService := authentication.NewTokenServiceImpl(authentication.TokenServiceConfig{
		RedisClient:   infra.redisClient,
		JWTSecret:     cfg.JWT.Secret,
		Issuer:        cfg.JWT.Issuer,
		AccessExpire:  cfg.JWT.AccessExpire,
		RefreshExpire: cfg.JWT.RefreshExpire,
	})

	// admin.Service 必须先创建：作为 PermissionQuerier 注入 authentication.Service，
	// 用于 Login 时查询用户权限（避免 auth 服务重复持有 roleRepo/menuRepo/enforcer）
	adminService := admin.NewService(repos.userRepo, repos.roleRepo, repos.menuRepo, infra.enforcer, infra.bus)

	authService := authentication.NewService(
		authentication.ServiceDeps{
			UserRepo:          repos.userRepo,
			TokenService:      tokenService,
			EventBus:          infra.bus,
			Metrics:           m,
			PermissionQuerier: adminService,
		},
		authentication.ServiceConfig{},
	)
	// 创建领域设置服务 → 应用层包装（注入操作日志）
	domainSettingSvc := setting.NewService(repos.settingRepo)
	settingRecorder := operationlog.NewOperationRecorder(infra.bus)
	settingSvc := appsetting.NewService(domainSettingSvc, settingRecorder)

	// 消息模块（WebSocket 启用时通过 Redis Pub/Sub 推送）
	notificationDomainSvc := notification.NewService(repos.messageRepo)
	var pushPublisher port.RealtimePusher
	if infra.hub != nil {
		pushPublisher = ws.NewRedisPublisher(infra.redisClient)
	}
	notificationSvc := notificationapp.NewAppService(notificationDomainSvc, pushPublisher)

	// 令牌管理器（共享邮箱验证 + 密码重置）
	tokenManager := tokenmanager.NewTokenManager(
		repository.NewVerificationTokenRepository(infra.db),
		repository.NewResetTokenRepository(infra.db),
		infra.redisClient,
		cfg.Auth.EmailVerificationExpire,
		cfg.Auth.PasswordResetExpire,
	)

	// 邮件发送器（开发环境使用 NoopSender）
	emailSender := mail.NewNoopSender()

	// 操作日志记录器
	verifRecorder := operationlog.NewOperationRecorder(infra.bus)
	resetRecorder := operationlog.NewOperationRecorder(infra.bus)

	// 邮箱验证服务
	emailVerifySvc := emailverification.NewService(
		tokenManager,
		repos.userRepo,
		emailSender,
		infra.bus,
		verifRecorder,
	)

	// 密码重置服务
	passwordResetSvc := passwordreset.NewService(
		tokenManager,
		repos.userRepo,
		emailSender,
		infra.bus,
		resetRecorder,
	)

	return &svcDeps{
		tokenService:     tokenService,
		authService:      authService,
		adminService:     adminService,
		settingSvc:       settingSvc,
		notificationSvc:  notificationSvc,
		emailVerifySvc:   emailVerifySvc,
		passwordResetSvc: passwordResetSvc,

		// 抓取模块（Go↔Python 集成）
		crawlSvc: crawlapp.NewService(
			repos.crawlSourceRepo,
			repos.crawlArticleRepo,
			repos.crawlTaskRepo,
			infra.redisClient,
			cfg.Scraper,
		),
	}
}

// initHandlers 初始化 HTTP 处理器
func initHandlers(svcs *svcDeps, repos *repoDeps, infra *infraDeps) *handlerDeps {
	hdls := &handlerDeps{
		authHandler:      handlers.NewAuthHandler(svcs.authService, svcs.emailVerifySvc, svcs.passwordResetSvc),
		adminHandler:     handlers.NewAdminHandler(svcs.adminService),
		operLogHdlr:      handlers.NewOperationLogHandler(repos.operLogRepo),
		settingHdlr:      handlers.NewSettingHandler(svcs.settingSvc),
		notificationHdlr: handlers.NewNotificationHandler(svcs.notificationSvc),
		crawlHdlr:       handlers.NewCrawlHandler(svcs.crawlSvc),
	}
	// WebSocket 处理器（仅在启用时创建）
	if infra.hub != nil {
		hdls.wsHandler = handlers.NewWSHandler(infra.hub, svcs.tokenService)
	}
	return hdls
}

// startDeps 服务器启动依赖
type startDeps struct {
	cfg          *config.Config
	metrics      *metrics.Metrics
	hdls         *handlerDeps
	tokenService authentication.TokenManager
	enforcer     *casbinenforcer.Enforcer
	db           *gorm.DB
	redisClient  *redis.Client
	hub          *ws.Hub
}

// startServer 创建并启动 HTTP 服务器（含优雅关闭）
func startServer(deps *startDeps) {
	cfg := deps.cfg
	m := deps.metrics
	hdls := deps.hdls

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	transhttp.Middlewares(engine, m, cfg.CORS)

	apiRouter := transhttp.NewRouter(&transhttp.RouterDeps{
		Engine:              engine,
		AuthHandler:         hdls.authHandler,
		AdminHandler:        hdls.adminHandler,
		OperationLogHandler: hdls.operLogHdlr,
		SettingHandler:      hdls.settingHdlr,
		NotificationHandler: hdls.notificationHdlr,
		WSHandler:           hdls.wsHandler,
		CrawlHandler:        hdls.crawlHdlr,
		TokenManager:        deps.tokenService,
		Enforcer:            deps.enforcer,
	})

	// 设置完整健康检查处理器（包含 DB/Redis 检查）
	healthHandler := health.NewHandler(deps.db, deps.redisClient, "1.0.0", cfg.Server.Mode)
	apiRouter.SetHealthHandler(healthHandler)

	apiRouter.Setup()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		pkglogger.Info("Starting HTTP server...", "port", cfg.Server.Port, "mode", cfg.Server.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 先关闭 WebSocket 连接，触发 HandleWS 退出
	if deps.hub != nil {
		deps.hub.Shutdown()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}

// initDatabase 初始化数据库连接
func initDatabase(cfg config.DatabaseConfig, m *metrics.Metrics) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 启动连接池监控
	if m != nil {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				stats := sqlDB.Stats()
				m.UpdateDBConnections(stats.OpenConnections, stats.MaxOpenConnections)
			}
		}()
	}

	log.Println("Database connection established and tables migrated")
	return db, nil
}

// initRedis 初始化 Redis 连接
func initRedis(cfg config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection established")
	return client
}
