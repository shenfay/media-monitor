package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-react-admin/internal/infra/authorize"
)

// RoutePermissions 路由 → 所需权限映射
var RoutePermissions = map[string]string{
	"GET:/api/v1/admin/users":              "user:manage",
	"POST:/api/v1/admin/users":             "user:manage",
	"PUT:/api/v1/admin/users/:id":          "user:manage",
	"PATCH:/api/v1/admin/users/:id/status": "user:manage",

	"GET:/api/v1/admin/roles":                 "permission:manage",
	"POST:/api/v1/admin/roles":                "permission:manage",
	"PUT:/api/v1/admin/roles/:id":             "permission:manage",
	"DELETE:/api/v1/admin/roles/:id":          "permission:manage",
	"PATCH:/api/v1/admin/roles/:id/status":    "permission:manage",
	"GET:/api/v1/admin/roles/:id/permissions": "permission:manage",
	"PUT:/api/v1/admin/roles/:id/permissions": "permission:manage",

	"GET:/api/v1/admin/menus":              "menu:manage",
	"POST:/api/v1/admin/menus":             "menu:manage",
	"PUT:/api/v1/admin/menus/:id":          "menu:manage",
	"DELETE:/api/v1/admin/menus/:id":       "menu:manage",
	"PATCH:/api/v1/admin/menus/:id/status": "menu:manage",
	"PUT:/api/v1/admin/menus/sort":         "menu:manage",

	// 抓取模块（管理后台）
	"GET:/api/v1/admin/crawl/sources":              "source:view",
	"POST:/api/v1/admin/crawl/sources":             "source:manage",
	"PUT:/api/v1/admin/crawl/sources/:id":          "source:manage",
	"DELETE:/api/v1/admin/crawl/sources/:id":       "source:manage",
	"POST:/api/v1/admin/crawl/sources/:id/run":     "source:manage",
	"POST:/api/v1/admin/crawl/tasks":               "task:manage",
	"GET:/api/v1/admin/crawl/tasks":                "task:view",
	"GET:/api/v1/admin/crawl/tasks/:id":            "task:view",
	"GET:/api/v1/admin/crawl/articles":             "article:view",
	"GET:/api/v1/admin/crawl/articles/:id":         "article:view",
	"GET:/api/v1/admin/crawl/dashboard/stats":      "crawl:view",
	"GET:/api/v1/admin/crawl/workers":              "worker:view",
	"GET:/api/v1/admin/crawl/workers/adapters":     "worker:view",
	"POST:/api/v1/admin/crawl/workers/:id/pause":   "worker:manage",
	"POST:/api/v1/admin/crawl/workers/:id/resume":  "worker:manage",
	"POST:/api/v1/admin/crawl/workers/:id/shutdown": "worker:manage",
}

// PermissionMiddleware 基于 Casbin 的权限检查中间件
// 根据请求方法+路径匹配所需权限，通过 Casbin Enforcer 鉴权
func PermissionMiddleware(enforcer *authorize.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			RespondError(c, http.StatusUnauthorized, "SYSTEM.UNAUTHORIZED", "缺少用户身份信息")
			return
		}

		// 构建路由 key：METHOD:/full/path
		routeKey := c.Request.Method + ":" + c.FullPath()

		// 查找所需权限
		requiredPerm, ok := RoutePermissions[routeKey]
		if !ok {
			// 未配置权限的路由，默认放行（可能是公开路由或内部路由）
			c.Next()
			return
		}

		// Casbin 鉴权
		allowed, err := enforcer.Enforce(userID, requiredPerm)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, "SYSTEM.INTERNAL_ERROR", "权限检查失败")
			return
		}

		if !allowed {
			RespondError(c, http.StatusForbidden, "SYSTEM.FORBIDDEN", "权限不足，无法访问该资源")
			return
		}

		c.Next()
	}
}
