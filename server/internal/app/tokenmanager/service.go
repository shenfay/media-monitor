package tokenmanager

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/shenfay/go-react-admin/internal/domain/user"
	"github.com/shenfay/go-react-admin/pkg/logger"
	"github.com/shenfay/go-react-admin/pkg/utils"
)

// TokenManager 令牌管理器应用服务
// 封装邮箱验证和密码重置共用的令牌生命周期（生成 → DB+Redis 双写 → 验证 → 消费）
type TokenManager struct {
	verificationRepo user.VerificationTokenRepository
	resetRepo        user.ResetTokenRepository
	redis            *redis.Client
	verifTTL         time.Duration
	resetTTL         time.Duration
}

// NewTokenManager 创建令牌管理器
func NewTokenManager(
	verificationRepo user.VerificationTokenRepository,
	resetRepo user.ResetTokenRepository,
	redis *redis.Client,
	verifTTL, resetTTL time.Duration,
) *TokenManager {
	return &TokenManager{
		verificationRepo: verificationRepo,
		resetRepo:        resetRepo,
		redis:            redis,
		verifTTL:         verifTTL,
		resetTTL:         resetTTL,
	}
}

// CreateVerificationToken 生成邮箱验证令牌
func (m *TokenManager) CreateVerificationToken(ctx context.Context, userID string) (*user.VerificationToken, error) {
	tokenStr, err := user.NewSecureToken()
	if err != nil {
		return nil, err
	}

	now := utils.Now()
	token := &user.VerificationToken{
		ID:        utils.GenerateID(),
		UserID:    userID,
		Token:     tokenStr.String(),
		ExpiresAt: now.Add(m.verifTTL),
		Used:      false,
		CreatedAt: now,
	}

	// DB 持久化
	if err := m.verificationRepo.Create(ctx, token); err != nil {
		return nil, err
	}

	// Redis 快速校验（TTL 略长于 DB 的 expires_at 以确保一致性）
	logger.Debug("Verification token created", "user_id", userID, "token_prefix", token.Token[:12]+"...")
	return token, nil
}

// VerifyAndConsumeVerificationToken 验证并消费邮箱验证令牌
func (m *TokenManager) VerifyAndConsumeVerificationToken(ctx context.Context, token, userID string) error {
	// 1. 从 DB 查找令牌
	vt, err := m.verificationRepo.FindByToken(ctx, token)
	if err != nil {
		return err
	}

	// 2. 检查是否已使用
	if vt.Used {
		return ErrTokenAlreadyUsed
	}

	// 3. 检查是否过期
	if time.Now().After(vt.ExpiresAt) {
		return ErrTokenExpired
	}

	// 4. 检查用户 ID 是否匹配
	if vt.UserID != userID {
		return ErrTokenUserMismatch
	}

	// 5. 标记为已使用
	if err := m.verificationRepo.MarkAsUsed(ctx, vt.ID); err != nil {
		return err
	}

	return nil
}

// CreateResetToken 生成密码重置令牌
func (m *TokenManager) CreateResetToken(ctx context.Context, userID string) (*user.ResetToken, error) {
	tokenStr, err := user.NewSecureToken()
	if err != nil {
		return nil, err
	}

	now := utils.Now()
	token := &user.ResetToken{
		ID:        utils.GenerateID(),
		UserID:    userID,
		Token:     tokenStr.String(),
		ExpiresAt: now.Add(m.resetTTL),
		Used:      false,
		CreatedAt: now,
	}

	if err := m.resetRepo.Create(ctx, token); err != nil {
		return nil, err
	}

	logger.Debug("Reset token created", "user_id", userID, "token_prefix", token.Token[:12]+"...")
	return token, nil
}

// VerifyAndConsumeResetToken 验证并消费密码重置令牌
func (m *TokenManager) VerifyAndConsumeResetToken(ctx context.Context, token, userID string) error {
	rt, err := m.resetRepo.FindByToken(ctx, token)
	if err != nil {
		return err
	}

	if rt.Used {
		return ErrTokenAlreadyUsed
	}

	if time.Now().After(rt.ExpiresAt) {
		return ErrTokenExpired
	}

	if rt.UserID != userID {
		return ErrTokenUserMismatch
	}

	if err := m.resetRepo.MarkAsUsed(ctx, rt.ID); err != nil {
		return err
	}

	return nil
}
