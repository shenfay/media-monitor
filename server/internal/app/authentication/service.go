package authentication

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/shenfay/go-react-admin/internal/app/shared/operationlog"
	"github.com/shenfay/go-react-admin/internal/domain/rbac"
	"github.com/shenfay/go-react-admin/internal/domain/shared/events"
	"github.com/shenfay/go-react-admin/internal/domain/user"
	authErr "github.com/shenfay/go-react-admin/pkg/errors/auth"
	userErr "github.com/shenfay/go-react-admin/pkg/errors/user"
	"github.com/shenfay/go-react-admin/pkg/logger"
	"github.com/shenfay/go-react-admin/pkg/metrics"
)

// JWTClaims JWT 自定义声明
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// TokenManager Token 核心操作接口（生成/验证/撤销）
type TokenManager interface {
	GenerateTokens(ctx context.Context, userID, email string) (*TokenPair, error)
	ValidateAccessToken(tokenString string) (*JWTClaims, error)
	ValidateRefreshTokenWithDevice(ctx context.Context, token string) (*DeviceInfo, error)
	RevokeToken(ctx context.Context, tokenID string) error
	RevokeDeviceByToken(ctx context.Context, token string) error
	RevokeAllDevices(ctx context.Context, userID string) error
}

// DeviceManager 设备会话管理接口（设备信息存储/查询/映射）
type DeviceManager interface {
	StoreDeviceInfo(ctx context.Context, token string, deviceInfo DeviceInfo) error
	GetUserDevices(ctx context.Context, userID string) ([]DeviceInfo, error)
	LinkAccessToDevice(ctx context.Context, accessToken, deviceTokenID string) error
	GetCurrentDeviceTokenID(ctx context.Context, accessToken string) (string, error)
}

// TokenService Token 服务完整接口（组合 TokenManager + DeviceManager）
type TokenService interface {
	TokenManager
	DeviceManager
}

// TokenPair 访问令牌和刷新令牌对
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
}

