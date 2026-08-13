package handlers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"

	"github.com/shenfay/go-react-admin/internal/domain/notification"
	"github.com/shenfay/go-react-admin/pkg/logger"
)

// NotificationHandler 消息通知 Worker 处理器
type NotificationHandler struct {
	repo notification.MessageRepository
}

// NewNotificationHandler 创建消息通知处理器
func NewNotificationHandler(repo notification.MessageRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

// notificationPayload 消息通知任务载荷
type notificationPayload struct {
	SenderID    *string                `json:"sender_id,omitempty"`
	RecipientID string                 `json:"recipient_id"`
	Type        string                 `json:"type"`
	Category    string                 `json:"category"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	RefType     string                 `json:"ref_type,omitempty"`
	RefID       string                 `json:"ref_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessTask 处理消息通知任务
func (h *NotificationHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload notificationPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logger.Error("Failed to unmarshal notification payload", "error", err)
		return err
	}

	msg := &notification.Message{
		SenderID:    payload.SenderID,
		RecipientID: payload.RecipientID,
		Type:        notification.MessageType(payload.Type),
		Category:    notification.MessageCategory(payload.Category),
		Title:       payload.Title,
		Content:     payload.Content,
		RefType:     payload.RefType,
		RefID:       payload.RefID,
		Metadata:    payload.Metadata,
		IsRead:      false,
	}

	if err := h.repo.Save(ctx, msg); err != nil {
		logger.Error("Failed to save notification message",
			"type", payload.Type,
			"category", payload.Category,
			"recipient_id", payload.RecipientID,
			"error", err,
		)
		return err
	}

	logger.Debug("Notification message saved",
		"type", payload.Type,
		"category", payload.Category,
		"recipient_id", payload.RecipientID,
	)
	return nil
}
