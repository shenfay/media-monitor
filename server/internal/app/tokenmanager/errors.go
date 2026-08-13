package tokenmanager

import (
	"net/http"

	"github.com/shenfay/go-react-admin/pkg/errors"
)

var (
	ErrTokenAlreadyUsed = &errors.AppError{
		Code:       errors.ErrCodeTokenAlreadyUsed,
		Message:    "令牌已被使用",
		HTTPStatus: http.StatusGone,
	}
	ErrTokenExpired = &errors.AppError{
		Code:       errors.ErrCodeTokenExpired,
		Message:    "令牌已过期",
		HTTPStatus: http.StatusGone,
	}
	ErrTokenUserMismatch = &errors.AppError{
		Code:       errors.ErrCodeTokenInvalid,
		Message:    "令牌与用户不匹配",
		HTTPStatus: http.StatusBadRequest,
	}
)
