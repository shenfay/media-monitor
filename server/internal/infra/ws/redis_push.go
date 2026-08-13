package ws

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/redis/go-redis/v9"

	"github.com/shenfay/go-react-admin/pkg/logger"
)

const pushChannel = "push:realtime"

// pushEnvelope Redis Pub/Sub 推送信封
type pushEnvelope struct {
	UserID string          `json:"user_id"`
	Data   json.RawMessage `json:"data"`
}

// RedisPublisher 基于 Redis Pub/Sub 的推送发布器
// 实现 port.RealtimePusher 接口，支持跨进程推送
type RedisPublisher struct {
	client *redis.Client
}

// NewRedisPublisher 创建 Redis 推送发布器
func NewRedisPublisher(client *redis.Client) *RedisPublisher {
	return &RedisPublisher{client: client}
}

// SendToUser 向指定用户推送消息（通过 Redis Pub/Sub 广播）
func (p *RedisPublisher) SendToUser(userID string, msg []byte) {
	envelope := pushEnvelope{
		UserID: userID,
		Data:   msg,
	}
	payload, err := json.Marshal(envelope)
	if err != nil {
		logger.Error("Failed to marshal push envelope", "error", err)
		return
	}

	ctx := context.Background()
	if err := p.client.Publish(ctx, pushChannel, payload).Err(); err != nil {
		logger.Error("Failed to publish push message", "error", err, "user_id", userID)
	}
}

// RedisSubscriber 基于 Redis Pub/Sub 的推送订阅器
// 订阅 Redis 频道，将消息转发到本地 Hub
type RedisSubscriber struct {
	client *redis.Client
	hub    *Hub
	sub    *redis.PubSub
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewRedisSubscriber 创建 Redis 推送订阅器
func NewRedisSubscriber(client *redis.Client, hub *Hub) *RedisSubscriber {
	return &RedisSubscriber{client: client, hub: hub}
}

// Start 启动订阅（非阻塞）
func (s *RedisSubscriber) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	sub := s.client.Subscribe(ctx, pushChannel)
	s.sub = sub

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ch := sub.Channel()
		for msg := range ch {
			s.handleMessage(msg)
		}
	}()

	logger.Info("Redis Pub/Sub subscriber started", "channel", pushChannel)
}

// Stop 停止订阅
func (s *RedisSubscriber) Stop() {
	// 显式关闭 PubSub，确保 ch 关闭 → goroutine 退出
	if s.sub != nil {
		if err := s.sub.Close(); err != nil {
			logger.Error("Failed to close PubSub subscription", "error", err)
		}
	}
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
	logger.Info("Redis Pub/Sub subscriber stopped")
}

// handleMessage 处理收到的推送消息
func (s *RedisSubscriber) handleMessage(msg *redis.Message) {
	var envelope pushEnvelope
	if err := json.Unmarshal([]byte(msg.Payload), &envelope); err != nil {
		logger.Error("Failed to unmarshal push envelope", "error", err)
		return
	}

	s.hub.SendToUser(envelope.UserID, envelope.Data)
}
