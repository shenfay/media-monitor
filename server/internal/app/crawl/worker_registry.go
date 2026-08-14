package crawl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	adapterKeyPrefix  = "scraper:adapters:"
	workerKeyPrefix   = "scraper:workers:"
	workerIDsKey      = "scraper:worker:ids"
	controlChannel    = "crawl:worker:control"
	configChannel     = "crawl:config:changed"
)

// WorkerInfo Worker 实例信息（从 Redis 读取）
type WorkerInfo struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Adapters       []string  `json:"adapters"`
	Capabilities   []string  `json:"capabilities"`
	Status         string    `json:"status"`
	Concurrency    int       `json:"concurrency"`
	CurrentTask    string    `json:"current_task"`
	ProcessedCount int64     `json:"processed_count"`
	LastHeartbeat  time.Time `json:"last_heartbeat"`
	StartedAt      time.Time `json:"started_at"`
}

// AdapterMeta 适配器元数据（从 Redis 读取）
type AdapterMeta struct {
	Name         string    `json:"name"`
	RequiredTags []string  `json:"required_tags"`
	PlatformType string    `json:"platform_type"`
	ClassName    string    `json:"class_name"`
	FirstSeen    time.Time `json:"first_seen"`
}

// WorkerRegistry Worker 注册表（从 Redis 读取 Worker/Adapter 信息）
type WorkerRegistry struct {
	rdb *redis.Client
}

// NewWorkerRegistry 创建 Worker 注册表
func NewWorkerRegistry(rdb *redis.Client) *WorkerRegistry {
	return &WorkerRegistry{rdb: rdb}
}

// ListWorkers 从 Redis 读取所有 Worker 实例
func (r *WorkerRegistry) ListWorkers(ctx context.Context) ([]WorkerInfo, error) {
	// 获取所有 worker ID
	ids, err := r.rdb.SMembers(ctx, workerIDsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get worker IDs: %w", err)
	}

	workers := make([]WorkerInfo, 0, len(ids))
	for _, id := range ids {
		key := workerKeyPrefix + id
		data, err := r.rdb.HGetAll(ctx, key).Result()
		if err != nil {
			continue
		}
		if len(data) == 0 {
			// key 已过期（TTL 到期），清理 ids 集合
			r.rdb.SRem(ctx, workerIDsKey, id)
			continue
		}

		w := WorkerInfo{ID: id}
		w.Name = data["name"]
		w.Status = data["status"]
		w.CurrentTask = data["current_task"]

		if data["adapters"] != "" {
			_ = json.Unmarshal([]byte(data["adapters"]), &w.Adapters)
		}
		if data["capabilities"] != "" {
			_ = json.Unmarshal([]byte(data["capabilities"]), &w.Capabilities)
		}
		if data["concurrency"] != "" {
			fmt.Sscanf(data["concurrency"], "%d", &w.Concurrency)
		}
		if data["processed_count"] != "" {
			fmt.Sscanf(data["processed_count"], "%d", &w.ProcessedCount)
		}
		if data["last_heartbeat"] != "" {
			var ts int64
			fmt.Sscanf(data["last_heartbeat"], "%d", &ts)
			if ts > 0 {
				w.LastHeartbeat = time.Unix(ts, 0)
			}
		}
		if data["started_at"] != "" {
			if t, err := time.Parse(time.RFC3339, data["started_at"]); err == nil {
				w.StartedAt = t
			}
		}

		// 检查是否已离线（TTL 过期但 SRem 尚未执行）
		ttl := r.rdb.TTL(ctx, key).Val()
		if ttl <= 0 {
			w.Status = "offline"
			r.rdb.SRem(ctx, workerIDsKey, id)
		}

		workers = append(workers, w)
	}
	return workers, nil
}

// ListAdapters 列出所有已知 adapter 元数据
func (r *WorkerRegistry) ListAdapters(ctx context.Context) ([]AdapterMeta, error) {
	// 扫描所有 scraper:adapters:* key
	var adapters []AdapterMeta
	iter := r.rdb.Scan(ctx, 0, adapterKeyPrefix+"*", 100).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		data, err := r.rdb.HGetAll(ctx, key).Result()
		if err != nil || len(data) == 0 {
			continue
		}

		name := key[len(adapterKeyPrefix):]
		meta := AdapterMeta{Name: name}
		meta.PlatformType = data["platform_type"]
		meta.ClassName = data["class_name"]

		if data["required_tags"] != "" {
			_ = json.Unmarshal([]byte(data["required_tags"]), &meta.RequiredTags)
		}
		if data["first_seen"] != "" {
			if t, err := time.Parse(time.RFC3339, data["first_seen"]); err == nil {
				meta.FirstSeen = t
			}
		}
		adapters = append(adapters, meta)
	}
	if err := iter.Err(); err != nil {
		return adapters, fmt.Errorf("failed to scan adapters: %w", err)
	}
	return adapters, nil
}

// SendCommand 向 Worker 发送控制指令（通过 Redis Pub/Sub）
func (r *WorkerRegistry) SendCommand(ctx context.Context, workerID, action string) error {
	msg := map[string]string{
		"worker_id": workerID,
		"action":    action, // pause | resume | shutdown
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.rdb.Publish(ctx, controlChannel, payload).Err()
}

// PublishConfigChanged 发布配置变更通知（Worker 收到后重新计算流列表）
func (r *WorkerRegistry) PublishConfigChanged(ctx context.Context) error {
	msg := `{"action":"source_changed"}`
	return r.rdb.Publish(ctx, configChannel, msg).Err()
}
