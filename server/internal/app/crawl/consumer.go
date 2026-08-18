// Package crawl — Stream 消费者：消费 Python Worker 回传的文章和事件。
//
// crawl:article:ingest → 反序列化 → UpsertBatch → XACK
// crawl:task:event    → 反序列化 → 更新 TaskRun → XACK
//
// 可靠性保障：
//   - 单实例锁：Redis SET NX 防止多个 Go Worker 同时消费
//   - 心跳上报：定期写入 Redis，管理面板可查看 consumer 状态
//   - 看门狗：定期检查队列积压和 worker 离线，输出告警日志
package crawl

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/shenfay/go-react-admin/internal/domain/crawl"
	"github.com/shenfay/go-react-admin/pkg/logger"
)

const (
	// articleConsumerGroup 文章消费组名
	articleConsumerGroup = "crawl:go:article"
	// eventConsumerGroup 事件消费组名
	eventConsumerGroup = "crawl:go:event"
	// consumerName 消费者名（单实例 Go worker）
	consumerName = "go-worker-1"
	// streamBatchSize 每次读取消息数
	streamBatchSize = 10
	// streamBlockMs 阻塞等待毫秒
	streamBlockMs = 5000

	// 单实例锁
	consumerLockKey    = "crawl:go:consumer:lock"
	consumerLockTTL    = 120 * time.Second
	consumerLockRenew  = 60 * time.Second

	// 心跳
	consumerHeartbeatKey    = "crawl:go:consumer:heartbeat"
	consumerHeartbeatTTL    = 90 * time.Second
	consumerHeartbeatInterval = 30 * time.Second

	// 看门狗
	watchdogInterval    = 60 * time.Second
	watchdogDetailLagThreshold   = 100
	watchdogIngestPendingThreshold = 200
)

// StartConsumers 启动所有 Stream 消费者（含单实例锁、心跳、看门狗）。
// 如果已有其他 Go Worker 持有锁，将打印错误并 os.Exit(1)。
func (s *Service) StartConsumers(ctx context.Context) {
	// 1. 获取单实例锁
	if err := s.acquireConsumerLock(ctx); err != nil {
		logger.Error("Failed to acquire consumer lock — another Go worker is running?", "err", err)
		os.Exit(1)
	}
	logger.Info("Consumer lock acquired, starting consumers...")

	// 2. 确保 consumer group 存在
	ensureGroup(s.redis, s.cfg.ArticleStream, articleConsumerGroup)
	ensureGroup(s.redis, s.cfg.EventStream, eventConsumerGroup)

	// 3. 启动消费循环
	go s.runConsumerLoop(ctx, s.cfg.ArticleStream, articleConsumerGroup, s.handleArticleMessage, "Article")
	go s.runConsumerLoop(ctx, s.cfg.EventStream, eventConsumerGroup, s.handleEventMessage, "Event")
	logger.Info("Stream consumers started (article + event)")

	// 4. 启动心跳
	go s.startConsumerHeartbeat(ctx)

	// 5. 启动看门狗
	go s.startWatchdog(ctx)
}

// ─── 单实例锁 ────────────────────────────────────────────────────────────────

// acquireConsumerLock 尝试获取消费者分布式锁（SET NX），失败时返回错误。
func (s *Service) acquireConsumerLock(ctx context.Context) error {
	ok, err := s.redis.SetNX(ctx, consumerLockKey, fmt.Sprintf("%d", os.Getpid()), consumerLockTTL).Result()
	if err != nil {
		return fmt.Errorf("redis SET NX failed: %w", err)
	}
	if !ok {
		// 读取当前锁持有者
		holder, _ := s.redis.Get(ctx, consumerLockKey).Result()
		return fmt.Errorf("lock already held by PID %s", holder)
	}
	// 后台续期
	go s.renewConsumerLock(ctx)
	return nil
}

// renewConsumerLock 每 60s 续期锁 TTL，直到 ctx 取消。
func (s *Service) renewConsumerLock(ctx context.Context) {
	ticker := time.NewTicker(consumerLockRenew)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			s.releaseConsumerLock(context.Background())
			return
		case <-ticker.C:
			if err := s.redis.Expire(ctx, consumerLockKey, consumerLockTTL).Err(); err != nil {
				logger.Warn("Failed to renew consumer lock", "err", err)
			}
		}
	}
}

