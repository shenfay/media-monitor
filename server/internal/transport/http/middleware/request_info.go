package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-react-admin/pkg/utils"
)

// RequestInfo 请求元数据注入中间件
// 从 HTTP 请求中提取 IP、User-Agent、设备/浏览器/操作系统信息，
// 注入到标准 context 中，供 Service 层操作日志记录使用。
// 必须在 JWTAuth 中间件之后使用（依赖 context 中已有的值）。
func RequestInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		uaInfo := utils.ParseUserAgent(userAgent)

		ctx := utils.WithRequestInfo(
			c.Request.Context(),
			ip,
			userAgent,
			uaInfo.Device,
			uaInfo.Browser,
			uaInfo.OS,
		)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
