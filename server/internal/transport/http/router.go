package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shenfay/go-react-admin/internal/app/authentication"
	"github.com/shenfay/go-react-admin/internal/infra/authorize"
	"github.com/shenfay/go-react-admin/internal/transport/http/handlers"
	"github.com/shenfay/go-react-admin/internal/transport/http/middleware"
	"github.com/shenfay/go-react-admin/pkg/health"
)

// RouterDeps 路由器依赖
type RouterDeps struct {
	Engine              *gin.Engine
	AuthHandler         *handlers.AuthHandler
	AdminHandler        *handlers.AdminHandler
	OperationLogHandler *handlers.OperationLogHandler
	SettingHandler      *handlers.SettingHandler
	NotificationHandler *handlers.NotificationHandler
	WSHandler           *handlers.WSHandler
	CrawlHandler        *handlers.CrawlHandler
	WorkerHandler       *handlers.WorkerHandler
	TokenManager        authentication.TokenManager
	Enforcer            *authorize.Enforcer
}

// Router 路由配置
type Router struct {
	engine              *gin.Engine
	authHandler         *handlers.AuthHandler
	adminHandler        *handlers.AdminHandler
	operationLogHandler *handlers.OperationLogHandler
	settingHandler      *handlers.SettingHandler
	notificationHandler *handlers.NotificationHandler
	wsHandler           *handlers.WSHandler
	crawlHandler        *handlers.CrawlHandler
	workerHandler       *handlers.WorkerHandler
	tokenManager        authentication.TokenManager
	enforcer            *authorize.Enforcer
	healthHandler       *health.Handler
}

// NewRouter 创建路由器
func NewRouter(deps *RouterDeps) *Router {
	return &Router{
		engine:              deps.Engine,
		authHandler:         deps.AuthHandler,
		adminHandler:        deps.AdminHandler,
		operationLogHandler: deps.OperationLogHandler,
		settingHandler:      deps.SettingHandler,
		notificationHandler: deps.NotificationHandler,
		wsHandler:           deps.WSHandler,
		crawlHandler:        deps.CrawlHandler,
		workerHandler:       deps.WorkerHandler,
		tokenManager:        deps.TokenManager,
		enforcer:            deps.Enforcer,
	}
}

// SetHealthHandler 设置健康检查处理器（可选，未设置时使用简单 fallback）
func (r *Router) SetHealthHandler(h *health.Handler) {
	r.healthHandler = h
}

// Setup 配置所有路由
func (r *Router) Setup() {
	// 注册全局中间件（顺序很重要！）
	r.engine.Use(middleware.TraceID())       // 1. 生成 trace_id
	r.engine.Use(middleware.RequestInfo())   // 2. 注入请求元数据（IP/UA/设备信息）
	r.engine.Use(middleware.ErrorHandling()) // 3. 错误处理

	// 健康检查（使用完整的 DB/Redis/Asynq 检查）
	if r.healthHandler != nil {
		r.healthHandler.RegisterRoutes(r.engine)
	} else {
		r.engine.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "healthy"})
		})
	}

	// Prometheus 指标端点
	r.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// WebSocket 端点（仅在启用时注册）
	if r.wsHandler != nil {
		r.engine.GET("/api/ws", r.wsHandler.ServeWS)
	}

	// API v1 路由组
	v1 := r.engine.Group("/api/v1")
	{
		r.setupAuthRoutes(v1)
		r.setupUserRoutes(v1)
		r.setupAdminRoutes(v1)
		r.setupOperationLogRoutes(v1)
		r.setupSettingRoutes(v1)
		r.setupNotificationRoutes(v1)
		r.setupCrawlRoutes(v1)
	}

	// 注册 Swagger UI 路由（开发环境）
	middleware.RegisterSwagger(r.engine, middleware.DefaultSwaggerConfig())
}

// setupAuthRoutes 配置认证相关路由
func (r *Router) setupAuthRoutes(v1 *gin.RouterGroup) {
	auth := v1.Group("/auth")
	{
		// 公开路由
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", middleware.LoginRateLimit(), r.authHandler.Login)
		auth.POST("/logout", r.authHandler.Logout)
		auth.POST("/refresh", r.authHandler.RefreshToken)

		// 邮箱验证（公开）
		auth.GET("/verify-email", r.authHandler.VerifyEmail)

		// 密码重置（公开）
		auth.POST("/forgot-password", r.authHandler.ForgotPassword)
		auth.POST("/reset-password", r.authHandler.ResetPassword)

		// 需要认证的路由
		authMiddleware := middleware.JWTAuthMiddleware(middleware.JWTAuthConfig{
			TokenService: r.tokenManager,
		})
		auth.GET("/me", authMiddleware, r.authHandler.GetCurrentUser)
		auth.GET("/devices", authMiddleware, r.authHandler.GetUserDevices)
		auth.DELETE("/devices/:token", authMiddleware, r.authHandler.RevokeDevice)
		auth.POST("/logout-all", authMiddleware, r.authHandler.LogoutAllDevices)
		auth.POST("/resend-verification", authMiddleware, r.authHandler.ResendVerification)
	}
}

