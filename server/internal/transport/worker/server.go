package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-react-admin/pkg/logger"
)

// Server Worker 服务器
type Server struct {
	srv *asynq.Server
}

// NewServer 创建 Worker 服务器
func NewServer(redisAddr, redisPwd string, db int, concurrency int) *Server {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPwd,
			DB:       db,
		},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			StrictPriority: true,
		},
	)

	return &Server{srv: srv}
}

// RegisterHandler 注册自定义任务处理器
func (s *Server) RegisterHandler(mux *asynq.ServeMux, taskType string, handler func(ctx context.Context, task *asynq.Task) error) {
	mux.HandleFunc(taskType, handler)
	logger.Info("Registered handler for type: ", taskType)
}

// Start 启动 Worker
func (s *Server) Start(mux *asynq.ServeMux) error {
	logger.Info("Starting Asynq Worker processor...")
	return s.srv.Run(mux)
}

// Shutdown 优雅关闭
func (s *Server) Shutdown() {
	s.srv.Shutdown()
	logger.Info("Worker stopped gracefully")
}
