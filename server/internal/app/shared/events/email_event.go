package events

import "time"

// SendEmailEvent 发送邮件事件
// 用于通知 Worker 发送验证或重置邮件
type SendEmailEvent struct {
	EmailType string    `json:"email_type"` // "verification" 或 "password_reset"
	To        string    `json:"to"`
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

// NewSendEmailEvent 创建发送邮件事件
func NewSendEmailEvent(emailType, to, token, userID string) *SendEmailEvent {
	return &SendEmailEvent{
		EmailType: emailType,
		To:        to,
		Token:     token,
		UserID:    userID,
		Timestamp: time.Now(),
	}
}

// EventName 返回事件名称
func (e *SendEmailEvent) EventName() string {
	return "auth:send_email"
}

// OccurredAt 返回事件发生时间
func (e *SendEmailEvent) OccurredAt() time.Time {
	return e.Timestamp
}
