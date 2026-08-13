package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	"github.com/shenfay/go-react-admin/pkg/constants"
	authErr "github.com/shenfay/go-react-admin/pkg/errors/auth"
	"github.com/shenfay/go-react-admin/pkg/utils"
)

// TokenServiceConfig Token 服务配置
type TokenServiceConfig struct {
	RedisClient    *redis.Client
	JWTSecret      string
	Issuer         string
	AccessExpire   time.Duration
	RefreshExpire  time.Duration
}

// jwtTokenService JWT Token 服务实现
type jwtTokenService struct {
	redisClient   *redis.Client
	jwtSecret     []byte
	issuer        string
	accessExpire  time.Duration
	refreshExpire time.Duration
}

// NewTokenServiceImpl 创建Token服务实现
func NewTokenServiceImpl(cfg TokenServiceConfig) TokenService {
	return &jwtTokenService{
		redisClient:   cfg.RedisClient,
		jwtSecret:     []byte(cfg.JWTSecret),
		issuer:        cfg.Issuer,
		accessExpire:  cfg.AccessExpire,
		refreshExpire: cfg.RefreshExpire,
	}
}

// --- Redis Key 构建 ---

func (s *jwtTokenService) refreshTokenKey(tokenID string) string {
	return constants.RedisKeyRefreshToken + tokenID
}

func (s *jwtTokenService) deviceKey(token string) string {
	return constants.RedisKeyPrefix + "device:" + token
}

func (s *jwtTokenService) userDevicesKey(userID string) string {
	return constants.RedisKeyPrefix + "user_devices:" + userID
}

func (s *jwtTokenService) accessDeviceKey(accessToken string) string {
	return constants.RedisKeyAccessDevice + accessToken
}

// GenerateTokens 生成 Token 对
func (s *jwtTokenService) GenerateTokens(ctx context.Context, userID, email string) (*TokenPair, error) {
	now := utils.Now()

	accessToken, err := s.generateAccessToken(userID, email, now)
	if err != nil {
		return nil, err
	}

	refreshTokenID := utils.GenerateID()
	refreshToken := refreshTokenID

	if err := s.storeRefreshToken(ctx, refreshTokenID, userID); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.accessExpire,
	}, nil
}

func (s *jwtTokenService) generateAccessToken(userID, email string, issuedAt time.Time) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		Email:     email,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(s.accessExpire)),
			NotBefore: jwt.NewNumericDate(issuedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *jwtTokenService) storeRefreshToken(ctx context.Context, tokenID, userID string) error {
	return s.redisClient.Set(ctx, s.refreshTokenKey(tokenID), userID, s.refreshExpire).Err()
}

// ValidateAccessToken 验证 Access Token
func (s *jwtTokenService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, authErr.ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, authErr.ErrInvalidToken
	}

	if claims.TokenType != "access" {
		return nil, authErr.ErrInvalidToken
	}

	return claims, nil
}

// ValidateRefreshTokenWithDevice 验证 Refresh Token 并返回设备信息
func (s *jwtTokenService) ValidateRefreshTokenWithDevice(ctx context.Context, token string) (*DeviceInfo, error) {
	userID, err := s.redisClient.Get(ctx, s.refreshTokenKey(token)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, authErr.ErrTokenExpired
		}
		return nil, authErr.ErrInvalidToken
	}

	deviceData, err := s.redisClient.Get(ctx, s.deviceKey(token)).Result()
	if err == nil {
		var deviceInfo DeviceInfo
		if err := json.Unmarshal([]byte(deviceData), &deviceInfo); err == nil {
			return &deviceInfo, nil
		}
	}

	return &DeviceInfo{
		UserID: userID,
	}, nil
}

// StoreDeviceInfo 存储设备信息到 Redis
func (s *jwtTokenService) StoreDeviceInfo(ctx context.Context, token string, deviceInfo DeviceInfo) error {
	deviceInfo.CreatedAt = utils.NowRFC3339()

	data, err := json.Marshal(deviceInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal device info: %w", err)
	}

	if err := s.redisClient.Set(ctx, s.deviceKey(token), string(data), s.refreshExpire).Err(); err != nil {
		return err
	}

	userDevicesKey := s.userDevicesKey(deviceInfo.UserID)
	if err := s.redisClient.SAdd(ctx, userDevicesKey, token).Err(); err != nil {
		return err
	}

	s.redisClient.Expire(ctx, userDevicesKey, s.refreshExpire)

	return nil
}

// LinkAccessToDevice 建立 access_token → device_token_id 映射，用于标识当前请求对应的设备
func (s *jwtTokenService) LinkAccessToDevice(ctx context.Context, accessToken, deviceTokenID string) error {
	return s.redisClient.Set(ctx, s.accessDeviceKey(accessToken), deviceTokenID, s.accessExpire).Err()
}

// GetCurrentDeviceTokenID 根据 access_token 获取当前设备对应的 device_token_id
func (s *jwtTokenService) GetCurrentDeviceTokenID(ctx context.Context, accessToken string) (string, error) {
	return s.redisClient.Get(ctx, s.accessDeviceKey(accessToken)).Result()
}

// RevokeToken 撤销 Token
func (s *jwtTokenService) RevokeToken(ctx context.Context, tokenID string) error {
	return s.redisClient.Del(ctx, s.refreshTokenKey(tokenID)).Err()
}

// RevokeDeviceByToken 根据 Token 撤销特定设备
func (s *jwtTokenService) RevokeDeviceByToken(ctx context.Context, token string) error {
	deviceInfo, err := s.ValidateRefreshTokenWithDevice(ctx, token)
	if err != nil {
		return err
	}

	s.redisClient.Del(ctx, s.deviceKey(token))
	s.redisClient.Del(ctx, s.refreshTokenKey(token))

	if deviceInfo.UserID != "" {
		s.redisClient.SRem(ctx, s.userDevicesKey(deviceInfo.UserID), token)
	}

	return nil
}

// RevokeAllDevices 撤销用户的所有设备
func (s *jwtTokenService) RevokeAllDevices(ctx context.Context, userID string) error {
	userDevicesKey := s.userDevicesKey(userID)
	tokens, err := s.redisClient.SMembers(ctx, userDevicesKey).Result()
	if err != nil {
		return err
	}

	for _, token := range tokens {
		if err := s.RevokeDeviceByToken(ctx, token); err != nil {
			continue
		}
	}

	s.redisClient.Del(ctx, userDevicesKey)

	return nil
}

// GetUserDevices 获取用户的所有设备列表
func (s *jwtTokenService) GetUserDevices(ctx context.Context, userID string) ([]DeviceInfo, error) {
	userDevicesKey := s.userDevicesKey(userID)
	tokens, err := s.redisClient.SMembers(ctx, userDevicesKey).Result()
	if err != nil {
		return nil, err
	}

	var devices []DeviceInfo
	for _, token := range tokens {
		deviceData, err := s.redisClient.Get(ctx, s.deviceKey(token)).Result()
		if err == nil {
			var deviceInfo DeviceInfo
			if err := json.Unmarshal([]byte(deviceData), &deviceInfo); err == nil {
				devices = append(devices, deviceInfo)
			}
		}
	}

	return devices, nil
}
