package user

import (
	"net/http"

	"github.com/shenfay/go-react-admin/pkg/errors"
)

// 用户域预定义错误
var (
	// ErrNotFound 用户不存在
	ErrNotFound = &errors.AppError{
		Code:       errors.ErrCodeUserNotFound,
		Message:    "用户不存在",
		HTTPStatus: http.StatusNotFound,
	}

	// ErrEmailAlreadyExists 邮箱已注册
	ErrEmailAlreadyExists = &errors.AppError{
		Code:       errors.ErrCodeUserEmailAlreadyExists,
		Message:    "该邮箱已被注册",
		HTTPStatus: http.StatusConflict,
	}
)

