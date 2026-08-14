// Package crawl — Stream 消费者：消费 Python Worker 回传的文章和事件。
//
// crawl:article:ingest → 反序列化 → UpsertBatch → XACK
// crawl:task:event    → 反序列化 → 更新 TaskRun → XACK
package crawl

import (
	"context"
	"encoding/json"
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
)

// StartArticleConsumer 消费 crawl:article:ingest，批量写入 PG
func (s *Service) StartArticleConsumer(ctx context.Context) {
	stream := s.cfg.ArticleStream
	rdb := s.redis

	ensureGroup(rdb, stream, articleConsumerGroup)
	logger.Info("Article consumer started", "stream", stream)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Article consumer shutting down")
			return
		default:
		}

		results, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    articleConsumerGroup,
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
			logger.Error("Article consumer read error", "err", err)
			time.Sleep(time.Second)
			continue
		}

		for _, xr := range results {
			for _, msg := range xr.Messages {
				if err := s.handleArticleMessage(ctx, msg); err != nil {
					logger.Error("Article message handle error", "msg_id", msg.ID, "err", err)
					continue // 不 XACK，等待重试
				}
				rdb.XAck(ctx, stream, articleConsumerGroup, msg.ID)
			}
		}
	}
}

// StartEventConsumer 消费 crawl:task:event，更新 TaskRun 状态
func (s *Service) StartEventConsumer(ctx context.Context) {
	stream := s.cfg.EventStream
	rdb := s.redis

	ensureGroup(rdb, stream, eventConsumerGroup)
	logger.Info("Event consumer started", "stream", stream)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Event consumer shutting down")
			return
		default:
		}

		results, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    eventConsumerGroup,
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
			logger.Error("Event consumer read error", "err", err)
			time.Sleep(time.Second)
			continue
		}

		for _, xr := range results {
			for _, msg := range xr.Messages {
				if err := s.handleEventMessage(ctx, msg); err != nil {
					logger.Error("Event message handle error", "msg_id", msg.ID, "err", err)
					continue
				}
				rdb.XAck(ctx, stream, eventConsumerGroup, msg.ID)
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
		if _, err := s.articles.UpsertBatch(ctx, arts); err != nil {
			return err
		}
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
			return s.tasks.UpdateStatus(ctx, event.TaskID, event.Status, 0, 0, 0, "")
		}

	case "list_synced":
		// 列表同步完成
		return s.tasks.UpdateStatus(ctx, event.TaskID, "running", event.ListCount, event.ListCount, 0, "")

	case "task_done":
		now := time.Now()
		failed := event.DetailFailed
		_ = s.tasks.UpdateStatus(ctx, event.TaskID, event.Status, event.Total, event.ListCount, failed, "")
		// 更新完成时间（通过 FindByID + 手动更新）
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

	return nil
}

// updateTaskFinished 更新任务完成时间
func (s *Service) updateTaskFinished(ctx context.Context, taskID string, t *time.Time) error {
	// 通过仓储的 UpdateStatus 无法设置 finished_at，需要直接更新
	// 这里复用 UpdateStatus 的 error 字段传空，不影响
	return nil // finished_at 在 UpdateStatus 中暂未支持，后续可扩展
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
