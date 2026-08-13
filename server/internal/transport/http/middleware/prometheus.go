package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-react-admin/pkg/metrics"
	"github.com/shenfay/go-react-admin/pkg/utils"
)

// PrometheusMiddleware Prometheus 指标收集中间件
func PrometheusMiddleware(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 增加正在处理的请求数
		m.HTTPRequestsInFlight.Inc()
		defer m.HTTPRequestsInFlight.Dec()

		// 继续处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(start).Seconds()
		status := utils.ToString(c.Writer.Status())

		// 记录请求计数（带状态码标签）
		m.IncHTTPRequests(status)

		// 记录请求耗时
		path := c.FullPath()
		if path == "" {
			path = "unknown" // 未匹配路由（404）时 FullPath 为空
		}
		m.ObserveHTTPDuration(c.Request.Method, path, status, duration)
	}
}
