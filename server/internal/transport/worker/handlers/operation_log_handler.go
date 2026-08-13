package handlers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"

	appevents "github.com/shenfay/go-react-admin/internal/app/shared/events"
	"github.com/shenfay/go-react-admin/internal/domain/operation"
	"github.com/shenfay/go-react-admin/pkg/logger"
)

// OperationLogHandler 统一操作日志 Worker 处理器
type OperationLogHandler struct {
	repo operation.LogRepository
}

// NewOperationLogHandler 创建操作日志处理器
func NewOperationLogHandler(repo operation.LogRepository) *OperationLogHandler {
	return &OperationLogHandler{repo: repo}
}

// ProcessTask 处理操作日志任务
func (h *OperationLogHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	return h.processOperationLog(ctx, task)
}

// processOperationLog 处理统一操作日志
func (h *OperationLogHandler) processOperationLog(ctx context.Context, task *asynq.Task) error {
	var evt appevents.OperationEvent
	if err := json.Unmarshal(task.Payload(), &evt); err != nil {
		logger.Error("Failed to unmarshal operation log payload", "error", err)
		return err
	}

	log := &operation.OperationLog{
		UserID:    evt.UserID,
		Email:     evt.Email,
		Action:    evt.Action,
		Category:  evt.Category,
		Status:    evt.Status,
		IP:        evt.IP,
		UserAgent: evt.UserAgent,
		Device:    evt.Device,
		Browser:   evt.Browser,
		OS:        evt.OS,
		Metadata:  evt.Metadata,
	}

	if err := h.repo.Save(ctx, log); err != nil {
		logger.Error("Failed to save operation log",
			"action", log.Action,
			"error", err,
		)
		return err
	}

	logger.Debug("Operation log saved",
		"action", log.Action,
		"user_id", log.UserID,
	)
	return nil
}
