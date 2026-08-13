package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TraceIDKey Gin Context 中 trace_id 的键
const TraceIDKey = "trace_id"

// BaseResponse 响应基础结构
type BaseResponse struct {
	TraceID   string `json:"trace_id"`  // 链路追踪 ID（用于日志关联和分布式追踪）
	Timestamp string `json:"timestamp"` // 响应时间（RFC3339 格式）
}

// SuccessResponse 成功响应结构
type SuccessResponse struct {
	BaseResponse             // 嵌套，JSON 自动扁平化
	Code         string      `json:"code"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	BaseResponse             // 嵌套
	Code         string      `json:"code"`
	Message      string      `json:"message"`
	Details      interface{} `json:"details,omitempty"`
}

// GetTraceID 从 Gin Context 获取 trace_id
func GetTraceID(c *gin.Context) string {
	if id, exists := c.Get(TraceIDKey); exists {
		return id.(string)
	}
	return ""
}

// newBaseResponse 创建基础响应结构
func newBaseResponse(c *gin.Context) BaseResponse {
	return BaseResponse{
		TraceID:   GetTraceID(c),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// Success 返回成功响应（200 OK）
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		BaseResponse: newBaseResponse(c),
		Code:         "SUCCESS",
		Message:      "Request successful",
		Data:         data,
	})
}

// Created 返回创建成功响应（201 Created）
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		BaseResponse: newBaseResponse(c),
		Code:         "CREATED",
		Message:      "Resource created successfully",
		Data:         data,
	})
}

// NoContent 返回无内容响应（204 No Content）
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error 返回错误响应（交由错误处理中间件处理）
func Error(c *gin.Context, err error) {
	c.Error(err)
}
