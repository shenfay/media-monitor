package auth

import (
	"net/http"

	"github.com/shenfay/go-react-admin/pkg/errors"
)

// 认证域预定义错误
var (
	// ErrInvalidCredentials 凭证无效（邮箱或密码错误）
	ErrInvalidCredentials = &errors.AppError{
		Code:       errors.ErrCodeAuthInvalidCredentials,
		Message:    "邮箱或密码错误",
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrInvalidToken 无效 Token
	ErrInvalidToken = &errors.AppError{
		Code:       errors.ErrCodeAuthInvalidToken,
		Message:    "无效的认证令牌",
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrTokenExpired Token 已过期
	ErrTokenExpired = &errors.AppError{
		Code:       errors.ErrCodeAuthTokenExpired,
		Message:    "登录已过期，请重新登录",
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrTokenRevoked Token 已被撤销
	ErrTokenRevoked = &errors.AppError{
		Code:       errors.ErrCodeAuthTokenRevoked,
		Message:    "令牌已被撤销",
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrAccountLocked 账户已锁定
	ErrAccountLocked = &errors.AppError{
		Code:       errors.ErrCodeAuthAccountLocked,
		Message:    "账户已锁定，请稍后重试",
		HTTPStatus: http.StatusLocked,
	}

	// ErrForbidden 操作无权限
	ErrForbidden = &errors.AppError{
		Code:       errors.ErrCodeSystemForbidden,
		Message:    "禁止操作：只能管理自己的设备",
		HTTPStatus: http.StatusForbidden,
	}
)

