package errors

// ==========================================
// 错误码定义（按领域分组）
// 格式：{DOMAIN}.{CATEGORY}.{ERROR}
// ==========================================

// 系统级错误码 (SYSTEM.*)
const (
	ErrCodeSystemInternal        = "SYSTEM.INTERNAL_ERROR"
	ErrCodeSystemInvalidRequest  = "SYSTEM.INVALID_REQUEST"
	ErrCodeSystemUnauthorized    = "SYSTEM.UNAUTHORIZED"
	ErrCodeSystemForbidden       = "SYSTEM.FORBIDDEN"
	ErrCodeSystemNotFound        = "SYSTEM.NOT_FOUND"
	ErrCodeSystemConflict        = "SYSTEM.CONFLICT"
	ErrCodeSystemTooManyRequests = "SYSTEM.TOO_MANY_REQUESTS"
)

// 认证域错误码 (AUTH.*)
const (
	ErrCodeAuthInvalidCredentials = "AUTH.INVALID_CREDENTIALS"
	ErrCodeAuthInvalidToken       = "AUTH.INVALID_TOKEN"
	ErrCodeAuthTokenExpired       = "AUTH.TOKEN_EXPIRED"
	ErrCodeAuthTokenRevoked       = "AUTH.TOKEN_REVOKED"
	ErrCodeAuthAccountLocked      = "AUTH.ACCOUNT_LOCKED"
)

// 用户域错误码 (USER.*)
const (
	ErrCodeUserNotFound           = "USER.NOT_FOUND"
	ErrCodeUserEmailAlreadyExists = "USER.EMAIL_ALREADY_EXISTS"
)

// 通知域错误码 (NOTIFICATION.*)
const (
	ErrCodeMessageNotFound     = "NOTIFICATION.MESSAGE_NOT_FOUND"
	ErrCodeMessageAccessDenied = "NOTIFICATION.MESSAGE_ACCESS_DENIED"
)

// 令牌域错误码 (TOKEN.*)
const (
	ErrCodeTokenAlreadyUsed = "TOKEN.ALREADY_USED"
	ErrCodeTokenExpired     = "TOKEN.EXPIRED"
	ErrCodeTokenInvalid     = "TOKEN.INVALID"
)

// 菜单域错误码 (MENU.*)
const (
	ErrCodeMenuKeyAlreadyExists = "MENU.KEY_ALREADY_EXISTS"
)
