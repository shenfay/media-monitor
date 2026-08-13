package authentication

import (
	"time"

	"github.com/shenfay/go-react-admin/internal/domain/rbac"
	"github.com/shenfay/go-react-admin/internal/domain/user"
)

// RegisterCommand 注册命令
type RegisterCommand struct {
	Email    string
	Password string
}

// LoginCommand 登录命令
type LoginCommand struct {
	Email      string
	Password   string
	IP         string
	UserAgent  string
	DeviceType string // 设备类型：desktop, mobile, tablet
}

// RefreshTokenCommand 刷新 Token 命令
type RefreshTokenCommand struct {
	RefreshToken string
}

// LogoutCommand 退出登录命令
type LogoutCommand struct {
	UserID string
}

// ServiceAuthResponse 服务层认证响应（内部使用）
type ServiceAuthResponse struct {
	User         *user.User
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
	Permissions  *rbac.UserPermission
}
