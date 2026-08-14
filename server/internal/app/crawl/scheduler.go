package crawl

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"

	"github.com/shenfay/go-react-admin/pkg/logger"
)

// SchedulerTaskType Asynq 调度任务类型（cron 触发后创建 Redis Stream 任务）
const SchedulerTaskType = "crawl:schedule"

// StartScheduler 启动基于 asynq.Scheduler 的定时抓取调度。
//
// 设计要点：调度器只负责在 cron 到点时「创建 TaskRun + 写入 Redis Stream」，
// 真正的抓取由 Python 常驻服务消费执行，避免 Go worker 阻塞在长耗时抓取上。
func StartScheduler(redisOpt asynq.RedisClientOpt, svc *Service) *asynq.Scheduler {
	sched := asynq.NewScheduler(redisOpt, &asynq.SchedulerOpts{Location: time.UTC})

	registerAll := func() {
		sources, err := svc.ListSources(context.Background(), true)
		if err != nil {
			logger.Error("crawl scheduler: load enabled sources failed", "err", err)
			return
		}
		for _, src := range sources {
			if src.Schedule == "" {
				continue
			}
			payload, _ := json.Marshal(map[string]string{"source_id": src.ID})
			task := asynq.NewTask(SchedulerTaskType, payload)
			// 同名已注册会返回错误；忽略之（运行时变更需等待下次 reload 或重启）
			if _, err := sched.Register(src.Schedule, task); err != nil {
				logger.Warn("crawl scheduler: register skipped", "source", src.ID, "err", err)
			}
		}
	}

	registerAll()
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			registerAll()
		}
	}()

	go func() {
		if err := sched.Run(); err != nil {
			logger.Error("crawl scheduler stopped", "err", err)
		}
	}()

	return sched
}

// ScheduleHandler 处理 cron 触发的 Asynq 任务：入队一次抓取
func ScheduleHandler(svc *Service) func(ctx context.Context, t *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var p struct {
			SourceID string `json:"source_id"`
		}
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}
		_, err := svc.EnqueueForSource(ctx, p.SourceID, "cron")
		return err
	}
}
