package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ResetTokenRedis 密码重置令牌 Redis 操作
type ResetTokenRedis struct {
	client redis.Cmdable
	prefix string
}

// NewResetTokenRedis 创建重置令牌 Redis 客户端
func NewResetTokenRedis(client redis.Cmdable) *ResetTokenRedis {
	return &ResetTokenRedis{
		client: client,
		prefix: "kiqi:password_reset:",
	}
}

// Set 存储令牌到 Redis（带 TTL）
func (r *ResetTokenRedis) Set(ctx context.Context, token, userID string, ttlSec int64) error {
	return r.client.Set(ctx, r.prefix+token, userID, 0).Err()
}

// Get 从 Redis 获取令牌对应的用户 ID
func (r *ResetTokenRedis) Get(ctx context.Context, token string) (string, error) {
	return r.client.Get(ctx, r.prefix+token).Result()
}

// Del 删除 Redis 中的令牌
func (r *ResetTokenRedis) Del(ctx context.Context, token string) error {
	return r.client.Del(ctx, r.prefix+token).Err()
}
