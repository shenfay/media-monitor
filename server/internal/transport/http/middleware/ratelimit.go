package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	visitors    map[string]*visitor
	mu          sync.Mutex
	rate        rate.Limit
	burst       int
	lastCleanup time.Time
}

// visitor 访问者信息
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(rate rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		visitors:    make(map[string]*visitor),
		rate:        rate,
		burst:       burst,
		lastCleanup: time.Now(),
	}
}

// allow 检查是否允许请求
func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		v = &visitor{limiter: limiter, lastSeen: time.Now()}
		rl.visitors[ip] = v
	} else {
		v.lastSeen = time.Now()
	}

	// 每分钟清理一次过期用户
	if time.Since(rl.lastCleanup) > time.Minute {
		rl.cleanup()
		rl.lastCleanup = time.Now()
	}

	return v.limiter.Allow()
}

// cleanup 清理长时间未访问的用户
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > 3*time.Minute {
			delete(rl.visitors, ip)
		}
	}
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.allow(ip) {
			RespondError(c, http.StatusTooManyRequests, "SYSTEM.TOO_MANY_REQUESTS", "请求过于频繁，请稍后重试")
			return
		}

		c.Next()
	}
}

// 包级单例限流器，避免每次请求创建新实例导致限流失效
var (
	loginLimiter   = NewRateLimiter(rate.Every(time.Minute/5), 10)   // 5 次/分钟，突发 10 次
	generalLimiter = NewRateLimiter(rate.Every(time.Minute/60), 100) // 60 次/分钟，突发 100 次
)

// LoginRateLimit 登录接口专用速率限制（更严格）
func LoginRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(loginLimiter)
}

// GeneralRateLimit 通用速率限制
func GeneralRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(generalLimiter)
}