// setupUserRoutes 配置用户相关路由
func (r *Router) setupUserRoutes(v1 *gin.RouterGroup) {
	users := v1.Group("/users")
	users.Use(middleware.JWTAuthMiddleware(middleware.JWTAuthConfig{
		TokenService: r.tokenManager,
	}))
	{
		users.GET("/:id", r.authHandler.GetUserByID)
	}
}

// setupAdminRoutes 配置管理员路由组
func (r *Router) setupAdminRoutes(v1 *gin.RouterGroup) {
	authMiddleware := middleware.JWTAuthMiddleware(middleware.JWTAuthConfig{
		TokenService: r.tokenManager,
	})
	permMiddleware := middleware.PermissionMiddleware(r.enforcer)

	adminGroup := v1.Group("/admin")
	adminGroup.Use(authMiddleware, permMiddleware)
	{
		// 用户管理
		adminGroup.GET("/users", r.adminHandler.ListUsers)
		adminGroup.POST("/users", r.adminHandler.CreateUser)
		adminGroup.PUT("/users/:id", r.adminHandler.UpdateUser)
		adminGroup.PATCH("/users/:id/status", r.adminHandler.ToggleUserStatus)

		// 角色管理
		adminGroup.GET("/roles", r.adminHandler.ListRoles)
		adminGroup.POST("/roles", r.adminHandler.CreateRole)
		adminGroup.PUT("/roles/:id", r.adminHandler.UpdateRole)
		adminGroup.DELETE("/roles/:id", r.adminHandler.DeleteRole)
		adminGroup.PATCH("/roles/:id/status", r.adminHandler.ToggleRoleStatus)

		// 权限管理
		adminGroup.GET("/roles/:id/permissions", r.adminHandler.GetRolePermissions)
		adminGroup.PUT("/roles/:id/permissions", r.adminHandler.UpdateRolePermissions)

		// 菜单管理
		adminGroup.GET("/menus", r.adminHandler.ListMenus)
		adminGroup.POST("/menus", r.adminHandler.CreateMenu)
		adminGroup.PUT("/menus/:id", r.adminHandler.UpdateMenu)
		adminGroup.DELETE("/menus/:id", r.adminHandler.DeleteMenu)
		adminGroup.PATCH("/menus/:id/status", r.adminHandler.ToggleMenuStatus)
		adminGroup.PUT("/menus/sort", r.adminHandler.UpdateMenuSort)
	}

	// 当前用户权限和菜单（放在 auth 组下，只需登录即可）
	auth := v1.Group("/auth")
	auth.Use(authMiddleware)
	{
		auth.GET("/permissions", r.adminHandler.GetCurrentUserPermissions)
		auth.GET("/menus", r.adminHandler.GetUserMenuTree)
	}
}

// setupOperationLogRoutes 配置操作日志路由（需要认证 + 管理员权限）
func (r *Router) setupOperationLogRoutes(v1 *gin.RouterGroup) {
	authMiddleware := middleware.JWTAuthMiddleware(middleware.JWTAuthConfig{
		TokenService: r.tokenManager,
	})
	permMiddleware := middleware.PermissionMiddleware(r.enforcer)

	operationLogs := v1.Group("/operation-logs")
	operationLogs.Use(authMiddleware, permMiddleware)
	{
		r.operationLogHandler.RegisterRoutes(operationLogs)
	}
}

// setupSettingRoutes 配置系统设置路由（需要认证 + setting:manage 权限）
func (r *Router) setupSettingRoutes(v1 *gin.RouterGroup) {
	authMiddleware := middleware.JWTAuthMiddleware(middleware.JWTAuthConfig{
		TokenService: r.tokenManager,
	})
	permMiddleware := middleware.PermissionMiddleware(r.enforcer)

	settings := v1.Group("/settings")
	settings.Use(authMiddleware, permMiddleware)
	{
		r.settingHandler.RegisterRoutes(settings)
	}
}

// setupNotificationRoutes 配置消息路由
func (r *Router) setupNotificationRoutes(v1 *gin.RouterGroup) {
	authMiddleware := middleware.JWTAuthMiddleware(middleware.JWTAuthConfig{
		TokenService: r.tokenManager,
	})
	permMiddleware := middleware.PermissionMiddleware(r.enforcer)

	// 用户消息接口（只需登录）
	messages := v1.Group("/messages")
	messages.Use(authMiddleware)
	{
		r.notificationHandler.RegisterRoutes(messages)
	}

	// 管理员消息管理（需要登录 + message:view 权限）
	adminMessages := v1.Group("/admin/messages")
	adminMessages.Use(authMiddleware, permMiddleware)
	{
		r.notificationHandler.RegisterAdminRoutes(adminMessages)
	}
}

// setupCrawlRoutes 配置抓取模块路由（仅管理路由，机器路由已由 Redis Stream 替代）
func (r *Router) setupCrawlRoutes(v1 *gin.RouterGroup) {
	authMiddleware := middleware.JWTAuthMiddleware(middleware.JWTAuthConfig{
		TokenService: r.tokenManager,
	})
	permMiddleware := middleware.PermissionMiddleware(r.enforcer)

	// 管理路由（管理员）
	admin := v1.Group("/crawl")
	admin.Use(authMiddleware, permMiddleware)
	{
		r.crawlHandler.RegisterAdminRoutes(admin)
		if r.workerHandler != nil {
			r.workerHandler.RegisterRoutes(admin)
		}
	}
}
