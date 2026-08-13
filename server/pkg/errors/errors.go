package errors

import (
	"fmt"
)

// AppError 应用层标准错误结构（所有领域共享）
// 示例：
//
//	return &AppError{
//	    Code:       "AUTH.INVALID_CREDENTIALS",
//	    Message:    "Invalid email or password",
//	    HTTPStatus: http.StatusUnauthorized,
//	}
type AppError struct {
	Code       string      `json:"code"`              // 错误码：DOMAIN.ERROR_TYPE
	Message    string      `json:"message"`           // 用户友好消息
	HTTPStatus int         `json:"-"`                 // HTTP 状态码（不返回给客户端）
	Err        error       `json:"-"`                 // 内部错误（不返回给客户端）
	Metadata   interface{} `json:"details,omitempty"` // 额外元数据
}

// Error 实现 Go error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Err.Error())
	}
	return e.Message
}

// Unwrap 实现 errors.Unwrap 接口（支持 errors.Is 和 errors.As）
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithError 设置内部错误（链式调用）
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// WithMetadata 设置元数据（链式调用）
func (e *AppError) WithMetadata(metadata interface{}) *AppError {
	e.Metadata = metadata
	return e
}

// NewAppError 创建应用层错误
func NewAppError(code string, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

