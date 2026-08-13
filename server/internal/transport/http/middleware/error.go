package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-react-admin/internal/transport/http/response"
	"github.com/shenfay/go-react-admin/pkg/errors"
)

// ErrorHandling 统一错误处理中间件
// 自动处理通过 c.Error() 设置的错误，并注入 trace_id 和 timestamp
func ErrorHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理 c.Error() 设置的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleAppError(c, err)
			c.Abort()
		}
	}
}

// handleAppError 处理应用错误（自动注入 trace_id 和 timestamp）
func handleAppError(c *gin.Context, err error) {
	baseResponse := response.BaseResponse{
		TraceID:   response.GetTraceID(c),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// 尝试转换为 AppError
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.HTTPStatus, response.ErrorResponse{
			BaseResponse: baseResponse,
			Code:         appErr.Code,
			Message:      appErr.Message,
			Details:      appErr.Metadata,
		})
		return
	}

	// 未知错误，返回 500
	c.JSON(http.StatusInternalServerError, response.ErrorResponse{
		BaseResponse: baseResponse,
		Code:         "SYSTEM.INTERNAL_ERROR",
		Message:      "服务器内部错误",
	})
}

// RespondError 中间件统一错误响应（含 trace_id 和 timestamp）
func RespondError(c *gin.Context, httpStatus int, code string, message string) {
	c.JSON(httpStatus, response.ErrorResponse{
		BaseResponse: response.BaseResponse{
			TraceID:   response.GetTraceID(c),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Code:    code,
		Message: message,
	})
	c.Abort()
}
