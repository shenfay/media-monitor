package handlers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/shenfay/go-react-admin/internal/app/port"
	appevents "github.com/shenfay/go-react-admin/internal/app/shared/events"
	"github.com/shenfay/go-react-admin/pkg/logger"
)

// SendEmailHandler 发送邮件 Worker 处理器
type SendEmailHandler struct {
	sender port.EmailSender
}

// NewSendEmailHandler 创建发送邮件处理器
func NewSendEmailHandler(sender port.EmailSender) *SendEmailHandler {
	return &SendEmailHandler{sender: sender}
}

// ProcessTask 处理发送邮件任务
func (h *SendEmailHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var evt appevents.SendEmailEvent
	if err := json.Unmarshal(task.Payload(), &evt); err != nil {
		logger.Error("Failed to unmarshal send email event", "error", err)
		return err
	}

	switch evt.EmailType {
	case "verification":
		logger.Info("Sending verification email", "to", evt.To, "user_id", evt.UserID)
		return h.sender.SendVerificationEmail(ctx, evt.To, evt.Token, evt.UserID)
	case "password_reset":
		logger.Info("Sending password reset email", "to", evt.To, "user_id", evt.UserID)
		return h.sender.SendPasswordResetEmail(ctx, evt.To, evt.Token, evt.UserID)
	default:
		logger.Warn("Unknown email type", "email_type", evt.EmailType)
		return nil
	}
}
