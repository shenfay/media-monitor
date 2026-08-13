package validation

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shenfay/go-react-admin/pkg/errors"
)

// 校验域预定义错误
var (
	// ErrFieldRequired 字段必填
	ErrFieldRequired = &errors.AppError{
		Code:       "VALIDATION.FIELD_REQUIRED",
		Message:    "必填字段不能为空",
		HTTPStatus: http.StatusBadRequest,
	}

	// ErrFieldInvalid 字段格式无效
	ErrFieldInvalid = &errors.AppError{
		Code:       "VALIDATION.FIELD_INVALID",
		Message:    "字段格式不正确",
		HTTPStatus: http.StatusBadRequest,
	}

	// ErrFieldTooShort 字段长度太短
	ErrFieldTooShort = &errors.AppError{
		Code:       "VALIDATION.FIELD_TOO_SHORT",
		Message:    "字段长度太短",
		HTTPStatus: http.StatusBadRequest,
	}

	// ErrFieldTooLong 字段长度太长
	ErrFieldTooLong = &errors.AppError{
		Code:       "VALIDATION.FIELD_TOO_LONG",
		Message:    "字段长度太长",
		HTTPStatus: http.StatusBadRequest,
	}
)


// FromGinError 从 Gin 绑定错误创建验证错误
func FromGinError(err error) *errors.AppError {
	if err == nil {
		return nil
	}

	// 检查是否是验证错误
	if ve, ok := err.(validator.ValidationErrors); ok {
		messages := make([]string, 0, len(ve))
		for _, fe := range ve {
			messages = append(messages, formatFieldError(fe))
		}
		return &errors.AppError{
			Code:       "VALIDATION.INVALID_REQUEST",
			Message:    strings.Join(messages, "; "),
			HTTPStatus: http.StatusBadRequest,
		}
	}

	// 其他错误
	return &errors.AppError{
		Code:       "VALIDATION.INVALID_REQUEST",
		Message:    err.Error(),
		HTTPStatus: http.StatusBadRequest,
	}
}

// formatFieldError 格式化字段错误
func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s 为必填字段", fe.Field())
	case "email":
		return fmt.Sprintf("%s 必须是有效的邮箱地址", fe.Field())
	case "min":
		return fmt.Sprintf("%s 长度不能少于 %s 个字符", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s 长度不能超过 %s 个字符", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s 格式不正确", fe.Field())
	}
}