// DeviceInfo 设备会话信息
type DeviceInfo struct {
	TokenID    string `json:"token_id"`
	UserID     string `json:"user_id"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	DeviceType string `json:"device_type"`
	CreatedAt  string `json:"created_at"`
}

// PermissionQuerier 权限查询接口（由 admin.Service 实现，消除 Login 中的重复逻辑）
type PermissionQuerier interface {
	GetUserPermissions(ctx context.Context, userID string) (*rbac.UserPermission, error)
}

// ServiceConfig 认证服务配置
type ServiceConfig struct {
	MaxLoginAttempts int
}

// ServiceDeps 认证服务依赖
type ServiceDeps struct {
	UserRepo          user.UserRepository
	TokenService      TokenService
	EventBus          events.Bus
	Metrics           metrics.Recorder
	PermissionQuerier PermissionQuerier
}

// Service 认证应用服务
type Service struct {
	userRepo          user.UserRepository
	tokenService      TokenService
	metrics           metrics.Recorder
	maxAttempts       int
	permissionQuerier PermissionQuerier
	recorder          *operationlog.OperationRecorder
}

// NewService 创建认证服务实例
func NewService(deps ServiceDeps, cfg ServiceConfig) *Service {
	if cfg.MaxLoginAttempts <= 0 {
		cfg.MaxLoginAttempts = 5
	}
	return &Service{
		userRepo:          deps.UserRepo,
		tokenService:      deps.TokenService,
		metrics:           deps.Metrics,
		maxAttempts:       cfg.MaxLoginAttempts,
		permissionQuerier: deps.PermissionQuerier,
		recorder:          operationlog.NewOperationRecorder(deps.EventBus),
	}
}

// Register 创建用户账户并返回认证令牌
//
// 注册流程：
// 1. 验证邮箱唯一性
// 2. 创建用户实体（密码已加密）
// 3. 生成访问令牌和刷新令牌
// 4. 记录操作日志
//
// 参数：
//   - ctx: 请求上下文
//   - cmd: 注册命令（包含邮箱和密码）
//
// 返回：
//   - *ServiceAuthResponse: 用户数据和认证令牌
//   - error: 注册失败时返回错误（邮箱已存在、验证错误等）
func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (*ServiceAuthResponse, error) {
	// 1. 检查邮箱是否已存在
	if s.userRepo.ExistsByEmail(ctx, cmd.Email) {
		return nil, userErr.ErrEmailAlreadyExists
	}

	// 2. 创建用户实体
	u, err := user.NewUser(cmd.Email, "", cmd.Password)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	// 记录用户注册指标
	s.metrics.IncUserRegistration()

	// 3. 生成 Token
	tokens, err := s.tokenService.GenerateTokens(ctx, u.ID, u.Email)
	if err != nil {
		return nil, err
	}

	// 4. 记录操作日志
	s.recorder.Record(ctx, "USER.REGISTER", "USER", "SUCCESS",
		u.ID, u.Email,
		map[string]interface{}{"email": u.Email},
	)

	return &ServiceAuthResponse{
		User:         u,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}

// Login 处理用户登录
func (s *Service) Login(ctx context.Context, cmd LoginCommand) (*ServiceAuthResponse, error) {
	// 1. 查找用户
	u, err := s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, authErr.ErrInvalidCredentials
	}

	// 2. 检查账户是否被锁定
	if u.IsLocked() {
		return nil, authErr.ErrAccountLocked
	}

	// 3. 验证密码
	if !u.VerifyPassword(cmd.Password) {
		u.IncrementFailedAttempts(s.maxAttempts)
		s.userRepo.Update(ctx, u)

		// 记录认证失败指标
		if u.IsLocked() {
			s.metrics.IncAuthFailure("password", "account_locked")
		} else {
			s.metrics.IncAuthFailure("password", "invalid_credentials")
		}

		if u.IsLocked() {
			return nil, authErr.ErrAccountLocked
		}

		return nil, authErr.ErrInvalidCredentials
	}

	// 4. 重置失败次数，更新最后登录时间
	u.ResetFailedAttempts()
	u.UpdateLastLogin()
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	// 5. 生成 Token
	tokens, err := s.tokenService.GenerateTokens(ctx, u.ID, u.Email)
	if err != nil {
		return nil, err
	}

	// 记录认证成功指标
	s.metrics.IncAuthSuccess("password")

	// 6. 存储设备信息到 Redis
	if err := s.tokenService.StoreDeviceInfo(ctx, tokens.RefreshToken, DeviceInfo{
		TokenID:    tokens.RefreshToken,
		UserID:     u.ID,
		IP:         cmd.IP,
		UserAgent:  cmd.UserAgent,
		DeviceType: cmd.DeviceType,
	}); err != nil {
		// 设备信息存储失败不影响登录流程，仅记录警告
		// 日志已在 StoreDeviceInfo 内部处理
	}

	// 6.1 建立 access_token → device_token 映射，用于标识当前设备
	s.tokenService.LinkAccessToDevice(ctx, tokens.AccessToken, tokens.RefreshToken)

	// 7. 记录操作日志
	s.recorder.Record(ctx, "AUTH.LOGIN.SUCCESS", "AUTH", "SUCCESS",
		u.ID, u.Email,
		nil,
	)

	// 8. 查询用户权限（委托给 PermissionQuerier，消除重复逻辑）
	var permissions *rbac.UserPermission
	if s.permissionQuerier != nil {
		perm, err := s.permissionQuerier.GetUserPermissions(ctx, u.ID)
		if err != nil {
			logger.Warn("Failed to query user permissions on login",
				"user_id", u.ID,
				"error", err,
			)
		} else {
			permissions = perm
		}
	}

	return &ServiceAuthResponse{
		User:         u,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		Permissions:  permissions,
	}, nil
}

// Logout 处理用户退出
func (s *Service) Logout(ctx context.Context, cmd LogoutCommand) error {
	// 1. 撤销该用户所有设备的登录状态（Refresh Token + 设备信息）
	if err := s.tokenService.RevokeAllDevices(ctx, cmd.UserID); err != nil {
		logger.Warn("Failed to revoke devices on logout",
			"user_id", cmd.UserID,
			"error", err,
		)
	}

	// 2. 记录操作日志（直接使用 userID，省去额外 DB 查询）
	s.recorder.Record(ctx, "AUTH.LOGOUT", "AUTH", "SUCCESS", cmd.UserID, "", nil)

	return nil
}

// RefreshToken 刷新 Access Token
func (s *Service) RefreshToken(ctx context.Context, cmd RefreshTokenCommand) (*ServiceAuthResponse, error) {
	// 1. 验证并解析 Refresh Token
	deviceInfo, err := s.tokenService.ValidateRefreshTokenWithDevice(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 2. 查找用户
	u, err := s.userRepo.FindByID(ctx, deviceInfo.UserID)
	if err != nil {
		return nil, userErr.ErrNotFound
	}

	// 3. 撤销旧的 Refresh Token
	if err := s.tokenService.RevokeDeviceByToken(ctx, cmd.RefreshToken); err != nil {
		return nil, err
	}

	// 4. 生成新的 Token 对
	tokens, err := s.tokenService.GenerateTokens(ctx, u.ID, u.Email)
	if err != nil {
		return nil, err
	}

	// 5. 存储新设备信息
	if err := s.tokenService.StoreDeviceInfo(ctx, tokens.RefreshToken, DeviceInfo{
		TokenID:    tokens.RefreshToken,
		UserID:     u.ID,
		IP:         deviceInfo.IP,
		UserAgent:  deviceInfo.UserAgent,
		DeviceType: deviceInfo.DeviceType,
	}); err != nil {
		return nil, err
	}

	// 5.1 建立 access_token → device_token 映射
	s.tokenService.LinkAccessToDevice(ctx, tokens.AccessToken, tokens.RefreshToken)

	// 6. 更新最后登录时间
	u.UpdateLastLogin()
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	// 7. 记录操作日志（已脱敏，不记录 token 明文）
	s.recorder.Record(ctx, "AUTH.TOKEN.REFRESHED", "AUTH", "SUCCESS",
		u.ID, u.Email,
		nil,
	)

	return &ServiceAuthResponse{
		User:         u,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}, nil
}

// GetUserByID 根据 ID 获取用户
func (s *Service) GetUserByID(ctx context.Context, userID string) (*user.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// ListUserDevices 获取用户所有登录设备
func (s *Service) ListUserDevices(ctx context.Context, userID string) ([]DeviceInfo, error) {
	return s.tokenService.GetUserDevices(ctx, userID)
}

// GetCurrentDeviceTokenID 获取当前请求对应的设备令牌 ID
func (s *Service) GetCurrentDeviceTokenID(ctx context.Context, accessToken string) (string, error) {
	return s.tokenService.GetCurrentDeviceTokenID(ctx, accessToken)
}

// RevokeDevice 撤销指定设备（校验设备归属）
func (s *Service) RevokeDevice(ctx context.Context, userID, token string) error {
	deviceInfo, err := s.tokenService.ValidateRefreshTokenWithDevice(ctx, token)
	if err != nil {
		return err
	}
	if deviceInfo.UserID != userID {
		return authErr.ErrForbidden
	}
	return s.tokenService.RevokeDeviceByToken(ctx, token)
}

// RevokeAllDevices 撤销用户所有设备
func (s *Service) RevokeAllDevices(ctx context.Context, userID string) error {
	return s.tokenService.RevokeAllDevices(ctx, userID)
}
