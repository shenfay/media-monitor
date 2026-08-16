package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-react-admin/internal/app/authentication"
	"github.com/shenfay/go-react-admin/internal/app/emailverification"
	"github.com/shenfay/go-react-admin/internal/app/passwordreset"
	"github.com/shenfay/go-react-admin/internal/transport/http/response"
	"github.com/shenfay/go-react-admin/pkg/errors"
	userErr "github.com/shenfay/go-react-admin/pkg/errors/user"
	validationErr "github.com/shenfay/go-react-admin/pkg/errors/validation"
)

// AuthHandler 认证 HTTP 处理器
type AuthHandler struct {
	service          *authentication.Service
	emailVerifySvc   *emailverification.Service
	passwordResetSvc *passwordreset.Service
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(service *authentication.Service, emailVerifySvc *emailverification.Service, passwordResetSvc *passwordreset.Service) *AuthHandler {
	return &AuthHandler{
		service:          service,
		emailVerifySvc:   emailVerifySvc,
		passwordResetSvc: passwordResetSvc,
	}
}

// RegisterPublicRoutes 注册公开认证路由（无需 JWT）
// 注意：/login 路由由 router.go 注册以附加限流中间件
func (h *AuthHandler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.Register)
	rg.POST("/logout", h.Logout)
	rg.POST("/refresh", h.RefreshToken)
	rg.GET("/verify-email", h.VerifyEmail)
	rg.POST("/forgot-password", h.ForgotPassword)
	rg.POST("/reset-password", h.ResetPassword)
}

// RegisterAuthRoutes 注册需 JWT 认证的认证路由
func (h *AuthHandler) RegisterAuthRoutes(rg *gin.RouterGroup) {
	rg.GET("/me", h.GetCurrentUser)
	rg.GET("/devices", h.GetUserDevices)
	rg.DELETE("/devices/:token", h.RevokeDevice)
	rg.POST("/logout-all", h.LogoutAllDevices)
	rg.POST("/resend-verification", h.ResendVerification)
}

// RegisterUserRoutes 注册用户查询路由（需 JWT 认证）
func (h *AuthHandler) RegisterUserRoutes(rg *gin.RouterGroup) {
	rg.GET("/:id", h.GetUserByID)
}

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`   // 用户邮箱
	Password string `json:"password" binding:"required,min=8,max=72"` // 用户密码（8-72字符）
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"` // 用户邮箱
	Password string `json:"password" binding:"required"`    // 用户密码
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"` // 刷新令牌
}

// Register 处理用户注册
//
// 创建新用户账户并返回认证令牌。
// 邮箱必须在系统中唯一。
//
// @Summary Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "用户注册数据"
// @Success 201 {object} response.SuccessResponse{data=authentication.AuthResponse} "注册成功"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 409 {object} response.ErrorResponse "邮箱已存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := authentication.RegisterCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.service.Register(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 触发发送验证邮件（异步非阻塞）
	_ = h.emailVerifySvc.SendVerificationEmail(c.Request.Context(), resp.User.ID, resp.User.Email)

	response.Created(c, authentication.ToAuthResponse(resp))
}

// Login 处理用户登录
//
// 验证用户凭据并返回访问/刷新令牌。
// 跟踪登录失败次数，失败过多时锁定账户。
//
// @Summary 用户登录并返回令牌
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} response.SuccessResponse{data=authentication.AuthResponse} "登录成功"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "账号或密码错误"
// @Failure 423 {object} response.ErrorResponse "账户已锁定"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := authentication.LoginCommand{
		Email:      req.Email,
		Password:   req.Password,
		IP:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		DeviceType: detectDeviceType(c.Request.UserAgent()),
	}

	resp, err := h.service.Login(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, authentication.ToAuthResponse(resp))
}

// Logout 处理用户退出
// @Summary User logout
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetString("user_id")

	cmd := authentication.LogoutCommand{
		UserID: userID,
	}

	if err := h.service.Logout(c.Request.Context(), cmd); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Logged out successfully"})
}

// RefreshToken 刷新 Access Token
// @Summary Refresh access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} authentication.AuthResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse "Invalid or expired token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := authentication.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, authentication.ToAuthResponse(resp))
}

// GetUserByID 根据 ID 获取用户信息
// @Summary Get user by ID
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} authentication.UserResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Router /users/{id} [get]
func (h *AuthHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.Error(c, errors.NewAppError(
			errors.ErrCodeSystemInvalidRequest,
			"用户 ID 不能为空",
			http.StatusBadRequest,
		))
		return
	}

	u, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, userErr.ErrNotFound)
		return
	}

	response.Success(c, authentication.ToUserResponse(u))
}

// GetCurrentUser 获取当前登录用户信息
// @Summary Get current user information
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} authentication.UserResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, errors.NewAppError(
			errors.ErrCodeSystemUnauthorized,
			"缺少用户身份信息",
			http.StatusUnauthorized,
		))
		return
	}

	u, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, userErr.ErrNotFound)
		return
	}

	response.Success(c, authentication.ToUserResponse(u))
}

// maskIP 脱敏 IP 地址（支持 IPv4 和 IPv6）
func maskIP(ip string) string {
	if ip == "" {
		return ""
	}
	// IPv6 包含冒号
	if strings.Contains(ip, ":") {
		parts := strings.Split(ip, ":")
		if len(parts) >= 3 {
			return parts[0] + ":" + parts[1] + ":***"
		}
		return ip[:len(ip)/2] + "***"
	}
	// IPv4
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		return parts[0] + "." + parts[1] + ".***"
	}
	return ip[:len(ip)/2] + "***"
}

