package user

import (
	"context"
	"time"
)

// VerificationToken 邮箱验证令牌实体
type VerificationToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

// VerificationTokenRepository 邮箱验证令牌仓储接口
type VerificationTokenRepository interface {
	// Create 创建验证令牌
	Create(ctx context.Context, token *VerificationToken) error

	// FindByToken 根据令牌字符串查找
	FindByToken(ctx context.Context, token string) (*VerificationToken, error)

	// MarkAsUsed 标记令牌为已使用
	MarkAsUsed(ctx context.Context, id string) error
}
