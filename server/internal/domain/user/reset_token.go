package user

import (
	"context"
	"time"
)

// ResetToken 密码重置令牌实体
type ResetToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

// ResetTokenRepository 密码重置令牌仓储接口
type ResetTokenRepository interface {
	// Create 创建重置令牌
	Create(ctx context.Context, token *ResetToken) error

	// FindByToken 根据令牌字符串查找
	FindByToken(ctx context.Context, token string) (*ResetToken, error)

	// MarkAsUsed 标记令牌为已使用
	MarkAsUsed(ctx context.Context, id string) error
}