// releaseConsumerLock 释放锁（进程退出时调用）。
func (s *Service) releaseConsumerLock(ctx context.Context) {
	// 仅删除自己持有的锁（Lua 脚本保证原子性）
	script := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		end
		return 0
	`)
	result, err := script.Run(ctx, s.redis, []string{consumerLockKey}, fmt.Sprintf("%d", os.Getpid())).Int()
	if err != nil {
		logger.Warn("Failed to release consumer lock", "err", err)
		return
	}
	if result == 1 {
		logger.Info("Consumer lock released")
	}
}

// ─── 心跳上报 ────────────────────────────────────────────────────────────────

// articlesProcessed / eventsProcessed 由消费循环原子递增
var articlesProcessed int64
var eventsProcessed int64

// IncArticlesProcessed 原子递增文章消费计数（由 handleArticleMessage 调用）
func IncArticlesProcessed(n int) { atomic.AddInt64(&articlesProcessed, int64(n)) }

// IncEventsProcessed 原子递增事件消费计数（由 handleEventMessage 调用）
func IncEventsProcessed() { atomic.AddInt64(&eventsProcessed, 1) }

// startConsumerHeartbeat 每 30s 向 Redis 写入 consumer 状态。
func (s *Service) startConsumerHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(consumerHeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			// 清除心跳 key
			s.redis.Del(ctx, consumerHeartbeatKey)
			return
		case <-ticker.C:
			data := map[string]interface{}{
				"last_heartbeat":     fmt.Sprintf("%d", time.Now().Unix()),
				"articles_processed": fmt.Sprintf("%d", atomic.LoadInt64(&articlesProcessed)),
				"events_processed":   fmt.Sprintf("%d", atomic.LoadInt64(&eventsProcessed)),
				"status":             "running",
				"pid":                fmt.Sprintf("%d", os.Getpid()),
				"consumer_name":      consumerName,
			}
			pipe := s.redis.Pipeline()
			pipe.HSet(ctx, consumerHeartbeatKey, data)
			pipe.Expire(ctx, consumerHeartbeatKey, consumerHeartbeatTTL)
			if _, err := pipe.Exec(ctx); err != nil {
				logger.Warn("Consumer heartbeat write failed", "err", err)
			}
		}
	}
}

// ─── 看门狗 ──────────────────────────────────────────────────────────────────

// startWatchdog 每 60s 检查关键指标，异常时输出 Warn 日志。
func (s *Service) startWatchdog(ctx context.Context) {
	ticker := time.NewTicker(watchdogInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.runWatchdogChecks(ctx)
		}
	}
}

func (s *Service) runWatchdogChecks(ctx context.Context) {
	// 1. 检查锁是否仍由自己持有
	holder, err := s.redis.Get(ctx, consumerLockKey).Result()
	myPID := fmt.Sprintf("%d", os.Getpid())
	if err != nil {
		logger.Warn("[watchdog] Consumer lock MISSING — another instance may have taken over", "err", err)
	} else if holder != myPID {
		logger.Warn("[watchdog] Consumer lock held by DIFFERENT PID — this instance is stale!",
			"my_pid", myPID, "lock_holder", holder)
	}

	// 2. 检查 detail queue lag
	detailLag := s.getStreamLag(ctx, s.cfg.DetailQueue, "crawl:detail:worker")
	if detailLag > watchdogDetailLagThreshold {
		logger.Warn("[watchdog] Detail queue backlog",
			"queue", s.cfg.DetailQueue, "lag", detailLag, "threshold", watchdogDetailLagThreshold)
	}

	// 3. 检查 ingest queue pending
	ingestPending := s.getGroupPending(ctx, s.cfg.ArticleStream, articleConsumerGroup)
	if ingestPending > watchdogIngestPendingThreshold {
		logger.Warn("[watchdog] Ingest queue pending backlog",
			"stream", s.cfg.ArticleStream, "pending", ingestPending, "threshold", watchdogIngestPendingThreshold)
	}

	// 4. 检查 Python worker 是否在线
	onlineWorkers := s.countOnlinePythonWorkers(ctx)
	if onlineWorkers == 0 {
		logger.Warn("[watchdog] No Python list worker online (scraper:worker:ids empty or all expired)")
	}

	// 5. 检查 detail worker 心跳
	detailWorkerOffline := s.isDetailWorkerOffline(ctx)
	if detailWorkerOffline {
		logger.Warn("[watchdog] Detail worker OFFLINE (heartbeat expired or missing)")
	}
}

// getStreamLag 获取 stream 的未读消息数（XINFO GROUPS → lag 字段）。
func (s *Service) getStreamLag(ctx context.Context, stream, group string) int64 {
	groups, err := s.redis.XInfoGroups(ctx, stream).Result()
	if err != nil {
		return 0
	}
	for _, g := range groups {
		if g.Name == group {
			return g.Lag
		}
	}
	return 0
}

// getGroupPending 获取 consumer group 的 pending 数（已读未 ACK）。
func (s *Service) getGroupPending(ctx context.Context, stream, group string) int64 {
	groups, err := s.redis.XInfoGroups(ctx, stream).Result()
	if err != nil {
		return 0
	}
	for _, g := range groups {
		if g.Name == group {
			return g.Pending
		}
	}
	return 0
}

// countOnlinePythonWorkers 检查有多少 Python list worker 在线。
func (s *Service) countOnlinePythonWorkers(ctx context.Context) int {
	ids, err := s.redis.SMembers(ctx, workerIDsKey).Result()
	if err != nil {
		return 0
	}
	count := 0
	for _, id := range ids {
		key := workerKeyPrefix + id
		ttl := s.redis.TTL(ctx, key).Val()
		if ttl > 0 {
			count++
		}
	}
	return count
}

// isDetailWorkerOffline 检查 detail worker 心跳是否过期。
func (s *Service) isDetailWorkerOffline(ctx context.Context) bool {
	// detail worker 注册时写入 scraper:workers:{id}，adapters 含 "detail"
	// 检查所有 worker 中是否有 status=online 且 capabilities 含 detail 的
	ids, err := s.redis.SMembers(ctx, workerIDsKey).Result()
	if err != nil {
		return false // 无法判断时不告警
	}
	for _, id := range ids {
		key := workerKeyPrefix + id
		data, err := s.redis.HGetAll(ctx, key).Result()
		if err != nil || len(data) == 0 {
			continue
		}
		// 检查 capabilities 是否包含 detail-worker 标识
		caps := data["capabilities"]
		var capList []string
		_ = json.Unmarshal([]byte(caps), &capList)
		for _, c := range capList {
			if c == "detail-worker" {
				ttl := s.redis.TTL(ctx, key).Val()
				return ttl <= 0 // TTL 过期 = offline
			}
		}
	}
	// 没有找到 detail worker 注册记录，也算 offline
	return true
}

// ─── 消费循环 ────────────────────────────────────────────────────────────────

// runConsumerLoop 通用 Stream 消费循环
func (s *Service) runConsumerLoop(ctx context.Context, stream, group string, handler func(context.Context, redis.XMessage) error, label string) {
	rdb := s.redis
	for {
		select {
		case <-ctx.Done():
			logger.Info(label + " consumer shutting down")
			return
		default:
		}

		results, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumerName,
			Streams:  []string{stream, ">"},
			Count:    streamBatchSize,
			Block:    time.Duration(streamBlockMs) * time.Millisecond,
		}).Result()

		if err != nil {
			if err == redis.Nil {
				continue
			}
			if isContextDone(err) {
				return
			}
			logger.Error(label+" consumer read error", "err", err)
			time.Sleep(time.Second)
			continue
		}

		for _, xr := range results {
			for _, msg := range xr.Messages {
				if err := handler(ctx, msg); err != nil {
					logger.Error(label+" message handle error", "msg_id", msg.ID, "err", err)
					continue // 不 XACK，等待重试
				}
				rdb.XAck(ctx, stream, group, msg.ID)
			}
		}
	}
}

// handleArticleMessage 处理单条文章消息
func (s *Service) handleArticleMessage(ctx context.Context, msg redis.XMessage) error {
	payloadStr := getField(msg, "payload")
	if payloadStr == "" {
		return nil
	}

	var ingest ArticleIngestMessage
	if err := json.Unmarshal([]byte(payloadStr), &ingest); err != nil {
		logger.Error("Failed to unmarshal article message", "err", err)
		return nil // 格式错误，跳过不重试
	}

	arts := make([]*crawl.Article, 0, len(ingest.Articles))
	for _, in := range ingest.Articles {
		a, err := mapArticleInput(in)
		if err != nil {
			logger.Warn("Skip invalid article", "url", in.URL, "err", err)
			continue
		}
		arts = append(arts, a)
	}

	if len(arts) > 0 {
		var err error
		if ingest.Phase == "detail" {
			// detail 阶段仅更新正文相关字段，不覆盖列表阶段的元数据
			_, err = s.articles.UpsertDetailBatch(ctx, arts)
		} else {
			// list 阶段更新全部字段
			_, err = s.articles.UpsertBatch(ctx, arts)
		}
		if err != nil {
			return err
		}
		IncArticlesProcessed(len(arts))
	}

	logger.Debug("Article batch ingested",
		"task_id", ingest.TaskID, "phase", ingest.Phase,
		"batch_seq", ingest.BatchSeq, "count", len(arts))
	return nil
}

// handleEventMessage 处理单条事件消息
func (s *Service) handleEventMessage(ctx context.Context, msg redis.XMessage) error {
	payloadStr := getField(msg, "payload")
	if payloadStr == "" {
		return nil
	}

	var event TaskEventMessage
	if err := json.Unmarshal([]byte(payloadStr), &event); err != nil {
		logger.Error("Failed to unmarshal event message", "err", err)
		return nil
	}

	switch event.Type {
	case "status":
		// 状态变更（如 running）
		if event.Status != "" {
			if err := s.tasks.UpdateStatus(ctx, event.TaskID, event.Status, 0, 0, 0, ""); err != nil {
				return err
			}
			// 状态变为 running 时记录开始时间
			if event.Status == "running" {
				now := time.Now()
				return s.tasks.UpdateTimestamps(ctx, event.TaskID, &now, nil)
			}
			return nil
		}

	case "list_synced":
		// 列表同步完成
		if err := s.tasks.UpdateStatus(ctx, event.TaskID, "running", event.ListCount, event.ListCount, 0, ""); err != nil {
			return err
		}
		// 列表同步完成意味着任务已开始运行，记录开始时间
		now := time.Now()
		return s.tasks.UpdateTimestamps(ctx, event.TaskID, &now, nil)

	case "task_done":
		now := time.Now()
		failed := event.DetailFailed
		if err := s.tasks.UpdateStatus(ctx, event.TaskID, event.Status, event.Total, event.ListCount, failed, ""); err != nil {
			logger.Error("Failed to update task_done status", "task_id", event.TaskID, "status", event.Status, "err", err)
		}
		return s.updateTaskFinished(ctx, event.TaskID, &now)

	case "task_failed":
		now := time.Now()
		_ = s.tasks.UpdateStatus(ctx, event.TaskID, "failed", 0, 0, 0, event.Error)
		return s.updateTaskFinished(ctx, event.TaskID, &now)

	case "phase_start", "phase_done", "detail_progress":
		// 进度事件：仅日志记录（可扩展为 WebSocket 通知）
		logger.Debug("Task event",
			"task_id", event.TaskID, "type", event.Type,
			"phase", event.Phase, "status", event.Status,
			"detail_count", event.DetailCount, "detail_failed", event.DetailFailed)
	}

	IncEventsProcessed()
	return nil
}

// updateTaskFinished 更新任务完成时间
func (s *Service) updateTaskFinished(ctx context.Context, taskID string, t *time.Time) error {
	return s.tasks.UpdateTimestamps(ctx, taskID, nil, t)
}

// ensureGroup 确保 Consumer Group 存在
func ensureGroup(rdb *redis.Client, stream, group string) {
	err := rdb.XGroupCreate(context.Background(), stream, group, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		// 如果 stream 不存在则先创建
		rdb.XAdd(context.Background(), &redis.XAddArgs{
			Stream: stream,
			Values: map[string]interface{}{"init": "1"},
		})
		rdb.XGroupCreate(context.Background(), stream, group, "0")
	}
}

// getField 从 XMessage 中取字段值（兼容 string/[]byte）
func getField(msg redis.XMessage, key string) string {
	v, ok := msg.Values[key]
	if !ok {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	default:
		return ""
	}
}

// isContextDone 检查错误是否为 context 取消
func isContextDone(err error) bool {
	if err == nil {
		return false
	}
	return err == context.Canceled || err == context.DeadlineExceeded
}