// DeviceResponse 设备响应
type DeviceResponse struct {
	TokenID    string `json:"token_id"`
	DeviceType string `json:"device_type"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	CreatedAt  string `json:"created_at"`
	IsCurrent  bool   `json:"is_current"`
}

// DevicesResponse 设备列表响应
type DevicesResponse struct {
	Devices []DeviceResponse `json:"devices"`
}

// GetUserDevices 获取当前用户的所有登录设备
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} DevicesResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /auth/devices [get]
func (h *AuthHandler) GetUserDevices(c *gin.Context) {
	userID := c.GetString("user_id")

	devices, err := h.service.ListUserDevices(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, errors.NewAppError(
			errors.ErrCodeSystemInternal,
			"获取设备列表失败",
			http.StatusInternalServerError,
		))
		return
	}

	// 通过 access_token → device_token 映射标识当前设备
	currentDeviceTokenID := ""
	authHeader := c.GetHeader("Authorization")
	if parts := strings.SplitN(authHeader, " ", 2); len(parts) == 2 && parts[0] == "Bearer" {
		if id, err := h.service.GetCurrentDeviceTokenID(c.Request.Context(), parts[1]); err == nil {
			currentDeviceTokenID = id
		}
	}

	var deviceResponses []DeviceResponse
	for _, device := range devices {
		deviceResponses = append(deviceResponses, DeviceResponse{
			TokenID:    device.TokenID,
			DeviceType: device.DeviceType,
			IP:         maskIP(device.IP),
			UserAgent:  device.UserAgent,
			CreatedAt:  device.CreatedAt,
			IsCurrent:  device.TokenID == currentDeviceTokenID,
		})
	}

	response.Success(c, DevicesResponse{
		Devices: deviceResponses,
	})
}

// RevokeDevice 踢出指定设备
// @Summary Revoke a specific device
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Param token path string true "Device token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /auth/devices/{token} [delete]
func (h *AuthHandler) RevokeDevice(c *gin.Context) {
	userID := c.GetString("user_id")
	token := c.Param("token")

	if token == "" {
		response.Error(c, errors.NewAppError(
			errors.ErrCodeSystemInvalidRequest,
			"设备令牌不能为空",
			http.StatusBadRequest,
		))
		return
	}

	if err := h.service.RevokeDevice(c.Request.Context(), userID, token); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Device revoked successfully"})
}

// LogoutAllDevices 退出所有设备
// @Summary Logout from all devices
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /auth/logout-all [post]
func (h *AuthHandler) LogoutAllDevices(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.service.RevokeAllDevices(c.Request.Context(), userID); err != nil {
		response.Error(c, errors.NewAppError(
			errors.ErrCodeSystemInternal,
			"退出所有设备失败",
			http.StatusInternalServerError,
		))
		return
	}

	response.Success(c, gin.H{"message": "Logged out from all devices successfully"})
}

// VerifyEmailRequest 邮箱验证请求
// GET /auth/verify-email?token=xxx&user_id=xxx
type VerifyEmailRequest struct {
	Token  string `form:"token" binding:"required"`
	UserID string `form:"user_id" binding:"required"`
}

// VerifyEmail 验证邮箱
// @Summary Verify email address
// @Tags Authentication
// @Produce json
// @Param token query string true "验证令牌"
// @Param user_id query string true "用户 ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 410 {object} response.ErrorResponse "令牌已过期或已使用"
// @Router /auth/verify-email [get]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	if err := h.emailVerifySvc.VerifyEmail(c.Request.Context(), req.Token, req.UserID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "邮箱验证成功"})
}

// ResendVerification 重新发送验证邮件
// @Summary Resend verification email
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /auth/resend-verification [post]
func (h *AuthHandler) ResendVerification(c *gin.Context) {
	userID := c.GetString("user_id")
	email := c.GetString("email")

	if userID == "" || email == "" {
		response.Error(c, errors.NewAppError(
			errors.ErrCodeSystemUnauthorized,
			"缺少用户身份信息",
			http.StatusUnauthorized,
		))
		return
	}

	if err := h.emailVerifySvc.SendVerificationEmail(c.Request.Context(), userID, email); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "验证邮件已发送"})
}

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgotPassword 请求密码重置
// @Summary Request password reset
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "注册邮箱"
// @Success 200 {object} response.SuccessResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	// 始终返回成功（防枚举）
	_ = h.passwordResetSvc.ForgotPassword(c.Request.Context(), req.Email)
	response.Success(c, gin.H{"message": "如果该邮箱已注册，重置密码邮件已发送"})
}

// ResetPasswordRequest 密码重置请求
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72"`
}

// ResetPassword 执行密码重置
// @Summary Reset password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "重置令牌和新密码"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 410 {object} response.ErrorResponse "令牌已过期或已使用"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	if err := h.passwordResetSvc.ResetPassword(c.Request.Context(), req.Token, req.UserID, req.NewPassword); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "密码已重置，请重新登录"})
}

// detectDeviceType 根据 User-Agent 检测设备类型
func detectDeviceType(userAgent string) string {
	if userAgent == "" {
		return "unknown"
	}

	if strings.Contains(userAgent, "iPad") || strings.Contains(userAgent, "Tablet") {
		return "tablet"
	}
	if strings.Contains(userAgent, "Mobile") || strings.Contains(userAgent, "Android") || strings.Contains(userAgent, "iPhone") {
		return "mobile"
	}

	return "desktop"
}
