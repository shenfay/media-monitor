package authentication

import (
	"time"

	"github.com/shenfay/go-react-admin/internal/domain/rbac"
	"github.com/shenfay/go-react-admin/internal/domain/user"
)

// UserResponse 用户响应 DTO
type UserResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Name          string     `json:"name"`
	EmailVerified bool       `json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// AuthResponse 认证响应 DTO
type AuthResponse struct {
	User         *UserResponse        `json:"user"`
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
	ExpiresIn    int64                `json:"expires_in"`
	Permissions  *rbac.UserPermission `json:"permissions,omitempty"`
}

// ToUserResponse 将领域实体转换为用户响应 DTO
func ToUserResponse(u *user.User) *UserResponse {
	if u == nil {
		return nil
	}

	return &UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		Name:          u.Name,
		EmailVerified: u.EmailVerified,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
	}
}

// ToAuthResponse 将服务层响应转换为认证响应 DTO
func ToAuthResponse(resp *ServiceAuthResponse) *AuthResponse {
	if resp == nil {
		return nil
	}

	return &AuthResponse{
		User:         ToUserResponse(resp.User),
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    int64(resp.ExpiresIn / time.Second),
		Permissions:  resp.Permissions,
	}
}
