package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// VerificationTokenRedis 邮箱验证令牌 Redis 操作
type VerificationTokenRedis struct {
	client redis.Cmdable
	prefix string
}

// NewVerificationTokenRedis 创建验证令牌 Redis 客户端
func NewVerificationTokenRedis(client redis.Cmdable) *VerificationTokenRedis {
	return &VerificationTokenRedis{
		client: client,
		prefix: "kiqi:email_verification:",
	}
}

// Set 存储令牌到 Redis（带 TTL）
func (r *VerificationTokenRedis) Set(ctx context.Context, token, userID string, ttlSec int64) error {
	return r.client.Set(ctx, r.prefix+token, userID, 0).Err()
}

// Get 从 Redis 获取令牌对应的用户 ID
// 返回空字符串表示令牌不存在或已过期
func (r *VerificationTokenRedis) Get(ctx context.Context, token string) (string, error) {
	return r.client.Get(ctx, r.prefix+token).Result()
}

// Del 删除 Redis 中的令牌
func (r *VerificationTokenRedis) Del(ctx context.Context, token string) error {
	return r.client.Del(ctx, r.prefix+token).Err()
}
